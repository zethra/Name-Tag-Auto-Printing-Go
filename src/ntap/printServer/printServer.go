package printServer

import (
	"ntap/data"
	"ntap/config"
	"time"
	"fmt"
	"ntap/service"
	"log"
)

var Quit chan struct{}
var Timer *time.Ticker

func Start(interval time.Duration, nameTagQueue *data.NameTagQueue, printerQueue *data.PrinterQueue, config *config.Config) {
	Quit = make(chan struct{})
	go run(interval, nameTagQueue, printerQueue, config)
}

func Stop() {
	close(Quit)
}

func run(interval time.Duration, nameTagQueue *data.NameTagQueue,
printerQueue *data.PrinterQueue, config *config.Config) {
	fmt.Println("Print Server Started")
	Timer = time.NewTicker(interval * time.Second)
	for {
		select {
		case <-Timer.C:
			printer := printerQueue.GetNext()
			if (printer == nil) {
				continue
			}
			var nameTag *data.NameTag
			getNew:
			if (printer.NameTag != nil && printer.NameTag.Error == false) {
				fmt.Printf("%v\n", printer.NameTag)
				tag, err := nameTagQueue.Find(printer.NameTag.Id, config)
				if (err != nil) {
					printer.NameTag.Error = true
					goto getNew
				}
				nameTag = tag
			} else {
				nameTag = nameTagQueue.GetNext()
				if (nameTag == nil) {
					continue
				}
				printer.NameTag = nameTag
			}
			nameTag.State = "Assigned to printer"
			if (nameTag.Stl == "") {
				nameTag.State = "Rendering STL"
				err := service.Export(nameTag, config)
				if (err != nil) {
					fmt.Println(err)
					nameTag.Error = true
					nameTag.State = "Has Error"
					goto save
				}
			}
			if (nameTag.Gcode == "") {
				nameTag.State = "Slicing"
				err := service.Slice(nameTag, printer, config)
				if (err != nil) {
					fmt.Println(err)
					nameTag.Error = true
					nameTag.State = "Has Error"
					goto save
				}
			}
			if (nameTag.Printing == false && printer.Printing == false) {
				nameTag.State = "Uploading"
				err := service.Upload(nameTag, printer, config)
				if (err != nil) {
					log.Println(err)
					printer.Active = false
					goto save
				}
				printer.Printing = true
				nameTag.Printing = true
				nameTag.State = "Printing"
			}
			save:
			printerQueue.Save(config)
			nameTagQueue.Save(config)
			config.DebugLog("Finnished Loop")
		case <-Quit:
			Timer.Stop()
			return
		}
	}

}