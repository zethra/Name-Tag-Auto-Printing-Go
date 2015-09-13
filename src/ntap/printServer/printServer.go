package printServer

import (
	"ntap/data"
	"ntap/config"
	"time"
	"fmt"
	"ntap/service"
	"log"
)

var Quit chan struct {}
var Timer time.Ticker

func Start(interval time.Duration, nameTagQueue *data.NameTagQueue, printerQueue *data.PrinterQueue, config *config.Config) {
	Quit = make(chan struct {})
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
			if (printer.NameTag != nil && printer.NameTag.Error == false) {
				nameTag = printer.NameTag
			} else {
				nameTag = nameTagQueue.GetNext()
				if (nameTag == nil) {
					continue
				}
				printer.NameTag = nameTag
			}
			if (nameTag.Stl == "") {
				err := service.Export(nameTag, config)
				if (err != nil) {
					fmt.Println(err)
					nameTag.Error = true
					goto save
				}
			}
			if (nameTag.Gcode == "") {
				err := service.Slice(nameTag, printer, config)
				if (err != nil) {
					fmt.Println(err)
					nameTag.Error = true
					goto save
				}
			}
			if(nameTag.Printing == false && printer.Printing == false) {
				err := service.Upload(nameTag, printer, config)
				if(err != nil) {
					log.Println(err)
					printer.Active = false
					goto save
				}
				printer.Printing = true
				nameTag.Printing = true
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