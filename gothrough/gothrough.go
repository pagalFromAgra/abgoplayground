package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/wearkinetic/stater/staterlogic/sdata"
)

var currTime float64
var npoints int64
var speed float64

type Message struct {
	Action string       `json:"action"`
	Body   sdata.Record `json:"body"`
}

func main() {

	fpPtr := flag.String("fp", "", "File/directory path")
	speedPtr := flag.Float64("speed", 1, "Specify speed")
	flag.Parse()

	speed = *speedPtr
	filepath := *fpPtr

	// Check if it's a file or a directory
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		log.Fatalf("Coulndn't get source of raw data")
	}

	npoints = 0

	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(filepath)
		if err != nil {
			log.Fatalf("Error reading files from %s", filepath)
		}

		// files ordered by the names
		for _, file := range files {
			// 6 - Read stdin, that will
			gothrough(path.Join(filepath, file.Name()))

		}
	} else {
		gothrough(filepath)
	}
}

func gothrough(filepath string) {

	log.Println("Opening", filepath)
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Error opening file %s", filepath)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line, t, err := Transform(scanner.Bytes())
		if err != nil {
			continue
		}
		npoints++
		if npoints > 1 {
			time.Sleep(time.Duration((t-currTime)/speed) * time.Millisecond)
		}
		currTime = t

		fmt.Println(string(line))

	}

}

func NewMessage(record sdata.Record) *Message {
	return &Message{
		Action: "i2c.point",
		Body:   record,
	}
}

func Transform(b []byte) ([]byte, float64, error) {

	// Unmarshal into a record
	var record sdata.Record
	err := json.Unmarshal(b, &record)
	if err != nil {
		return nil, currTime, err
	}

	out, err := json.Marshal(NewMessage(record))
	if err != nil {
		return nil, currTime, err
	}

	return out, record.Time, err
}
