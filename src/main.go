package main

import (

	"os"
	"fmt"
	"strings"
	"net/http"
	"html/template"
	"log"

//	"ntap/data"
	"ntap/config"
	"ntap/service"

	"github.com/kardianos/osext"
	"path"
)

var configImpl config.Config

func main() {

	//Get Bim directory path
	filename, _ := osext.Executable()
	var pos int
	if pos = strings.LastIndex(filename, "\\"); pos == -1 {
		if pos = strings.LastIndex(filename, "/"); pos == -1 {
			fmt.Println("Cannot find base directory")
			return
		}
	}
	binDirectory := filename[0:pos] + "/"
	fmt.Println(binDirectory)

	//Make Config
	configImpl.PrintersFile = binDirectory + "../config/printer.xml"
	configImpl.QueueFile = binDirectory + "../config/queue.xml"
	configImpl.ImagesDirectory = binDirectory + "../web/static/images"
	configImpl.ScadDirectory = binDirectory + "../data/scad"
	configImpl.StlDirectory = binDirectory + "../data/stl"
	configImpl.GcodeDirectory = binDirectory + "../data/gcode"

	//Generate Files
	if _, err := os.Stat(configImpl.ImagesDirectory); os.IsNotExist(err) {
		fmt.Println("Makeing Images Directory")
		os.MkdirAll(configImpl.ImagesDirectory, 666)
	}
	if _, err := os.Stat(configImpl.ScadDirectory); os.IsNotExist(err) {
		fmt.Println("Makeing Scad Directory")
		os.MkdirAll(configImpl.ScadDirectory, 666)
	}
	if _, err := os.Stat(configImpl.StlDirectory); os.IsNotExist(err) {
		fmt.Println("Makeing STL Directory")
		os.MkdirAll(configImpl.StlDirectory, 666)
	}
	if _, err := os.Stat(configImpl.GcodeDirectory); os.IsNotExist(err) {
		fmt.Println("Makeing Gcode Directory")
		os.MkdirAll(configImpl.GcodeDirectory, 666)
	}
	if _, err := os.Stat(configImpl.PrintersFile); os.IsNotExist(err) {
		fmt.Println("Makeing Printers File")
		os.Create(configImpl.PrintersFile)
	}
	if _, err := os.Stat(configImpl.QueueFile); os.IsNotExist(err) {
		fmt.Println("Makeing Queue File")
		os.Create(configImpl.QueueFile)
	}

//	nameTag := data.NameTag{Name:"Ben"}
//	nameTag.Export(&configImpl)

	//Start HTTP Server
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.HandleFunc("/", serveTemplate)
	http.HandleFunc("/preview", preview)
	http.ListenAndServe(":8080", nil)
}

func preview(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Recieved preview request")
	output, code := service.Preview(request.FormValue("name"), &configImpl)
	fmt.Println("Server output: " + string(output))
	if (code != 200) {
		http.Error(writer, string(output), code)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(output)
}

func serveTemplate(writer http.ResponseWriter, request *http.Request) {
	includesPath := path.Join("web", "dynamic", "includes.html")
	filePath := path.Join("web", "dynamic", request.URL.Path)
	if (request.URL.Path == "/") {
		filePath = path.Join("web", "dynamic", "index.html")
	}

//	fmt.Println(filePath)

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(writer, request)
			return
		}
	}
	if info.IsDir() {
		http.NotFound(writer, request)
		return
	}
	tmpl, err := template.ParseFiles(includesPath, filePath)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(writer, http.StatusText(500), 500)
		return
	}

	if err := tmpl.ExecuteTemplate(writer, "main", nil); err != nil {
		log.Println(err.Error())
		http.Error(writer, http.StatusText(500), 500)
	}
}