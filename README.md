# go-whosonfirst-index-bucket

Go package implementing go-whosonfirst-iterate/emitter functionality for GoCloud blob resources.

## Important

Documentation for this package is incomplete and will be updated shortly.

## Example

```
package main

import (
	_ "gocloud.dev/blob/fileblob"
)

import (
	"context"
	"flag"
	"fmt"
	_ "github.com/whosonfirst/go-whosonfirst-iterate-bucket"
	"github.com/whosonfirst/go-whosonfirst-iterate/iterator"
	"io"
	"log"
	"os"
	"strings"
	"sync/atomic"
)

func main() {

	flag.Parse()
	uris := flag.Args()
	
	emitter_uri = "bucket-file:///"

	var count int64
	count = 0

	emitter_cb := func(ctx context.Context, fh io.ReadSeeker, args ...interface{}) error {
		atomic.AddInt64(&count, 1)
		return nil
	}

	ctx := context.Background()

	iter, _ := iterator.NewIterator(ctx, emitter_uri, emitter_cb)
	iter.IterateURIs(ctx, uris...)

	log.Printf("Counted %d records (saw %d records)\n", count, iter.Seen)
}
```

_Error handling omitted for the sake of brevity._

## Schemes (and `blob.Bucket` providers)

This package does not load _any_ [blob.Bucket](https://gocloud.dev/howto/blob/) providers by default. You will need to do that explicitly in your code.

Importantly, `gocloud.dev/blob` import statements need to be declared _before_ you import `whosonfirst/go-whosonfirst-iterate-bucket`. The easiest way to do this is with two separate `import()` statements. It's not elegant, but it works.

The reason this is necessary is because the initializer function in the `bucket.go` package registers available emitter schemes by iterating over the registered GoCloud `blob.Bucket` schemes and and prefixing them with "bucket-".

For example the GoCloud "file://" scheme becomes "bucket-file://" in order to ensure that there aren't namespace collisions with schemes registered by other packages (for example "file://" is already registered by the `go-whosonfirst-iterate` package).

## Tools

### count

Count files in one or more whosonfirst/go-whosonfirst-iterate/emitter sources.

```
$> ./bin/count -h
Count files in one or more whosonfirst/go-whosonfirst-iterate/emitter sources.
Usage:
	 ./bin/count [options] uri(N) uri(N)
Valid options are:

  -emitter-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/emitter URI. Supported emitter URI schemes are: bucket-file://,directory://,featurecollection://,file://,filelist://,geojsonl://,repo:// (default "bucket-file:///")
```

For example:

```
$> ./bin/count -emitter-uri bucket-file:///usr/local/data/sfomuseum-data-architecture/ data
2021/02/18 13:34:15 time to index paths (1) 810.932513ms
2021/02/18 13:34:15 Counted 1072 records (saw 1072 records)
```

### emit

Publish features from one or more whosonfirst/go-whosonfirst-iterate/emitter sources.

```
> ./bin/emit -h
Publish features from one or more whosonfirst/go-whosonfirst-iterate/emitter sources.
Usage:
	 ./bin/emit [options] uri(N) uri(N)
Valid options are:

  -emitter-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/emitter URI. Supported emitter URI schemes are: bucket-file://,directory://,featurecollection://,file://,filelist://,geojsonl://,repo:// (default "bucket-file:///")
  -geojson
    	Emit features as a well-formed GeoJSON FeatureCollection record.
  -json
    	Emit features as a well-formed JSON array.
  -null
    	Publish features to /dev/null
  -stdout
    	Publish features to STDOUT. (default true)
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-iterate
* https://gocloud.dev/howto/blob/