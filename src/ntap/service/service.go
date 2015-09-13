package service
import (
	"ntap/config"
	"sync"
	"fmt"
	"strings"
	"os/exec"
	"encoding/json"
	"os"
	"ntap/data"
	"errors"
	"bytes"
	"net/http"
	"mime/multipart"
	"io/ioutil"
	"net/http/httputil"
	"log"
)

type jsonResponse struct {
	Image, Error string
	Code         int
}

func Preview(name string, config *config.Config) ([]byte, int) {
	var response jsonResponse

	if (name == "") {
		fmt.Println("No name previded")
		response = jsonResponse{Code:-1, Error:"No name previded"}
		json, err := json.Marshal(response)
		if (err != nil) {
			return make([]byte, 0), 500
		}
		return json, 400
	}

	image := config.ImagesDirectory + "/" + name + ".png"

	if _, err := os.Stat(image); err == nil {
		fmt.Println("Image already exstits")
		response = jsonResponse{Code:0, Image: "/static/images/" + name + ".png"}
		json, err := json.Marshal(response)
		if (err != nil) {
			return make([]byte, 0), 500
		}
		return json, 200
	}

	zoom := 100
	if (len(name) > 8) {
		zoom = 150
	} else if (len(name) > 5) {
		zoom = 130
	}

	pngargs := fmt.Sprintf(" -o %s/%s.png -D name=\"%s\" -D chars=%d --camera=0,0,0,0,0,0,%d %s/name.scad",
		config.ImagesDirectory, name, name, len(name), zoom, config.ScadDirectory)

	fmt.Println("Generating preview...")
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go exe_cmd("openscad" + pngargs, wg)
	wg.Wait()
	fmt.Println("Done generating preview")

	response = jsonResponse{Code:0, Image: "/static/images/" + name + ".png"}
	json, err := json.Marshal(response)
	if (err != nil) {
		return make([]byte, 0), 500
	}
	return json, 200
}

func Export(nameTag *data.NameTag, config *config.Config) error {
	if (nameTag.Name == "") {
		return errors.New("Name tag has no name set")
	}
	stlArgs := fmt.Sprintf(" -o %s/%s.stl -D name=\"%s\" -D chars=%d %s/name.scad", config.StlDirectory,
		nameTag.Name, nameTag.Name, len(nameTag.Name), config.ScadDirectory)

	fmt.Println("Exporting STL...")
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go exe_cmd("openscad" + stlArgs, wg)
	wg.Wait()
	stl := config.StlDirectory + "/" + nameTag.Name + ".stl"
	if _, err := os.Stat(stl); os.IsNotExist(err) {
		return errors.New("An Error occured while exporting STL")
	}
	nameTag.Stl = nameTag.Name + ".stl"
	fmt.Println("Done exporting STL")
	return nil
}

func Slice(nameTag *data.NameTag, printer *data.Printer, config *config.Config) error {
	if (nameTag.Name == "") {
		return errors.New("Name tag has no name set")
	}
	var configFile string
	if (printer.ConfigFile == "") {
		configFile = config.DefaultConfig
	} else {
		configFile = printer.ConfigFile
	}

	slic3rArgs := fmt.Sprintf(" %s/%s.stl --output %s/%s.gcode --load %s", config.StlDirectory,
		nameTag.Name, config.GcodeDirectory, nameTag.Name, configFile)

	fmt.Println("Slicing STL...")
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go exe_cmd("slic3r" + slic3rArgs, wg)
	wg.Wait()
	gcode := config.GcodeDirectory + "/" + nameTag.Name + ".gcode"
	if _, err := os.Stat(gcode); os.IsNotExist(err) {
		return errors.New("An Error occured while slicing STL")
	}
	nameTag.Gcode = nameTag.Name + ".gcode"
	fmt.Println("Done slicing STL")
	return nil
}

func Upload(nameTag *data.NameTag, printer *data.Printer, config *config.Config) error {
	fmt.Println("Uploading Gcode...")
	uri := fmt.Sprintf("http://%s:%d/api/files/local", printer.Ip, printer.Port)
	file, err := os.Open(config.GcodeDirectory + "/" + nameTag.Gcode)
	if err != nil {
		return err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fi, err := file.Stat()
	if err != nil {
		return err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		return err
	}
	part.Write(fileContents)
	_ = writer.WriteField("print", "true")
	err = writer.Close()
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return err
	}
	request.Header.Add("X-Api-Key", printer.ApiKey)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	data, err := httputil.DumpRequest(request, true)
	if err == nil {
		config.DebugLog(data)
	} else {
		config.DebugLog(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	body = &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		log.Println(err)
	}
	resp.Body.Close()
	fmt.Printf("Upload finnished: %d\n", resp.StatusCode)
	config.DebugLog(body)
	return nil
}

func Delete(nameTag data.NameTag, printer data.Printer, config *config.Config) {
	fmt.Printf("Sending deleting request for %s to printer %s", nameTag, printer)
	uri := fmt.Sprintf("http://%s:%d/api/files/local", printer.Ip, printer.Port)
	request, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Add("X-Api-Key", printer.ApiKey)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Delete finnished: %d\n", resp.StatusCode)

}

func exe_cmd(cmd string, wg *sync.WaitGroup) {
	fmt.Println("command is: ", cmd, "\n")
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%s\n", out)
	wg.Done() // Need to signal to waitgroup that this goroutine is done
}

