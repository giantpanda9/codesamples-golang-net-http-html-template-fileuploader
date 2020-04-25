package main

import (
	"fmt"
	"io"
	"image"
	"path/filepath"
	"os"
	"syscall"
	"time"
	"strconv"
	"github.com/gosexy/exif"
	"github.com/disintegration/imaging"
)

func timespecToTime(ts syscall.Timespec) time.Time {   
	return time.Unix(int64(ts.Sec), int64(0))
}

type filesDataReturnType []map[string]string

func getFilesData() filesDataReturnType {
		outputDataSlice := make([]map[string]string, 1, 1)
		
		root := "static/data/images/"
		err := filepath.Walk(root, func(existingImageFilePath string, info os.FileInfo, err error) error {
			if existingImageFilePath != "static/data/images/" {
				outputFileInformation:= map[string]string{}
				stat_t := info.Sys().(*syscall.Stat_t)
				outputFileInformation["FileName"] = string(filepath.Base(existingImageFilePath))
				outputFileInformation["PathToThumbnail"] = "static/data/thumbnails/thumbnail_" + outputFileInformation["FileName"]
				outputFileInformation["UploadDate"] = (timespecToTime(stat_t.Ctim)).String()
				outputFileInformation["FileSize"] = strconv.FormatInt(stat_t.Size,10)
				fmt.Println("Obtaining EXIF data for " + existingImageFilePath)
				parser := exif.New()
				err := parser.Open(existingImageFilePath)
				if err != nil {
					fmt.Println("Error obtaining EXIF data for " + existingImageFilePath)
					fmt.Println(err)
					outputFileInformation["EXIFDateCreated"] = ""
					outputFileInformation["EXIFCameraMake"] = ""
					outputFileInformation["EXIFCameraModel"] = ""
				} 
				for k, v := range parser.Tags {					
					if k == "Date and Time (Digitized)" {
						outputFileInformation["EXIFDateCreated"] = v
					}
					if k == "Manufacturer" {
						outputFileInformation["EXIFCameraMake"] = v
					}
					if k == "Model" {
						outputFileInformation["EXIFCameraModel"] = v
					}
				}
				outputDataSlice = append(outputDataSlice, outputFileInformation)
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error opening running through existing files")
			fmt.Println(err)
		}

		
	return outputDataSlice
}

func diff(a, b uint32) int64 {
	if a > b {
		return int64(a - b)
	}
	return int64(b - a)
}

func checkImagesEqual(imageObject1 image.Image, imageObject2 image.Image) int {
	// many thanks to `https://www.rosettacode.org/wiki/Percentage_difference_between_images#Go` for best possible way to compare two images
	if imageObject1.ColorModel() != imageObject2.ColorModel() {
		fmt.Println("Images are different by color scheme")
		return 0
	}
	boundsImageObject1 := imageObject1.Bounds()
	if !boundsImageObject1.Eq(imageObject1.Bounds()) {
		fmt.Println("Images are different by sizes")
		return 0
	}
	var sum int64
	for y := boundsImageObject1.Min.Y; y < boundsImageObject1.Max.Y; y++ {
		for x := boundsImageObject1.Min.X; x < boundsImageObject1.Max.X; x++ {
			r1, g1, b1, _ := imageObject1.At(x, y).RGBA()
			r2, g2, b2, _ := imageObject2.At(x, y).RGBA()
			sum += diff(r1, r2)
			sum += diff(g1, g2)
			sum += diff(b1, b2)
		}
	}
	nPixels := (boundsImageObject1.Max.X - boundsImageObject1.Min.X) * (boundsImageObject1.Max.Y - boundsImageObject1.Min.Y)
	// At this point I need to accept a certain limitations as sometimes uploaded image somewhat different from the same image being loading
	// for example at test run I got Image difference of 0.213817% for exactly the same images, but were already uploaded and the other one was the source
	// Difference appeared due to many factors - poor test photo quality, changes applied by github.com/disintegration/imaging I can of course use os.Create ...
	// and other Go lang native ways to store the images, but I want image storage to be format independent and store it in the same format as uploaded image
	// while Go natively supports Encode function only as "jpeg.Encode(", "png.Encode" and not supports image.Encode, while I want exactly format independent,
	// github.com/disintegration/imaging allows me not to worry about the image formate while I store the image, but somehow changes the image -
	// - most probably, improves quality, thus resulting in 0.213817% difference
	imagesDifference := float64(sum*100)/(float64(nPixels)*0xffff*3)
	// therefore I think it should be fine to check if imagesDifference variable is not greater than 1 and if so consider the images equal
	// if imagesDifference variable is greater than zero than images are different, by for example contract changes or gamma changes, etc.
	// i.e. if we have imagesDifference like 0.2xxxxx% then consider that images are equal if imagesDifference is like 1.xxxxxxx% then consider images different
	if (imagesDifference < 1) {
		fmt.Printf("Images difference is %f%% considering equal\n", imagesDifference)
		return 1
	}
	return 0
}

func checkExists(imageObject image.Image) int {
	var files []string
	exists := 0
	root := "static/data/images/"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		fmt.Println("Error opening running through existing files")
		fmt.Println(err)
	}
	for _, existingImageFilePath := range files[1:] {
		existingImageFile, err := os.Open(existingImageFilePath)
		fmt.Printf("Comparing vs existing %+v\n", existingImageFilePath)
		if err != nil {
			fmt.Println("Error opening " + existingImageFilePath)
			fmt.Println(err)
		}
		defer existingImageFile.Close()
		existingImageData, _, err := image.Decode(existingImageFile)
		if err != nil {
			fmt.Println("Error creating image object for " + existingImageFilePath)
			fmt.Println(err)
		}
		equal := checkImagesEqual(imageObject, existingImageData)
		if (equal == 1) {
			exists = 1
			break
		}
	}
	return exists
}

func createThubmnail(imageObject image.Image, customFileName string) int {
	thumbnail := imaging.Thumbnail(imageObject, 80, 80, imaging.CatmullRom)
	imaging.Save(thumbnail, "static/data/thumbnails/thumbnail_" + customFileName)
	return 0
}

func uploadFile(fileLink io.ReadSeeker, customFileName string) int {
	imageObject, _, err := image.Decode(fileLink)
	fileLink.Seek(0,0)
	if err != nil {
		fmt.Println("Error creating image object")
		fmt.Println(err)
		return 1
	} else {
		exists := checkExists(imageObject)
		if (exists == 1) {
			return 2
		}
		// Copy the file as it was to preserve the EXIF data
		destinationImage, err := os.Create("static/data/images/" + customFileName)
		if err != nil {
			fmt.Println(err)
		}
		defer destinationImage.Close()
		_, err = io.Copy(destinationImage, fileLink)
		if err != nil {
			fmt.Println(err)
		}
		defer fileLink.Seek(0,0)
		// Creating thumbnail using the github.com module to demonstrate ability to use those, if needed
		createThubmnail(imageObject,  customFileName) 
	}
	return 0
}
