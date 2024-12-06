package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync/atomic"

	_ "github.com/whosonfirst/go-whosonfirst-iterate-bucket/v3"
	_ "gocloud.dev/blob/fileblob"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

func main() {

	var iterator_uri = flag.String("iterator-uri", "bucket-file:///", "")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Count files in one or more whosonfirst/go-whosonfirst-iterate/v3 sources.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options] uri(N) uri(N)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	var count int64
	count = 0

	ctx := context.Background()

	it, err := iterator.NewIterator(ctx)

	if err != nil {
		log.Fatal(err)
	}

	sources := flag.Args()

	for _, err := range iter.Iterate(ctx, iterator_uri, sources...) {
		atomic.AddInt64(&count, 1)
	}

	log.Printf("Counted %d records (saw %d records)\n", count, iter.Seen)
}
