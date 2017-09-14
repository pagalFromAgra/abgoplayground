package main

import (
	"log"
	"os"
	"time"

	"github.com/wearkinetic/awss3"
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

	DEFAULT_S3_BUCKET = "kinetic-device-logs/analytics"
)

func main() {

	s3 := awss3.NewSession(os.Getenv(ENV_AWS_S3_BUCKET_REGION))

	s3bucket := os.Getenv(ENV_AWS_S3_BUCKET)

	localFile := os.Args[1]
	f, err := os.Open(localFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	putError := s3.Put(s3bucket, f, localFile, ACCESS_CONTROL_PUBLIC_READ, map[string]*string{}, 1*time.Minute)
	if putError != nil {
		log.Fatalln(putError)
	}

}
