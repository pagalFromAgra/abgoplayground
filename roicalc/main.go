package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// StartTime = "2017-05-15T06:00:00.000Z"
	// EndTime   = "2017-05-23T00:00:00.000Z"

	StartTime = "2017-05-24T06:00:00.000Z"
	EndTime   = "2017-07-23T00:00:00.000Z"

	TIME_LAYOUT = "2006-01-02T15:04:05.000Z"

	METRIC_MAX_SAG_VEL   = "maximum_sagittal_velocity"
	METRIC_MAX_SAG_ANGLE = "max_sagittal_angle"

	MIN_SAGITTAL_ANGLE = 40.0
	MAX_SAGITTAL_ANGLE = 120.0

	HRL_SAGITTAL_ANGLE = 72.0
)

type Metric struct {
	avg float64
	n   float64
}

type Metrics struct {
	Name        string `json:"name"`
	NHRLWindows int    `json:"no_hrl_windows"`
	MaxSagAngle Metric `json:"max_sagittal_angle"`
	MaxSagVel   Metric `json:"max_sagittal_velocity"`
}

type SkipWindow struct {
	skip bool
	wuid string
}

func main() {
	m := make(map[string]*Metrics)

	_ = LoadDevices(os.Args[1], m)

	_ = LoadMetrics(os.Args[2], m)

	// Sort the keys for printing
	mk := sortKeys(m)

	groupedMaxSagAngle := Metric{
		avg: 0.0,
		n:   0.0,
	}

	groupedMaxSagVel := Metric{
		avg: 0.0,
		n:   0.0,
	}

	for _, k := range mk {
		// fmt.Println(k, m[k].Name, m[k].MaxSagAngle.avg, m[k].MaxSagVel.avg, m[k].MaxSagAngle.n, m[k].MaxSagVel.n)
		fmt.Printf("%s\t%0.1f\t%0.1f\t%d\n", k, m[k].MaxSagAngle.avg, m[k].MaxSagVel.avg, m[k].NHRLWindows)

		if m[k].MaxSagAngle.n > 0 {
			groupedMaxSagAngle.avg += m[k].MaxSagAngle.avg * m[k].MaxSagAngle.n
			groupedMaxSagAngle.n += m[k].MaxSagAngle.n

			groupedMaxSagVel.avg += m[k].MaxSagVel.avg * m[k].MaxSagVel.n
			groupedMaxSagVel.n += m[k].MaxSagVel.n
		}
	}

	groupedMaxSagAngle.avg /= groupedMaxSagAngle.n
	groupedMaxSagVel.avg /= groupedMaxSagVel.n

	fmt.Printf("grouped Max Sag Angle: %0.2f\n", groupedMaxSagAngle.avg)
	fmt.Printf("grouped Max Sag Velocity: %0.2f\n", groupedMaxSagVel.avg)

	// for k, v := range m {
	// 	fmt.Println(k, v.Name, v.MaxSagAngle.avg, v.MaxSagVel.avg)
	// }

}

// LoadData creates all the payloads
func LoadDevices(file string, m map[string]*Metrics) error {

	f, _ := os.Open(file)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		line := scanner.Text()
		arr := strings.Split(line, ",")
		if len(arr) < 4 {
			log.Fatal("Err line size, arr: ", arr)
		}

		m[arr[3]] = &Metrics{
			Name:        arr[0],
			MaxSagAngle: Metric{0.0, 0},
			MaxSagVel:   Metric{0.0, 0},
		}

	}

	return nil
}

func LoadMetrics(file string, data map[string]*Metrics) error {

	s, _ := time.Parse(TIME_LAYOUT, StartTime)
	e, _ := time.Parse(TIME_LAYOUT, EndTime)

	f, _ := os.Open(file)
	defer f.Close()

	skipThisWindow := false

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		line := scanner.Text() // Format: 0VeZZD26dOsSbQVG__2017-07-06T10:41:27.327Z__2017-07-06T10:41:28.286Z__average_sagittal_acceleration	-63.60103408208558
		if strings.Contains(line, "discontinuity") {
			continue
		}

		kvarr := strings.Split(line, "\t")

		karr := strings.Split(kvarr[0], "__")
		wuid := karr[0]
		stime, _ := time.Parse(TIME_LAYOUT, karr[1])
		etime, _ := time.Parse(TIME_LAYOUT, karr[2])
		metric := karr[3]

		if stime.Unix() < s.Unix() || etime.Unix() > e.Unix() {
			// log.Println("time beyond range")
			continue
		}

		if _, ok := data[wuid]; !ok {
			// log.Println("map for this key doesn't exist")
			continue
		}

		if metric != METRIC_MAX_SAG_ANGLE && metric != METRIC_MAX_SAG_VEL {
			// log.Println("not the metric, ", metric)
			continue
		}

		v, _ := strconv.ParseFloat(kvarr[1], 64)

		if metric == METRIC_MAX_SAG_ANGLE {
			if v < MIN_SAGITTAL_ANGLE || v > MAX_SAGITTAL_ANGLE {
				skipThisWindow = true
				continue
			}
			// Now recalculate the averages
			data[wuid].MaxSagAngle.avg = (data[wuid].MaxSagAngle.avg*data[wuid].MaxSagAngle.n + v) / (data[wuid].MaxSagAngle.n + 1)
			data[wuid].MaxSagAngle.n++

			if v > HRL_SAGITTAL_ANGLE {
				data[wuid].NHRLWindows++
			}

		}

		if metric == METRIC_MAX_SAG_VEL {
			if skipThisWindow { // ASSUMPTION: Velocity is listed after the angle
				skipThisWindow = false
				continue
			}
			// Now recalculate the averages
			data[wuid].MaxSagVel.avg = (data[wuid].MaxSagVel.avg*data[wuid].MaxSagVel.n + v) / (data[wuid].MaxSagVel.n + 1)
			data[wuid].MaxSagVel.n++
		}

	}
	return nil
}

func sortKeys(m map[string]*Metrics) []string {
	mk := make([]string, len(m))
	i := 0
	for k, _ := range m {
		mk[i] = k
		i++
	}
	sort.Strings(mk)

	return mk
}
