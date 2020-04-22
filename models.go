package main

import (
	"io/ioutil"
	"fmt"
	"github.com/disintegration/imaging"
)

func createThubmnail(pathToImageUploaded string, customFileName string) int {
	uploadedImage := pathToImageUploaded  + customFileName
	uploadedImageLink, err := imaging.Open(uploadedImage)
	if err != nil {
		fmt.Println("Error creating thumbnail")
		fmt.Println(err)
	}
	thumbnail := imaging.Thumbnail(uploadedImageLink, 80, 80, imaging.CatmullRom)
	imaging.Save(thumbnail, "static/data/thumbnails/thumbnail_" + customFileName)
	return 0
}
//fileLink io.Reader for now not needed {name type} reserved for possible future use
func uploadFile(newFileContents []uint8, customFileName string) int {	
	ioutil.WriteFile("static/data/images/" + customFileName, newFileContents, 0644)
	createThubmnail("static/data/images/",  customFileName) 
	return 0
}
