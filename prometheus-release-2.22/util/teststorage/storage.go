package teststorage

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/prometheus/prometheus/tsdb"
	"github.com/prometheus/prometheus/util/testutil"
)

// New returns a new TestStorage for testing purposes
// that removes all associated files on closing.
func New(t testutil.T) *TestStorage {
	dir, err := ioutil.TempDir("", "test_storage")
	if err != nil {
		t.Fatalf("Opening test dir failed: %s", err)
	}

	// Tests just load data for a series sequentially. Thus we
	// need a long appendable window.
	opts := tsdb.DefaultOptions()
	opts.MinBlockDuration = int64(24 * time.Hour / time.Millisecond)
	opts.MaxBlockDuration = int64(24 * time.Hour / time.Millisecond)
	db, err := tsdb.Open(dir, nil, nil, opts)
	if err != nil {
		t.Fatalf("Opening test storage failed: %s", err)
	}
	return &TestStorage{DB: db, dir: dir}
}

type TestStorage struct {
	*tsdb.DB
	dir string
}

func (s TestStorage) Close() error {
	if err := s.DB.Close(); err != nil {
		return err
	}
	return os.RemoveAll(s.dir)
}
