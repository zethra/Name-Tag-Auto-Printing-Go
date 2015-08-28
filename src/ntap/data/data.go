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

func (queue *nameTagQueue) Add(nameTag NameTag) {
	queue.Queue = append(queue.Queue, nameTag)
}

func (queue *nameTagQueue) Remove(id uuid.UUID) {
	for i := 0; i < len(queue.Queue); i++ {
		if (uuid.Equal(queue.Queue[i].Id, id)) {
			queue.Queue = append(queue.Queue[:i], queue.Queue[i + 1:]...)
		}
	}
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
	fmt.Println(string(xml))
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
	Name, Ip, apiKey string
	Port             byte
	Active, Printing bool
	NameTag          *NameTag
}

func (printer *Printer)Slice(config *config.Config) {
	fmt.Println("Sliced")
}

type printerQueue struct {
	queue []Printer
}

func NewPrinterQueue() printerQueue {
	queue := printerQueue{queue:make([]Printer, 0)}
	return queue
}

func (queue *printerQueue) Add(printer Printer) {
	queue.queue = append(queue.queue, printer)
}

func (queue *printerQueue) Remove(id uuid.UUID) {
	for i := 0; i < len(queue.queue); i++ {
		if (uuid.Equal(queue.queue[i].Id, id)) {
			queue.queue = append(queue.queue[:i], queue.queue[i + 1:]...)
		}
	}
}

func (queue *printerQueue)Get(i int) *Printer {
	return &queue.queue[i]
}

func (queue *printerQueue) GetNext() (Printer, error) {
	for i := 0; i < len(queue.queue); i++ {
		if (queue.queue[i].Printing == false) {
			return queue.queue[i], nil
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
	fmt.Println(string(xml))
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