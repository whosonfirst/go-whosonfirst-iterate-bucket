# go-whosonfirst-index-bucket

Go package implementing go-whosonfirst-iterate/emitter functionality for GoCloud blob resources.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/whosonfirst/go-whosonfirst-iterate.svg)](https://pkg.go.dev/github.com/whosonfirst/go-whosonfirst-iterate/v3)

## Example

Version 3.x of this package introduce major, backward-incompatible changes from earlier releases. That said, migragting from version 2.x to 3.x should be relatively straightforward as a the _basic_ concepts are still the same but (hopefully) simplified. Where version 2.x relied on defining a custom callback for looping over records version 3.x use Go's [iter.Seq2](https://pkg.go.dev/iter) iterator construct to yield records as they are encountered.

For example:

```
import (
	"context"
	"flag"
	"log"

	_ "github.com/whosonfirst/go-whosonfirst-iterate-bucket/v3"
	_ "gocloud.dev/blob/fileblob"
	
	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

func main() {

     	var iterator_uri string

	flag.StringVar(&iterator_uri, "iterator-uri", "bucket-file:///". "A registered whosonfirst/go-whosonfirst-iterate/v3.Iterator URI.")
	ctx := context.Background()
	
	iter, _:= iterate.NewIterator(ctx, iterator_uri)

	paths := flag.Args()
	
	for rec, _ := range iter.Iterate(ctx, paths...) {
	    	defer rec.Body.Close()
		log.Printf("Indexing %s\n", rec.Path)
	}
}
```

_Error handling removed for the sake of brevity._

### Version 2.x (the old way)

This is how you would do the same thing using the older version 2.x code:

```
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

	_ "github.com/whosonfirst/go-whosonfirst-iterate-bucket/v2"
	_ "gocloud.dev/blob/fileblob"
	
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
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

This package only loads one [blob.Bucket](https://gocloud.dev/howto/blob/) provider by default: `fileblob://` for reading files from the local disk. All other blob providers will need to be imported manually.

All `gocloud.dev/blob` providers are registered with the "bucket-" prefix. For example the GoCloud "file://" scheme becomes "bucket-file://" in order to ensure that there aren't namespace collisions with schemes registered by other packages (for example "file://" is already registered by the `go-whosonfirst-iterate` package).

Depending on the order of your import statements you may need to explicitly register bucket providers before working with iterators. This is a by-product of the changes (in Go 1.21) to how import statements are handled. It's annoying but the remedy is simply to call the `bucket.RegisterSchemes(ctx)` method.

Under the hood this package is iterating over a `gocloud.dev/blob.Bucket` instance using the `whosonfirst/go-whosonfirst-iterate/v3.FSIterator` instance.

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