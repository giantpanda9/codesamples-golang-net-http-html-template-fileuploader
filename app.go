package main

import (
	"html/template"
	"fmt"
	"net/http"
	"log"
	"path/filepath"

)

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	customFileName := r.Form.Get("fileName")
	fileLink, fileHeader, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer fileLink.Close()
	if customFileName != "" {
		fileExtension := filepath.Ext(fileHeader.Filename)
		customFileName = customFileName + fileExtension
	} else {
		customFileName = fileHeader.Filename
	}
	fmt.Printf("Name: %+v\n", customFileName)
	fmt.Printf("Size: %+v\n", fileHeader.Size)
	fmt.Printf("MIME: %+v\n", fileHeader.Header)

	response := uploadFile(fileLink , customFileName)
	fmt.Printf("Upload file response: %+v\n", response)
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

func readmeHandler(w http.ResponseWriter, r *http.Request) {
	readmePageTemplate:= template.Must(template.ParseFiles("templates/instructions.html", "templates/base.html"))
	var data []int
	err := readmePageTemplate.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/uploadFile", uploadFileHandler)
	http.HandleFunc("/readme", readmeHandler)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
