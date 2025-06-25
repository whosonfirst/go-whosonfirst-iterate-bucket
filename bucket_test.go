package bucket

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"testing"

	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

func TestBucketIterator(t *testing.T) {

	ctx := context.Background()

	abs_path, err := filepath.Abs("fixtures/data")

	if err != nil {
		t.Fatalf("Failed to derive absolute path for fixtures, %v", err)
	}

	iter_uri := fmt.Sprintf("bucket-file://%s", abs_path)

	it, err := iterate.NewIterator(ctx, iter_uri)

	if err != nil {
		t.Fatalf("Failed to create bucket iterator, %v", err)
	}

	count := 0

	for rec, err := range it.Iterate(ctx, ".") {

		if err != nil {
			t.Fatalf("Failed to walk filesystem, %v", err)
			break
		}

		_, err = io.ReadAll(rec.Body)

		if err != nil {
			t.Fatalf("Failed to read body for %s, %v", rec.Path, err)
		}

		_, err = rec.Body.Seek(0, 0)

		if err != nil {
			t.Fatalf("Failed to rewind body for %s, %v", rec.Path, err)
		}

		_, err = io.ReadAll(rec.Body)

		if err != nil {
			t.Fatalf("Failed second read body for %s, %v", rec.Path, err)
		}

		count += 1
	}

	expected := 37

	if count != expected {
		t.Fatalf("Expected %d records, but counted %d", expected, count)
	}

}
