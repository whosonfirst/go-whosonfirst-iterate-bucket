package bucket

import (
	"context"
	"net/url"
	"strings"
	"log/slog"
	
	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
	"gocloud.dev/blob"
)

const PREFIX string = "bucket-"

func init() {
	ctx := context.Background()

	for _, scheme := range blob.DefaultURLMux().BucketSchemes() {
		scheme = PREFIX + scheme
		slog.Info("Register iterator scheme", "scheme", scheme)
		iterate.RegisterIterator(ctx, scheme, NewBucketIterator)
	}
}

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

	return iterate.NewFSIterator(ctx, bucket_uri, bucket)
}
