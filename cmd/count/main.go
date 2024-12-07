package main

import (
	"context"
	"log"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3/app/count"
	_ "gocloud.dev/blob/fileblob"
	_ "github.com/whosonfirst/go-whosonfirst-iterate-bucket/v3"		
)

func main() {

	ctx := context.Background()
	err := count.Run(ctx)

	if err != nil {
		log.Fatalf("Failed to count record, %v", err)
	}
}
