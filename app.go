package main

import (
	"html/template"
	"fmt"
	"net/http"
	"log"
	"path/filepath"
	"strconv"
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
	response := uploadFile(fileLink , customFileName)

	http.Redirect(w, r, "/?msg=" + strconv.Itoa(response), 301)
}

func deleteImageHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	fileName := r.Form.Get("delete")
	response := deleteImages(fileName)
	http.Redirect(w, r, "/?msg=" + strconv.Itoa(response), 301)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexPageTemplate:= template.Must(template.ParseFiles("templates/index.html", "templates/base.html"))

	filesData := getFilesData()
	msg := ""
	for k, v := range r.URL.Query() {
		if k == "msg" {
			msg = errorsInterpreter(v[0])
			break
		}
	}
	data := map[string]interface{}{
		"message": msg,
		"filesData": filesData,
	}
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
	http.HandleFunc("/deleteImage", deleteImageHandler)
	http.HandleFunc("/readme", readmeHandler)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
