package main

import (
	"context"
	"log"

	"github.com/whosonfirst/go-whosonfirst-iterate-bucket/v3"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3/app/emit"
	_ "gocloud.dev/blob/fileblob"
)

func main() {

	ctx := context.Background()

	err := bucket.RegisterBucketIterators(ctx)

	if err != nil {
		log.Fatalf("Failed to register bucket schemes, %v", err)
	}

	err = emit.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to emit record, %v", err)
	}
}
