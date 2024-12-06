package bucket

import (
	"context"
	"io"
	"iter"
	"log/slog"
	"net/url"
	"strings"

	"github.com/whosonfirst/go-ioutil"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3/filters"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3/iterator"
	"gocloud.dev/blob"
)

const PREFIX string = "bucket-"

type BucketIterator struct {
	iterator.Iterator
	bucket  *blob.Bucket
	filters filters.Filters
}

func init() {
	ctx := context.Background()

	for _, scheme := range blob.DefaultURLMux().BucketSchemes() {
		scheme = PREFIX + scheme
		iterator.RegisterIterator(ctx, scheme, NewBucketIterator)
	}
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

		// add go routines
		// add throttles

		// Update to use https://github.com/aaronland/gocloud-blob/tree/main/walk

		var list func(context.Context, *blob.Bucket, string) error

		list = func(ctx context.Context, b *blob.Bucket, prefix string) error {

			iter := b.List(&blob.ListOptions{
				Delimiter: "/",
				Prefix:    prefix,
			})

			for {

				logger := logger.Default()
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

				iter_r := iterator.NewRecord(obj.Key, fh)
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
