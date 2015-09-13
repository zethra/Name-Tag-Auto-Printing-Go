package data
import (
	"fmt"
	"ntap/config"
	"github.com/satori/go.uuid"
	"encoding/xml"
	"io/ioutil"
)

type NameTag struct {
	Id               uuid.UUID
	Name, Stl, Gcode string
	Printing, Error  bool
}

func (nameTag *NameTag) String() string {
	if (nameTag == nil) {
		return ""
	}
	return nameTag.Name
}

func (nameTag *NameTag) Export(config     *config.Config) {
	fmt.Println("Exported")
}

type NameTagQueue struct {
	Queue []NameTag `xml:"NameTag", json:"NamneTag"`
}

func NewNameTagQueue() NameTagQueue {
	queue := NameTagQueue{Queue:make([]NameTag, 0)}
	return queue
}

func (queue *NameTagQueue) Add(nameTag NameTag, config *config.Config) {
	fmt.Println("Adding name tag: ", nameTag.Name)
	queue.Queue = append(queue.Queue, nameTag)
	fmt.Printf("Name tags: %v\n", queue)
	queue.Save(config)
}

func (queue *NameTagQueue) Remove(id uuid.UUID, config *config.Config) {
	for i := 0; i < len(queue.Queue); i++ {
		if (uuid.Equal(queue.Queue[i].Id, id)) {
			fmt.Println("Removing name tag: ", queue.Queue[i].Name)
			queue.Queue = append(queue.Queue[:i], queue.Queue[i + 1:]...)
		}
	}
	fmt.Printf("Name tags: %v\n", queue)
	queue.Save(config)
}

func (queue *NameTagQueue) GetNext() *NameTag {
	for i := 0; i < len(queue.Queue); i++ {
		if (queue.Queue[i].Error == false && queue.Queue[i].Stl == "" || queue.Queue[i].Gcode == "" ||
		queue.Queue[i].Printing == false) {
			return &queue.Queue[i]
		}
	}
	return nil
}

func (queue *NameTagQueue) Save(config *config.Config) {
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

func (queue *NameTagQueue) Load(config *config.Config) {
	data, err := ioutil.ReadFile(config.QueueFile)
	if (err != nil) {
		panic(err)
		return
	}
	if (string(data) == "") {
		return
	}
	err = xml.Unmarshal(data, queue)
	if (err != nil) {
		panic(err)
		return
	}
}


type Printer struct {
	Id                           uuid.UUID
	Name, Ip, ApiKey, ConfigFile string
	Port                         int
	Active, Printing             bool
	NameTag                      *NameTag
}

func (printer *Printer)String() string {
	if (printer == nil) {
		return ""
	}
	return printer.Name
}

func (printer *Printer) Slice(config *config.Config) {
	fmt.Println("Sliced")
}

type PrinterQueue struct {
	Queue []Printer `xml:"Printer", json:"Printer"`
}

func NewPrinterQueue() PrinterQueue {
	queue := PrinterQueue{Queue:make([]Printer, 0)}
	return queue
}

func (queue *PrinterQueue) Add(printer Printer, config *config.Config) {
	fmt.Printf("Adding printer: %s\n", printer.Name)
	queue.Queue = append(queue.Queue, printer)
	fmt.Printf("Name tags: %v\n", queue.Queue)
	queue.Save(config)
}

func (queue *PrinterQueue) Remove(id uuid.UUID, config *config.Config) {
	for i := 0; i < len(queue.Queue); i++ {
		if (uuid.Equal(queue.Queue[i].Id, id)) {
			fmt.Printf("Removing printer: %s\n", queue.Queue[i].Name)
			queue.Queue = append(queue.Queue[:i], queue.Queue[i + 1:]...)
		}
	}
	fmt.Printf("Name tags: %v\n", queue.Queue)
	queue.Save(config)
}

func (queue *PrinterQueue) GetNext() *Printer {
	for i := 0; i < len(queue.Queue); i++ {
		if (queue.Queue[i].Printing == false) {
			return &queue.Queue[i]
		}
	}
	return nil
}

func (queue *PrinterQueue) Save(config *config.Config) {
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

func (queue *PrinterQueue) Load(config *config.Config) {
	data, err := ioutil.ReadFile(config.PrintersFile)
	if (err != nil) {
		panic(err)
		return
	}
	if (string(data) == "") {
		return
	}
	err = xml.Unmarshal(data, queue)
	if (err != nil) {
		panic(err)
		return
	}
}

type DataWrapper struct {
	NameTagQueue NameTagQueue
	PrinterQueue PrinterQueue
	Delete       []bool
}