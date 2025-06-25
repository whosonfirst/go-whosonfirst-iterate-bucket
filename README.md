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
	_ "github.com/whosonfirst/go-whosonfirst-iterate-bucket/v2"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
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

	emitter_cb := func(ctx context.Context, path string fh io.ReadSeeker, args ...interface{}) error {
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

Count files in one or more whosonfirst/go-whosonfirst-iterate/v3.Iterator sources.

```
$> ./bin/count -h
Count files in one or more whosonfirst/go-whosonfirst-iterate/v3.Iterator sources.
Usage:
	 ./bin/count [options] uri(N) uri(N)
Valid options are:

  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI. Supported iterator URI schemes are: bucket-file://,cwd://,directory://,featurecollection://,file://,filelist://,geojsonl://,null://,repo:// (default "repo://")
  -verbose
    	Enable verbose (debug) logging.
``

### emit

Emit records in one or more whosonfirst/go-whosonfirst-iterate/v3.Iterator sources as structured data.

```
$> ./bin/emit -h
Emit records in one or more whosonfirst/go-whosonfirst-iterate/v3.Iterator sources as structured data.
Usage:
	 ./bin/emit [options] uri(N) uri(N)
Valid options are:

  -geojson
    	Emit features as a well-formed GeoJSON FeatureCollection record.
  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI. Supported iterator URI schemes are: bucket-file://,cwd://,directory://,featurecollection://,file://,filelist://,geojsonl://,null://,repo:// (default "repo://")
  -json
    	Emit features as a well-formed JSON array.
  -null
    	Publish features to /dev/null
  -stdout
    	Publish features to STDOUT. (default true)
  -verbose
    	Enable verbose (debug) logging.
```	

## See also

* https://github.com/whosonfirst/go-whosonfirst-iterate
* https://gocloud.dev/howto/blob/