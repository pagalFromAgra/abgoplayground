package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"os"
	"time"
)

const (
	TIME_FORMAT = "2006-01-02T15:04:05.000Z"
)

func main() {

	csvdata := flag.String("csv", "", "a string")
	flag.Parse()

	this1, _ := time.Parse(TIME_FORMAT, "2017-08-24T07:50:01.000Z")
	that1, _ := time.Parse(TIME_FORMAT, "2017-09-12T18:35:19.000Z")
	tdiff1 := that1.Unix() - this1.Unix() // for 34d3n426LdxamFFN

	this2, _ := time.Parse(TIME_FORMAT, "2017-08-23T09:00:01.000Z")
	that2, _ := time.Parse(TIME_FORMAT, "2017-09-12T15:58:58.000Z")
	tdiff2 := that2.Unix() - this2.Unix() // for pZV9PXGOmuEU0sqr

	f, _ := os.Open(*csvdata)

	r := csv.NewReader(bufio.NewReader(f))
	result, _ := r.ReadAll()

	fout, _ := os.Create("/Users/adityabansal/kineticdevs/go/src/github.com/wearkinetic/myplayground/addtokeendate/testout.csv")
	writer := csv.NewWriter(fout)
	defer writer.Flush()

	var td int64

	for _, row := range result {
		tkeentimestamp, _ := time.Parse(TIME_FORMAT, row[0])
		tend, _ := time.Parse(TIME_FORMAT, row[7])
		tstart, _ := time.Parse(TIME_FORMAT, row[10])

		if row[5] == "34d3n426LdxamFFN" {
			td = tdiff1
		} else {
			td = tdiff2
		}
		row[0] = time.Unix(tkeentimestamp.Unix()+td, 0).UTC().Format(TIME_FORMAT)
		row[7] = time.Unix(tend.Unix()+td, 0).UTC().Format(TIME_FORMAT)
		row[10] = time.Unix(tstart.Unix()+td, 0).UTC().Format(TIME_FORMAT)

		// newepochtime := t.Unix() + tdiff
		// fmt.Printf("t: %s, oldtime: %d, newtime: %d (tdiff: %d), now new: %s\n", row[10], t.Unix(), newepochtime, tdiff, time.Unix(newepochtime, 0).UTC().Format(TIME_FORMAT))

		_ = writer.Write(row)
	}
}

func tFromF(tInt int64) time.Time {
	millisec := tInt % 1000
	sec := (tInt - millisec) / 1000
	nsec := millisec * 1000000
	return time.Unix(sec, nsec).UTC()
}
