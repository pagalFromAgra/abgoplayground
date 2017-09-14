package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/wearkinetic/awss3"
	"github.com/wearkinetic/beutils"
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
)

func main() {

	s3 := awss3.NewSession(awss3.REGION_US_EAST_1)
	s3bucket := "kinetic-device-logs"

	dates, err := beutils.GetDateRange(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatalln("Failed to get date range, err: ", err)
	}

	// STEP 1: Get all the dates in the given start and end dates

	for _, date := range dates {
		fkey := fmt.Sprintf("papertrail/logs/3931621/dt=%s", date)

		// STEP 2: Get all the .tsv.gz files (created every hour) for a given date

		list, listError := s3.List(s3bucket, fkey)
		if listError != nil {
			log.Fatalln("List error ", listError)
		}
		// fmt.Printf("List of files: %s", list)

		if len(list) > 0 {
			_ = os.MkdirAll(fkey, 0777)
		}

		// STEP 3: Download all the hourly .tsv.gz files and unzip them

		for i, logfile := range list {
			fmt.Printf("%s: downloading %d/%d\n", date, i, len(list))

			got, getError := s3.Get(s3bucket, logfile)
			if getError != nil {
				log.Fatalln("Get error", getError)
			}

			f, err := os.Create(logfile)
			if err != nil {
				log.Fatalln("Create error ", err)
			}

			if _, err := io.Copy(f, got.Body); err != nil {
				log.Fatalln("Failed to copy object to file", err)
			}
			got.Body.Close()
			f.Close()

			unzipCmdStr := fmt.Sprintf("gunzip %s", logfile)
			cmdUnzip := exec.Command("bash", "-c", unzipCmdStr)
			err = cmdUnzip.Run()
			if err != nil {
				log.Fatalln("Error running unzip command, err: ", err)
			}
		}

		// STEP 4: Concatenate all the hourly unzipped .tsv files into a single file of name: <date>.tsv

		catCmdStr := fmt.Sprintf("cat %s/* > %s.tsv", fkey, date)
		cmdCat := exec.Command("bash", "-c", catCmdStr)
		fmt.Println("Running ", catCmdStr)
		err = cmdCat.Run()
		if err != nil {
			log.Fatalln("Error running cat command, err: ", err)
		}

		// file, err := os.Open(fmt.Sprintf("%s.tsv", date))
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer file.Close()
		//
		// client, errc := beutils.NewHTTPClient()
		// if errc != nil {
		// 	log.Fatal(errc)
		// }
		//
		// scanner := bufio.NewScanner(file)
		// for scanner.Scan() {
		// 	valarray := strings.Fields(scanner.Text())
		//
		// 	// Remove "-1" or "-2" etc. from the key added by papertrail
		// 	key := strings.Split(valarray[1], "-")
		//
		// 	location, sku := beutils.GetLocationSKU(client, key[0])
		// 	fmt.Printf("2017-06-%s,%s,%s,%s,%s\n", os.Args[2], key[0], valarray[0], sku, location)
		// }

	}
}
