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

	outdata := flag.String("out", "", "a string")

	flag.Parse()

	this, _ := time.Parse(TIME_FORMAT, "2017-10-12T00:00:00.000Z")
	that, _ := time.Parse(TIME_FORMAT, "2017-09-28T00:00:00.000Z")
	tdiff := that.Unix() - this.Unix()

	f, _ := os.Open(*csvdata)

	r := csv.NewReader(bufio.NewReader(f))
	result, _ := r.ReadAll()

	fout, _ := os.Create(*outdata)
	writer := csv.NewWriter(fout)
	defer writer.Flush()

	var td int64

	for _, row := range result {
		tkeentimestamp, _ := time.Parse(TIME_FORMAT, row[0])
		tend, _ := time.Parse(TIME_FORMAT, row[7])
		tstart, _ := time.Parse(TIME_FORMAT, row[10])

		td = tdiff

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
