package bucket

import (
	"context"
	"fmt"
	"iter"
	"log/slog"
	"net/url"
	"strings"
	"sync"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
	"gocloud.dev/blob"
)

const PREFIX string = "bucket-"

// In principle this could also be done with a sync.OnceFunc call but that will
// require that everyone uses Go 1.21 (whose package import changes broke everything)
// which is literally days old as I write this. So maybe a few releases after 1.21.
//
// Also, _not_ using a sync.OnceFunc means we can call RegisterSchemes multiple times
// if and when multiple gomail-sender instances register themselves.

var register_mu = new(sync.RWMutex)
var register_map = map[string]bool{}

// BucketIterator implements the `Iterator` interface for crawling records in a `gocloud.dev/blob.Bucket` bucket.
type BucketIterator struct {
	// bucket is the `gocloud.dev/blob.Bucket` instance where records are stored.
	bucket *blob.Bucket
	// iterator is the underlying `DirectoryIterator` instance for crawling records.
	iterator iterate.Iterator
}

func init() {

	ctx := context.Background()
	err := RegisterSchemes(ctx)

	if err != nil {
		panic(err)
	}
}

// RegisterSchemes will explicitly register all the schemes associated with the `gocloud.dev/blob` interface.
func RegisterSchemes(ctx context.Context) error {

	register_mu.Lock()
	defer register_mu.Unlock()

	for _, scheme := range blob.DefaultURLMux().BucketSchemes() {

		scheme = PREFIX + scheme
		slog.Debug("Register bucket iterator scheme", "scheme", scheme)

		_, exists := register_map[scheme]

		if exists {
			continue
		}

		err := iterate.RegisterIterator(ctx, scheme, NewBucketIterator)

		if err != nil {
			return fmt.Errorf("Failed to register blob writer for '%s', %w", scheme, err)
		}

		register_map[scheme] = true
	}

	return nil
}

// NewBucketIterator() returns a new `BucketIterator` where 'uri' takes the form of:
//
//	bucket-{SCHEME}://?{PARAMETERS}
//
// Where {SCHEME} is a registered `gocloud.dev/blob` driver and {PARAMETERS} may be:
// * `?include=` Zero or more `aaronland/go-json-query` query strings containing rules that must match for a document to be considered for further processing.
// * `?exclude=` Zero or more `aaronland/go-json-query`	query strings containing rules that if matched will prevent a document from being considered for further processing.
// * `?include_mode=` A valid `aaronland/go-json-query` query mode string for testing inclusion rules.
// * `?exclude_mode=` A valid `aaronland/go-json-query` query mode string for testing exclusion rules.
// * `?processes=` An optional number assigning the maximum number of database rows that will be processed simultaneously. (Default is defined by `runtime.NumCPU()`.)
func NewBucketIterator(ctx context.Context, uri string) (iterate.Iterator, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	u.Scheme = strings.Replace(u.Scheme, PREFIX, "", 1)

	bucket_uri := u.String()

	bucket, err := blob.OpenBucket(ctx, bucket_uri)

	if err != nil {
		return nil, err
	}

	fs_iter, err := iterate.NewFSIterator(ctx, bucket_uri, bucket)

	if err != nil {
		return nil, err
	}

	it := &BucketIterator{
		iterator: fs_iter,
		bucket:   bucket,
	}

	return it, nil
}

// Iterate will return an `iter.Seq2[*Record, error]` for each record encountered in 'uris'.
func (it *BucketIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[*iterate.Record, error] {
	return it.iterator.Iterate(ctx, uris...)
}

// Seen() returns the total number of records processed so far.
func (it *BucketIterator) Seen() int64 {
	return it.iterator.Seen()
}

// IsIterating() returns a boolean value indicating whether 'it' is still processing documents.
func (it *BucketIterator) IsIterating() bool {
	return it.iterator.IsIterating()
}

// Close performs any implementation specific tasks before terminating the iterator.
func (it *BucketIterator) Close() error {

	err := it.iterator.Close()

	if err != nil {
		return err
	}

	return it.bucket.Close()
}
