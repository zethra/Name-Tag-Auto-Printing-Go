package main

import (

	"os"
	"fmt"
	"strings"
	"net/http"
	"html/template"
	"log"

	"ntap/data"
	"ntap/config"
	"ntap/service"

	"github.com/kardianos/osext"
	"github.com/gorilla/schema"
	"path"
	"sync"
	"os/signal"
	"syscall"
//	"strconv"
)

var configImpl config.Config
var nameTagQueue = data.NewNameTagQueue()
var printerQueue = data.NewPrinterQueue()

func main() {
	//Get Bin directory path
	filename, _ := osext.Executable()
	var pos int
	if pos = strings.LastIndex(filename, "\\"); pos == -1 {
		if pos = strings.LastIndex(filename, "/"); pos == -1 {
			fmt.Println("Cannot find base directory")
			return
		}
	}
	binDirectory := filename[0:pos] + "/"
	//	fmt.Println(binDirectory)

	//Make Config
	configImpl.PrintersFile = binDirectory + "../config/printer.xml"
	configImpl.QueueFile = binDirectory + "../config/queue.xml"
	configImpl.ImagesDirectory = binDirectory + "../web/static/images"
	configImpl.ScadDirectory = binDirectory + "../data/scad"
	configImpl.StlDirectory = binDirectory + "../data/stl"
	configImpl.GcodeDirectory = binDirectory + "../data/gcode"
	configDirectory := binDirectory + "../config"

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
	if _, err := os.Stat(configDirectory); os.IsNotExist(err) {
		fmt.Println("Makeing Config Directory")
		os.MkdirAll(configDirectory, 666)
	}
	if _, err := os.Stat(configImpl.PrintersFile); os.IsNotExist(err) {
		fmt.Println("Makeing Printers File")
		os.Create(configImpl.PrintersFile)
	}
	if _, err := os.Stat(configImpl.QueueFile); os.IsNotExist(err) {
		fmt.Println("Makeing Queue File")
		os.Create(configImpl.QueueFile)
	}

	printerQueue.Load(&configImpl)
	nameTagQueue.Load(&configImpl)
	//	printerQueue.Add(data.Printer{Name:"Test"}, &configImpl)
	//	nameTagQueue.Add(data.NameTag{Name:"Ben"}, &configImpl)

	//Start HTTP Server
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
		http.HandleFunc("/", serveTemplate)
		http.HandleFunc("/preview", preview)
		http.HandleFunc("/queue/add", addToQueue)
		http.HandleFunc("/manager/nameTagSubmit", nameTagSubmit)
		http.HandleFunc("/manager/printersSubmit", printersSubmit)
		http.ListenAndServe(":8080", nil)
	}()

	fmt.Println("HTTP Server Started")

	killchan := make(chan os.Signal, 2)
	signal.Notify(killchan, os.Interrupt, syscall.SIGTERM)
	// wait for kill signal
	<-killchan
	log.Println("Kill sig!")
	fmt.Println("Saving")
	printerQueue.Save(&configImpl)
	nameTagQueue.Save(&configImpl)
	os.Exit(0)
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
	//	fmt.Println(request.URL.Path)
	includesPath := path.Join("web", "dynamic", "includes.html")
	data := data.DataWrapper{}
	filePath := path.Join("web", "dynamic", request.URL.Path)
	if (request.URL.Path == "/") {
		filePath = path.Join("web", "dynamic", "index.html")
	} else if (request.URL.Path == "/manager") {
		//		fmt.Println("Loading manager")
		filePath = path.Join("web", "dynamic", "manager.html")
		data.PrinterQueue = printerQueue
		data.NameTagQueue = nameTagQueue
	} else if (request.URL.Path == "/manager/nameTags") {
		filePath = path.Join("web", "dynamic", "nameTags.html")
		data.NameTagQueue = nameTagQueue
	} else if (request.URL.Path == "/manager/printers") {
		filePath = path.Join("web", "dynamic", "printers.html")
		data.PrinterQueue = printerQueue
	}

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

	if err := tmpl.ExecuteTemplate(writer, "main", data); err != nil {
		log.Println(err.Error())
		http.Error(writer, http.StatusText(500), 500)
	}
}

func addToQueue(writer http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	if (name == "") {
		http.Error(writer, http.StatusText(400), 400)
		return
	}
	nameTag := data.NameTag{Name:name}
	nameTagQueue.Add(nameTag, &configImpl)
}

func nameTagSubmit(writer http.ResponseWriter, request *http.Request) {
	defer http.Redirect(writer, request, "/manager", 301)
	fmt.Println("Name Tags Submited")
	err := request.ParseForm()
	if err != nil {
		fmt.Println("Parsing form data failed:", err)
		http.Error(writer, http.StatusText(500), 500)
		return
	}
	wrapper := new(data.DataWrapper)
	decoder := schema.NewDecoder()
	err = decoder.Decode(wrapper, request.PostForm)
	if err != nil {
		fmt.Println("Decoding form data failed:", err)
		http.Error(writer, http.StatusText(400), 400)
		return
	}
	for i := 0; i < len(wrapper.NameTagQueue.Queue); i++ {
		if (len(wrapper.NameTagQueue.Queue) >= i + 1 && wrapper.NameTagQueue.Queue[i].Name != "") {
			if (len(wrapper.Delete) >= i + 1 && wrapper.Delete[i] == true) {
				nameTagQueue.Remove(wrapper.NameTagQueue.Queue[i].Id, &configImpl)
			} else {
				nameTagQueue.Queue[i] = wrapper.NameTagQueue.Queue[i]
			}
		}
	}
	nameTagQueue.Save(&configImpl)
	fmt.Println("Name Tags written")
}

func printersSubmit(writer http.ResponseWriter, request *http.Request)  {
	defer http.Redirect(writer, request, "/manager#printersTab", 301)
	fmt.Println("Printers Submited")
	err := request.ParseMultipartForm(0)
	if err != nil {
		fmt.Println("Parsing form data failed:", err)
		http.Error(writer, http.StatusText(500), 500)
		return
	}
	fmt.Printf("%v\n", request.MultipartForm.Value)
	wrapper := new(data.DataWrapper)
	decoder := schema.NewDecoder()
	err = decoder.Decode(wrapper, request.MultipartForm.Value)
	if err != nil {
		fmt.Println("Decoding form data failed:", err)
		http.Error(writer, http.StatusText(400), 400)
		return
	}
//	fmt.Printf("%v\n", wrapper.PrinterQueue.Queue)
	for i := 0; i < len(wrapper.PrinterQueue.Queue); i++ {
		if (len(wrapper.PrinterQueue.Queue) >= i + 1 && wrapper.PrinterQueue.Queue[i].Name != "") {
			if (len(wrapper.Delete) >= i + 1 && wrapper.Delete[i] == true) {
				printerQueue.Remove(wrapper.PrinterQueue.Queue[i].Id, &configImpl)
			} else {
				fmt.Printf("Setting printer: %s to positon: %d\n", wrapper.PrinterQueue.Queue[i].Name, i)
				fmt.Printf("%v\n", wrapper.PrinterQueue.Queue[i])
				printerQueue.Queue[i] = wrapper.PrinterQueue.Queue[i]
			}
		}
	}
	printerQueue.Save(&configImpl)
	fmt.Println("Printers written")
}