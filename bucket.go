package bucket

import (
	"context"
	"github.com/whosonfirst/go-ioutil"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/emitter"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/filters"
	"gocloud.dev/blob"
	"io"
	"net/url"
	"strings"
)

const PREFIX string = "bucket-"

type BucketEmitter struct {
	emitter.Emitter
	bucket  *blob.Bucket
	filters filters.Filters
}

func init() {
	ctx := context.Background()

	for _, scheme := range blob.DefaultURLMux().BucketSchemes() {
		scheme = PREFIX + scheme
		emitter.RegisterEmitter(ctx, scheme, NewBucketEmitter)
	}
}

func NewBucketEmitter(ctx context.Context, uri string) (emitter.Emitter, error) {

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

	em := &BucketEmitter{
		bucket:  bucket,
		filters: f,
	}

	return em, nil
}

func (em *BucketEmitter) WalkURI(ctx context.Context, emitter_cb emitter.EmitterCallbackFunc, uri string) error {

	// add go routines
	// add throttles

	var list func(context.Context, *blob.Bucket, string) error

	list = func(ctx context.Context, b *blob.Bucket, prefix string) error {

		iter := b.List(&blob.ListOptions{
			Delimiter: "/",
			Prefix:    prefix,
		})

		for {
			obj, err := iter.Next(ctx)

			if err == io.EOF {
				break
			}

			if err != nil {
				return err
			}

			if obj.IsDir {

				err := list(ctx, b, obj.Key)

				if err != nil {
					return err
				}

				continue
			}

			r, err := em.bucket.NewReader(ctx, obj.Key, nil)

			if err != nil {
				return err
			}

			defer r.Close()

			fh, err := ioutil.NewReadSeekCloser(r)

			if err != nil {
				return err
			}

			if em.filters != nil {

				ok, err := em.filters.Apply(ctx, fh)

				if err != nil {
					return err
				}

				if !ok {
					return nil
				}

				_, err = fh.Seek(0, 0)

				if err != nil {
					return err
				}
			}

			err = emitter_cb(ctx, obj.Key, fh)

			if err != nil {
				return err
			}
		}

		return nil
	}

	return list(ctx, em.bucket, uri)
}
