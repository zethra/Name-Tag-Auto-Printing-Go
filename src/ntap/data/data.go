package data
import (
	"fmt"
	"ntap/config"
	"github.com/satori/go.uuid"
	"encoding/xml"
	"io/ioutil"
	"errors"
)

type NameTag struct {
	Id                      uuid.UUID
	Name, Stl, Gcode, State string
	Printing, Error         bool
}

func (nameTag *NameTag) String() string {
	if (nameTag == nil) {
		return ""
	}
	return nameTag.Name
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
	fmt.Printf("Name tags: %v\n", queue.Queue)
	queue.Save(config)
}

func (queue *NameTagQueue) Update(nameTag NameTag, config *config.Config) error {
	err := errors.New("No name tag found")
	for i := 0; i < len(queue.Queue); i++ {
		if (uuid.Equal(queue.Queue[i].Id, nameTag.Id)) {
			fmt.Println("Updating name tag: ", queue.Queue[i].Name)
			nameTag.Id = queue.Queue[i].Id
			queue.Queue[i] = nameTag
			err = nil
		}
	}
	fmt.Printf("Name tags: %v\n", queue.Queue)
	queue.Save(config)
	return err
}

func (queue *NameTagQueue) Remove(id uuid.UUID, config *config.Config) error {
	err := errors.New("No name tag found")
	for i := 0; i < len(queue.Queue); i++ {
		if (uuid.Equal(queue.Queue[i].Id, id)) {
			fmt.Println("Removing name tag: ", queue.Queue[i].Name)
			queue.Queue = append(queue.Queue[:i], queue.Queue[i + 1:]...)
			err = nil
		}
	}
	fmt.Printf("Name tags: %v\n", queue.Queue)
	queue.Save(config)
	return err
}

func (queue *NameTagQueue) Find(id uuid.UUID, config *config.Config) (*NameTag, error) {
	for i := 0; i < len(queue.Queue); i++ {
		if (uuid.Equal(queue.Queue[i].Id, id)) {
			return &queue.Queue[i], nil
		}
	}
	return &NameTag{}, errors.New("No name tag found\n")
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
	Id                                  uuid.UUID
	Name, Ip, ApiKey, ConfigFile, State string
	Port                                int
	Active, Printing                    bool
	NameTag                             *NameTag
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
	fmt.Printf("Printers: %v\n", queue.Queue)
	queue.Save(config)
}

func (queue *PrinterQueue) Update(printer Printer, config *config.Config) error {
	err := errors.New("No printer found")
	for i := 0; i < len(queue.Queue); i++ {
		if (uuid.Equal(queue.Queue[i].Id, printer.Id)) {
			fmt.Printf("Updating printer: %s\n", queue.Queue[i].Name)
			printer.Id = queue.Queue[i].Id
			queue.Queue[i] = printer
			err = nil
		}
	}
	fmt.Printf("Printers: %v\n", queue.Queue)
	queue.Save(config)
	return err
}

func (queue *PrinterQueue) Remove(id uuid.UUID, config *config.Config) error {
	err := errors.New("No printer found")
	for i := 0; i < len(queue.Queue); i++ {
		if (uuid.Equal(queue.Queue[i].Id, id)) {
			fmt.Printf("Removing printer: %s\n", queue.Queue[i].Name)
			queue.Queue = append(queue.Queue[:i], queue.Queue[i + 1:]...)
			err = nil
		}
	}
	fmt.Printf("Printers: %v\n", queue.Queue)
	queue.Save(config)
	return err
}

func (queue *PrinterQueue) FindByIp(ip string, config *config.Config) (*Printer, error) {
	for i := 0; i < len(queue.Queue); i++ {
		if (queue.Queue[i].Ip == ip) {
			return &queue.Queue[i], nil
		}
	}
	return &Printer{}, errors.New("No Printer found")
}

func (queue *PrinterQueue) GetNext() *Printer {
	for i := 0; i < len(queue.Queue); i++ {
		if (queue.Queue[i].Active == true && queue.Queue[i].Printing == false) {
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