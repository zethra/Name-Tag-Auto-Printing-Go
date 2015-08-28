package data
import (
	"fmt"
	"ntap/config"
)

type NameTag struct {
	Id               int
	Name, Stl, Gcode string
	Printer          *Printer
	Printing         bool
}

func (nameTag *NameTag) Export(config     *config.Config) {
	fmt.Println("Exported")
}

type nameTagQueue struct {
	queue []NameTag
}

func NewNameTagQueue() nameTagQueue {
	queue := nameTagQueue{queue:make([]NameTag, 0)}
	return queue
}

func (queue *nameTagQueue) add(nameTag NameTag) {
	queue.queue = append(queue.queue, nameTag)
}

func (queue *nameTagQueue) remove(id int) {
	for i := 0; i < len(queue.queue); i++ {
		if (queue.queue[i].Id == id) {
			queue.queue = append(queue.queue[:i], queue.queue[i + 1]...)
		}
	}
}

func (queue *nameTagQueue) getNext() NameTag{
	for i := 0; i < len(queue.queue); i++ {
		if (queue.queue[i].Printer == nil && (queue.queue[i].Stl == "" ||
				queue.queue[i].Gcode == "" || queue.queue[i].Printing == false)) {
			return queue.queue[i]
		}
	}
	return nil
}

type Printer struct {
	Name, Ip, apiKey string
	Port             byte
	Active, Printing bool
	NameTag          *NameTag
}

func (printer *Printer)slice(config *config.Config) {
	fmt.Println("Sliced")
}

