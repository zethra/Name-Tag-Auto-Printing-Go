package service
import (
	"ntap/config"
	"sync"
	"fmt"
"strings"
	"os/exec"
	"encoding/json"
	"os"
)

type jsonResponse struct {
	Image, Error string
	Code int
}

func Preview(name string, config *config.Config) ([]byte, int) {
	var response jsonResponse

	if(name == "") {
		fmt.Println("No name previded")
		response = jsonResponse{Code:-1, Error:"No name previded"}
		json, err := json.Marshal(response)
		if(err != nil) {
			return make([]byte, 0), 500
		}
		return json, 400
	}

	image := config.ImagesDirectory + "/" + name + ".png"

	if _, err := os.Stat(image); err == nil {
		fmt.Println("Image already exstits")
		response = jsonResponse{Code:0, Image: "/static/images/" + name + ".png"}
		json, err := json.Marshal(response)
		if(err != nil) {
			return make([]byte, 0), 500
		}
		return json, 200
	}

	zoom := 100
	if(len(name) > 8) {
		zoom = 150
	} else if(len(name) > 5) {
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
	if(err != nil) {
		return make([]byte, 0), 500
	}
	return json, 200
}

func exe_cmd(cmd string, wg *sync.WaitGroup) {
	fmt.Println("command is: ",cmd, "\n")
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head,parts...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
	wg.Done() // Need to signal to waitgroup that this goroutine is done
}