package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/kurin/blazer/b2"
)

var (
	accountID  = os.Getenv("B2_ACCOUNT_ID")
	accountKey = os.Getenv("B2_ACCOUNT_KEY")

	maxAge     time.Duration
	bucketName string
)

func init() {
	flag.DurationVar(&maxAge, "maxAge", time.Hour*24*31, "max age of files")
	flag.StringVar(&bucketName, "bucket", "", "bucket name")
	flag.Parse()
}

func main() {
	if bucketName == "" {
		log.Fatalf("bucket name must be set")
	}

	ctx := context.Background()

	client, err := b2.NewClient(ctx, accountID, accountKey)
	checkError("init client", err)

	bucket, err := client.Bucket(ctx, bucketName)
	checkError("find bucket: "+bucketName, err)

	iterator := bucket.List(ctx)

	for iterator.Next() {
		obj := iterator.Object()

		attr, err := obj.Attrs(ctx)
		checkError("list attributes: "+obj.Name(), err)

		if attr.UploadTimestamp.Before(time.Now().Add(-maxAge)) {
			log.Printf("deleting file: %s (last modified: %s)", obj.Name(), attr.UploadTimestamp.String())
			err := obj.Delete(ctx)
			checkError("delete file: "+obj.Name(), err)
		}
	}
}

func checkError(what string, err error) {
	if err == nil {
		//log.Printf("completed: %s", what)
		return
	}

	log.Fatalf("could not: %s, err: %s", what, err.Error())
}
