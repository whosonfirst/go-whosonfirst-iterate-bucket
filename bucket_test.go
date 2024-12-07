package bucket

import (
	"context"
	"log/slog"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3/iterator"
	_ "gocloud.dev/blob/fileblob"
)

func TestBucketIterator(t *testing.T) {

	slog.SetLogLoggerLevel(slog.LevelDebug)

	ctx := context.Background()

	it, err := iterator.NewIterator(ctx, "bucket-file:///")

	if err != nil {
		t.Fatalf("Failed to create directory emitter, %v", err)
	}

	path_data, err := filepath.Abs("fixtures/data")

	if err != nil {
		t.Fatalf("Failed to derive absolute path for fixtures, %v", err)
	}

	path_data = strings.TrimLeft(path_data, "/")

	expected := int32(37)
	count := int32(0)

	for _, err := range it.Iterate(ctx, path_data) {

		if err != nil {
			t.Fatalf("Failed to walk directory, %v", err)
		}

		atomic.AddInt32(&count, 1)
	}

	if count != expected {
		t.Fatalf("Unexpected count: %d", count)
	}

}
