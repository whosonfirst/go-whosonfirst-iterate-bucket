# go-whosonfirst-index-bucket

Go package implementing `whosonfirst/go-whosonfirst-iterate/v3` functionality for GoCloud blob resources.

## Important

Documentation for this package is incomplete and will be updated shortly.

## Example

```
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync/atomic"

	_ "github.com/whosonfirst/go-whosonfirst-iterate-bucket/v3"
	_ "gocloud.dev/blob/fileblob"
	
	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

func main() {

	var count int64
	var iterator_uri string

	flag.Parse()
	iterator_sources := flag.Args()
	
	iterator_uri = "bucket-file:///"
	count = 0

	ctx := context.Background()

	it, _ := iterator.NewIterator(ctx)

	for _, err := range it.Iterate(ctx, iterator_uri, iterator_sources...){

		if err != nil {
			log.Fatal(err)
		}
		
		atomic.AddInt64(&count, 1)
	}

	log.Printf("Counted %d records (saw %d records)\n", count, iter.Seen)
}
```

_Error handling omitted for the sake of brevity._

## Schemes (and `blob.Bucket` providers)

This package does not load _any_ [blob.Bucket](https://gocloud.dev/howto/blob/) providers by default. You will need to do that explicitly in your code.

Importantly, `gocloud.dev/blob` import statements need to be declared _before_ you import `whosonfirst/go-whosonfirst-iterate-bucket`. The easiest way to do this is with two separate `import()` statements. It's not elegant, but it works.

The reason this is necessary is because the initializer function in the `bucket.go` package registers available emitter schemes by iterating over the registered GoCloud `blob.Bucket` schemes and and prefixing them with "bucket-".

For example the GoCloud "file://" scheme becomes "bucket-file://" in order to ensure that there aren't namespace collisions with schemes registered by other packages (for example "file://" is already registered by the `go-whosonfirst-iterate` package).

## See also

* https://github.com/whosonfirst/go-whosonfirst-iterate
* https://gocloud.dev/howto/blob/