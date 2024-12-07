package bucket

import (
	"context"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"net/url"
	"strings"
	"sync"

	"github.com/whosonfirst/go-ioutil"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3/filters"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3/iterator"
	"gocloud.dev/blob"
)

const PREFIX string = "bucket-"

var mu_iterators = new(sync.Map)

type BucketIterator struct {
	iterator.Iterator
	bucket  *blob.Bucket
	filters filters.Filters
}

func init() {

	ctx := context.Background()
	err := RegisterBucketIterators(ctx)

	if err != nil {
		panic(err)
	}
}

// RegisterBucketIterators will explicitly register all the schemes...
func RegisterBucketIterators(ctx context.Context) error {

	to_register := make([]string, 0)

	for _, scheme := range blob.DefaultURLMux().BucketSchemes() {
		to_register = append(to_register, scheme)
	}

	for _, scheme := range to_register {

		scheme = PREFIX + scheme

		_, registered := mu_iterators.Load(scheme)

		if registered {
			continue
		}

		err := iterator.RegisterIterator(ctx, scheme, NewBucketIterator)

		if err != nil {
			return fmt.Errorf("Failed to register iterator for '%s', %w", scheme, err)
		}

		mu_iterators.Store(scheme, true)
	}

	return nil
}

func NewBucketIterator(ctx context.Context, uri string) (iterator.Iterator, error) {

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

	q := u.Query()

	f, err := filters.NewQueryFiltersFromQuery(ctx, q)

	if err != nil {
		return nil, err
	}

	it := &BucketIterator{
		bucket:  bucket,
		filters: f,
	}

	return it, nil
}

func (it *BucketIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[iterator.Record, error] {

	return func(yield func(iterator.Record, error) bool) {

		for _, uri := range uris {
			for r, err := range it.iterate(ctx, uri) {
				yield(r, err)
			}
		}
	}

}

func (it *BucketIterator) iterate(ctx context.Context, uri string) iter.Seq2[iterator.Record, error] {

	logger := slog.Default()
	logger = logger.With("uri", uri)

	return func(yield func(iterator.Record, error) bool) {

		// Convenience to account for gocloud.dev/blob -isms
		uri = strings.TrimLeft(uri, "/")

		var list func(context.Context, *blob.Bucket, string) error

		list = func(ctx context.Context, b *blob.Bucket, prefix string) error {

			iter := b.List(&blob.ListOptions{
				Delimiter: "/",
				Prefix:    prefix,
			})

			for {

				logger := slog.Default()
				logger = logger.With("uri", uri)

				obj, err := iter.Next(ctx)

				if err == io.EOF {
					break
				}

				if err != nil {
					logger.Debug("Iterator reported an error", "error", err)
					return err
				}

				logger = logger.With("key", obj.Key)

				if obj.IsDir {

					err := list(ctx, b, obj.Key)

					if err != nil {
						logger.Debug("Listing directory failed", "error", err)
						return err
					}

					continue
				}

				bucket_r, err := it.bucket.NewReader(ctx, obj.Key, nil)

				if err != nil {
					logger.Debug("Failed to create new reader", "error", err)
					return err
				}

				defer bucket_r.Close()

				r, err := ioutil.NewReadSeekCloser(bucket_r)

				if err != nil {
					logger.Debug("Failed to create ReadSeekCloser from bucket reader", "error", err)
					return err
				}

				if it.filters != nil {

					ok, err := it.filters.Apply(ctx, r)

					if err != nil {
						logger.Debug("Failed to apply filters", "error", err)
						return err
					}

					if !ok {
						logger.Debug("No matches after applying filters, skipping")
						return nil
					}

					_, err = r.Seek(0, 0)

					if err != nil {
						logger.Debug("Failed to rewind reader", "error", err)
						return err
					}
				}

				logger.Debug("Yield new record")

				iter_r := iterator.NewRecord(obj.Key, r)
				yield(iter_r, nil)
			}

			return nil
		}

		err := list(ctx, it.bucket, uri)

		if err != nil {
			logger.Debug("Failed to list bucket", "error", err)
			yield(nil, err)
		}
	}
}
