package main

import (
	"html/template"
	"io/ioutil"
	"fmt"
	"net/http"
	"log"
	"path/filepath"
)

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	customFileName := r.Form.Get("fileName")
	fileLink, fileObject, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer fileLink.Close()
	if customFileName != "" {
		fileExtension := filepath.Ext(fileObject.Filename)
		customFileName = customFileName + fileExtension
	} else {
		customFileName = fileObject.Filename
	}
	fmt.Printf("Name: %+v\n", customFileName)
	fmt.Printf("Size: %+v\n", fileObject.Size)
	fmt.Printf("MIME: %+v\n", fileObject.Header)
	
	newFileContents, err := ioutil.ReadAll(fileLink)
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile("static/data/images/" + customFileName, newFileContents, 0644)
	
	http.Redirect(w, r, "/", 301)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexPageTemplate:= template.Must(template.ParseFiles("templates/index.html", "templates/base.html"))
	var data []int
	err := indexPageTemplate.ExecuteTemplate(w, "base", data)
	if err != nil {		
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/uploadFile", uploadFileHandler)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
