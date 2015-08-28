package data
import (
	"fmt"
	"ntap/config"
	"github.com/satori/go.uuid"
	"encoding/xml"
	"errors"
	"io/ioutil"
)

type NameTag struct {
	Id               uuid.UUID
	Name, Stl, Gcode string
	Printer          *Printer
	Printing         bool
}

func (nameTag *NameTag) String() string {
	if(nameTag == nil) {
		return ""
	}
	return nameTag.Name
}

func (nameTag *NameTag) Export(config     *config.Config) {
	fmt.Println("Exported")
}

type nameTagQueue struct {
	Queue []NameTag `xml:"NameTag", json:"NamneTag"`
}

func NewNameTagQueue() nameTagQueue {
	queue := nameTagQueue{Queue:make([]NameTag, 0)}
	return queue
}

func (queue *nameTagQueue) Add(nameTag NameTag, config *config.Config) {
	queue.Queue = append(queue.Queue, nameTag)
	queue.Save(config)
}

func (queue *nameTagQueue) Remove(id uuid.UUID, config *config.Config) {
	for i := 0; i < len(queue.Queue); i++ {
		if (uuid.Equal(queue.Queue[i].Id, id)) {
			queue.Queue = append(queue.Queue[:i], queue.Queue[i + 1:]...)
		}
	}
	queue.Save(config)
}

func (queue *nameTagQueue)Get(i int) *NameTag {
	return &queue.Queue[i]
}

func (queue *nameTagQueue) GetNext() (NameTag, error) {
	for i := 0; i < len(queue.Queue); i++ {
		if (queue.Queue[i].Stl == "" || queue.Queue[i].Gcode == "" || queue.Queue[i].Printing == false) {
			return queue.Queue[i], nil
		}
	}
	return NameTag{}, errors.New("No names tags")
}

func (queue *nameTagQueue) Save(config *config.Config) {
	xml, err := xml.MarshalIndent(queue, "", "    ")
	if (err != nil) {
		panic(err)
		return
	}
//	fmt.Println(string(xml))
	err = ioutil.WriteFile(config.QueueFile, xml, 666)
	if (err != nil) {
		panic(err)
		return
	}
}

func (queue *nameTagQueue) Load(config *config.Config) {
	data, err := ioutil.ReadFile(config.QueueFile)
	if (err != nil) {
		panic(err)
		return
	}
	if(string(data) == "") {
		return
	}
	err = xml.Unmarshal(data, queue)
	if (err != nil) {
		panic(err)
		return
	}
}


type Printer struct {
	Id               uuid.UUID
	Name, Ip, ApiKey string
	Port             byte
	Active, Printing bool
	NameTag          *NameTag
}

func (printer *Printer)String() string {
	if(printer == nil) {
		return ""
	}
	return printer.Name
}

func (printer *Printer) Slice(config *config.Config) {
	fmt.Println("Sliced")
}

type printerQueue struct {
	Queue []Printer `xml:"Printer", json:"Printer"`
}

func NewPrinterQueue() printerQueue {
	queue := printerQueue{Queue:make([]Printer, 0)}
	return queue
}

func (queue *printerQueue) Add(printer Printer, config *config.Config) {
	queue.Queue = append(queue.Queue, printer)
	queue.Save(config)
}

func (queue *printerQueue) Remove(id uuid.UUID, config *config.Config) {
	for i := 0; i < len(queue.Queue); i++ {
		if (uuid.Equal(queue.Queue[i].Id, id)) {
			queue.Queue = append(queue.Queue[:i], queue.Queue[i + 1:]...)
		}
	}
	queue.Save(config)
}

func (queue *printerQueue)Get(i int) *Printer {
	return &queue.Queue[i]
}

func (queue *printerQueue) GetNext() (Printer, error) {
	for i := 0; i < len(queue.Queue); i++ {
		if (queue.Queue[i].Printing == false) {
			return queue.Queue[i], nil
		}
	}
	return Printer{}, errors.New("No printers")
}

func (queue *printerQueue) Save(config *config.Config) {
	xml, err := xml.MarshalIndent(queue, "", "    ")
	if (err != nil) {
		panic(err)
		return
	}
//	fmt.Println(string(xml))
	err = ioutil.WriteFile(config.PrintersFile, xml, 666)
	if (err != nil) {
		panic(err)
		return
	}
}

func (queue *printerQueue) Load(config *config.Config) {
	data, err := ioutil.ReadFile(config.PrintersFile)
	if (err != nil) {
		panic(err)
		return
	}
	if(string(data) == "") {
		return
	}
	err = xml.Unmarshal(data, queue)
	if (err != nil) {
		panic(err)
		return
	}
}

type DataWrapper struct {
	NameTagQueue nameTagQueue
	PrinterQueue printerQueue
}