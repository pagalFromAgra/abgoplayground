package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/wearkinetic/awss3"
	"github.com/wearkinetic/keen"

	"github.com/wearkinetic/beutils"
)

type Result struct {
	ID        string `json:"id"`
	Timeframe struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}
	Lifts       int `json:"lifts"`
	Time_Active int `json:"time_active"`
	Lift_Rate   int `json:"lift_rate"`
}

func main() {

	// --------
	// STEP 1. Get the data from Keen
	// --------
	k, err := keen.NewFromEnv()
	company := keen.Company{
		Keen:       k,          // Keen instance
		Name:       os.Args[1], // company name as stored in keen
		ShiftHours: 8}          // shift length in hours
	if err != nil {
		log.Fatal(err)
	}

	client, errc := beutils.NewHTTPClient()
	if errc != nil {
		log.Fatal(errc)
	}

	devicesAtlocation, err := beutils.GetAllDevicesAtLocation(client, company.Name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("devicesAtlocation = ", devicesAtlocation)

	var startDate string
	var endDate string

	if len(os.Args) < 3 {
		ty := time.Now().AddDate(0, 0, -1) // Yesterday's date
		startDate = ty.AddDate(0, 0, -1).Format("2006-01-02")
		endDate = ty.Format("2006-01-02")
	} else {
		startDate = os.Args[2]
		endDate = os.Args[3]
	}

	dates, err := beutils.GetDateRange(startDate, endDate)
	if err != nil {
		log.Fatal(err)
	}

	for _, checkdate := range dates {

		fmt.Printf("Checking for the date: %s\n", checkdate)

		response, err := company.GetData(
			checkdate+"T00:00:00-00:00", // start of timeframe to get
			checkdate+"T23:59:59-00:00", // end of timeframe to get
			"daily") // interval to group into
		if err != nil {
			log.Fatal(err)
		}

		// responseByTimeframe := response.ByTimeframeByEmployee() // group first by timeframe then by employee for easy marshalling to JSON
		// responseByEmployee := response.ByTimeframeByEmployee() // the reverse

		// fmt.Println(goutil.Pretty(*responseByEmployee))

		// --------
		// STEP 2. Setup S3 session
		// --------
		session := awss3.NewSession(awss3.REGION_US_EAST_1)

		// --------
		// STEP 4. For each device key assigned to the company, compare the data between Keen and S3
		// --------
		countNoData := 0
		serialCounter := 0

		for _, device := range devicesAtlocation {

			// First check if this device is assigned to an employee
			employeeName, employeeID := beutils.GetEmployeeInfo(client, device)
			if employeeName == "" {
				continue
			}

			lifts := 0
			activetime := 0

			for _, dt := range *response.Employees {
				if dt.ID == employeeID {
					lifts = dt.Lifts
					activetime = dt.ActiveSeconds
				}
			}

			list, err := session.List("kinetic-device-data", "raw/"+device+"/"+checkdate)
			if err != nil {
				log.Println("Couldn't read file list")
			}

			serialCounter++

			// Each file has either 1 data point (40ms) or 5mins (300s) of data
			// len(list)*0.04 <= activetime <= len(list)*5*60

			missingS3data := ""
			missingKeendata := ""
			over15hrs := ""
			lessThan_1_HRPperHr := ""

			marker := false

			if activetime > (len(list)*5*60 + 3600) { // Because there can be 1 hr of overlap from the other day in Keen data
				if len(list) == 0 {
					missingS3data = "NO-S3"
				} else {
					missingS3data = "PART-S3"
				}
				marker = true
			}

			if activetime < len(list)*1*6 && len(list) > 1 { // At least 1 sec of data and > 1 file in S3
				missingKeendata = "NO-KEEN"
				marker = true
			}

			if activetime > 3600*15 {
				over15hrs = "GT-15HRS"
				marker = true
			}

			if activetime > 0 && (lifts*3600/activetime) < 2 {
				lessThan_1_HRPperHr = "LT-1-HRP"
				marker = true
			}

			if marker || !marker {
				employeeName, _ = beutils.GetEmployeeInfo(client, device)
				fmt.Printf("%d.\t%.15s\t%s\t%d\t%d\t%d\t\t\t%s\t%s\t%s\t%s\n", serialCounter, employeeName, device, lifts, activetime, len(list), missingS3data, missingKeendata, over15hrs, lessThan_1_HRPperHr)
			}

			if activetime == 0 && lifts == 0 && len(list) == 0 {
				countNoData++
			}
		}
		fmt.Printf("%d/%d with no data\n\n", countNoData, len(devicesAtlocation))
	}

	//
	// for _, row := range allinfo {
	//
	// 	if strings.Contains(row.side, "left") {
	// 		fmt.Printf("%s;", row.device)
	// 	}
	// }

	// i := 0
	// for _, dt := range *response.Employees {
	// 	// fmt.Printf("%d %s %d %d\n", i, dt.ID, dt.Lifts, dt.ActiveSeconds)
	// 	fmt.Println(dt)
	// 	i++
	// }

}
