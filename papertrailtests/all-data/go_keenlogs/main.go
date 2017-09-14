package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wearkinetic/awss3"
	"github.com/wearkinetic/beutils"
	"github.com/wearkinetic/keenlogs/uploadlogs"
)

const (
	ENV_AWS_KEY              = "AWS_KEY"
	ENV_AWS_SECRET           = "AWS_SECRET"
	ENV_AWS_S3_BUCKET        = "AWS_S3_BUCKET"
	ENV_AWS_S3_BUCKET_REGION = "AWS_S3_BUCKET_REGION"

	REGION_US_EAST_1 = "us-east-1"
	REGION_US_WEST_2 = "us-west-2"

	ACCESS_CONTROL_PUBLIC_READ_WRITE = "public-read-write"
	ACCESS_CONTROL_PUBLIC_READ       = "public-read"
	ACCESS_CONTROL_PRIVATE           = "private"

	S3_PUBLIC_URL = "https://s3.amazonaws.com"

	DEFAULT_S3_BUCKET = "kinetic-device-logs"

	S3_FILE_KEY_PREFIX = "analytics" // S3 path: kinetic-device-logs/analytics/localFile
)

type inlogdata struct {
	Time string `jston:"time"`
	Type string `json:"type"`
	Log  string `jston:"log"`
}

var listButtonEvents = []string{"BUTTON_PRESS_RISKY_LIFTS_GOAL", "BUTTON_PRESS_TIME", "BUTTON_PRESS_STAY_SAFE"}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Need the source CSV file")
		os.Exit(3)
	}

	sourcefile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln("Could not open sourcefile, err: ", err)
	}
	defer sourcefile.Close()

	// To get the location for a device
	client, errc := beutils.NewHTTPClient()
	if errc != nil {
		log.Fatal(errc)
	}

	// err := copyToS3(os.Args[1])
	// if err != nil {
	// 	fmt.Printf("copying to S3 failed, err: %v", err)
	// }

	scanner := bufio.NewScanner(sourcefile)
	for scanner.Scan() {
		// fmt.Println(scanner.Text())

		if scanner.Text() == "" { // Keep going through empty line
			fmt.Println(err)
			continue
		}

		splitted := strings.Split(scanner.Text(), "|@|")

		if len(splitted) < 2 { // Keep going if can't parse the line
			fmt.Println(err)
			continue
		}

		tempkey := strings.Split(splitted[0], "-") // Some keys have -1 or -2 by papertrail
		key := tempkey[0]

		indata := &inlogdata{}

		err := json.Unmarshal([]byte(splitted[1]), indata)
		if err != nil {
			fmt.Println(err)
			continue
		}

		record := &uploadlogs.Payload{
			DateTime: indata.Time,
			Device:   key,
			Event:    "display.switch",
		}

		record.Company, _, err = beutils.GetLocationSKU(client, key)
		if err != nil {
			fmt.Println(err)
			continue
		}

		record.Value = 1

		for _, button := range listButtonEvents {
			if strings.Contains(indata.Log, button) {
				record.Info = button
				recordJson, err := json.Marshal(record)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("%s\n", string(recordJson))

				erru := uploadlogs.Handler(record)
				if erru != nil {
					fmt.Printf("Error uploading to keen: %v", erru)
				}
			}
		}
	}

	// for {
	// 	line, err := reader.Read()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatalln("Error reading sourcefile, err: ", err)
	// 	}
	//
	// 	log.Println(line)
	//
	// 	record := &uploadlogs.Payload{
	// 		DateTime: line[0],
	// 		Company:  line[1],
	// 		Device:   line[2],
	// 		Event:    line[3],
	// 		Info:     line[7],
	// 	}

	// record.Vbat, _ = strconv.ParseInt(line[4], 10, 64)
	// record.Temprature, _ = strconv.ParseInt(line[5], 10, 64)
	// record.Value, _ = strconv.ParseInt(line[6], 10, 64)

	// recordJson, _ := json.Marshal(record)
	// fmt.Printf("%s\n", string(recordJson))
}

func copyToS3(localFile string) error {

	s3 := awss3.NewSession(os.Getenv(ENV_AWS_S3_BUCKET_REGION))

	s3bucket := os.Getenv(ENV_AWS_S3_BUCKET)

	f, err := os.Open(localFile)
	if err != nil {
		return err
	}

	s3lockey := fmt.Sprintf("%s/%s", S3_FILE_KEY_PREFIX, localFile)

	putError := s3.Put(s3bucket, f, s3lockey, ACCESS_CONTROL_PUBLIC_READ, map[string]*string{})
	if putError != nil {
		return putError
	}

	return nil

}
