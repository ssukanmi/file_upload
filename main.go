package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	indexTemplate *template.Template
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := indexTemplate.Execute(w, nil); err != nil {
		panic(err)
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Uploading File\n")

	// parse input, type multilpart/from-data
	r.ParseMultipartForm(10 << 20)

	// retrieve file from posted form-data
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving file from form-date")
		fmt.Println(err)
		log.Fatal(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// write temporty file on our server
	tempFile, err := ioutil.TempFile("temp-images", "upload-*.png")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	tempFile.Write(fileBytes)

	//return whether or not this has been successful
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func setupRoutes() {
	var err error
	indexTemplate, err = template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Go file upload server started")
	setupRoutes()
}
