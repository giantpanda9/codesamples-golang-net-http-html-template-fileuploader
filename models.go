package main

import (
	//"io/ioutil"
	"fmt"
	"github.com/disintegration/imaging"
	"io"
	"image"

)



func createThubmnail(imageObject image.Image, customFileName string) int {
	thumbnail := imaging.Thumbnail(imageObject, 80, 80, imaging.CatmullRom)
	imaging.Save(thumbnail, "static/data/thumbnails/thumbnail_" + customFileName)
	return 0
}

func createImage(imageObject image.Image, customFileName string) int {
	err := imaging.Save(imageObject, "static/data/images/" + customFileName)
	if err != nil {
		fmt.Println(err)
	}
	return 0
}
//fileLink io.Reader for now not needed {name type} reserved for possible future use
func uploadFile(fileLink io.Reader, customFileName string) int {
	imageObject, _, err := image.Decode(fileLink)
	if err != nil {
		fmt.Println("Error creating image object")
		fmt.Println(err)
	}

	createImage(imageObject,  customFileName) 
	createThubmnail(imageObject,  customFileName) 
	return 0
}
