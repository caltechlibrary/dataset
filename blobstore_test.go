package dataset

import (
	"flag"
	"fmt"
	"os"
	"testing"

	// Caltech Library Packages
	"github.com/caltechlibrary/storage"
)

var (
	S3Bucket string
	GSBucket string
)

func TestS3(t *testing.T) {
	verbose := false
	if S3Bucket == "" {
		fmt.Fprintf(os.Stderr, "Skipping S3 tests, no bucket\n")
		return
	}
	// Cleanup stale test collections collections
	bucketPath := fmt.Sprintf("%s/testdata/", S3Bucket)
	store, err := storage.GetStore(bucketPath)
	if err == nil {
		err = store.RemoveAll("testdata")
		if err != nil {
			t.Errorf("Unable to remove %s before running tests, %s", bucketPath, err)
			t.FailNow()
		}

	}

	collectionURI1 := fmt.Sprintf("%s/testdata/blob_b.ds", S3Bucket)

	c1, err := InitCollection(collectionURI1)
	if err != nil {
		t.Errorf("expected to create %q, got %s", collectionURI1, err)
		t.FailNow()
	}
	defer func() {
		err = c1.Close()
		if err != nil {
			t.Errorf("expected to close collection, %s", err)
			t.FailNow()
		}
	}()
	if c1.Store.Type != storage.S3 {
		t.Errorf("expected storaged type S3 (%d), got %d", storage.S3, c1.Store.Type)
		t.FailNow()
	}

	key := "one"
	objSrc := []byte(`{"one": 1}`)

	if c1.KeyExists(key) {
		err = c1.UpdateJSON(key, objSrc)
	} else {
		err = c1.CreateJSON(key, objSrc)
	}
	if err != nil {
		t.Errorf("expect err == nil for create/update (%q, %q), %s", key, objSrc, err)
		t.FailNow()
	}

	err = analyzer(collectionURI1, verbose)
	if err != nil {
		t.Errorf("shouldn't have an error for Analyser on %s, %s", collectionURI1, err)
		t.FailNow()
	}

	collectionURI2 := fmt.Sprintf("%s/testdata/blob_p.ds", S3Bucket)
	keys1 := c1.Keys()
	err = c1.Clone(collectionURI2, keys1, verbose)
	if err != nil {
		t.Errorf("should be able to clone %q to %q, %s", collectionURI1, collectionURI2, err)
		t.FailNow()
	}

	c2, err := openCollection(collectionURI2)
	if err != nil {
		t.Errorf("expected err == nil, got %s for %s", err, collectionURI2)
	}
	keys2 := c2.Keys()
	if len(keys1) != len(keys2) {
		t.Errorf("expected %d keys1, got %d keys2", len(keys1), len(keys2))
	}

	err = analyzer(collectionURI2, verbose)
	if err != nil {
		t.Errorf("shouldn't have an error for Analyser on %s, %s", collectionURI2, err)
		t.FailNow()
	}
}

func TestGS(t *testing.T) {
	verbose := false
	if GSBucket == "" {
		fmt.Fprintf(os.Stderr, "Skipping GS tests, no bucket\n")
		return
	}
	// Cleanup stale test collections collections
	bucketPath := fmt.Sprintf("%s/testdata/", GSBucket)
	store, err := storage.GetStore(bucketPath)
	if err == nil {
		err = store.RemoveAll("testdata")
		if err != nil {
			t.Errorf("Unable to remove %s before running tests, %s", bucketPath, err)
			t.FailNow()
		}

	}

	collectionURI1 := fmt.Sprintf("%s/testdata/blob_b.ds", GSBucket)

	c1, err := InitCollection(collectionURI1)
	if err != nil {
		t.Errorf("expected to create %q, got %s", collectionURI1, err)
		t.FailNow()
	}
	defer func() {
		err = c1.Close()
		if err != nil {
			t.Errorf("expected to close collection, %s", err)
			t.FailNow()
		}
	}()
	if c1.Store.Type != storage.GS {
		t.Errorf("expected storaged type GS (%d), got %d", storage.GS, c1.Store.Type)
		t.FailNow()
	}

	key := "one"
	objSrc := []byte(`{"one": 1}`)

	if c1.KeyExists(key) {
		err = c1.UpdateJSON(key, objSrc)
	} else {
		err = c1.CreateJSON(key, objSrc)
	}
	if err != nil {
		t.Errorf("expect err == nil for create/update (%q, %q), %s", key, objSrc, err)
		t.FailNow()
	}

	err = analyzer(collectionURI1, verbose)
	if err != nil {
		t.Errorf("shouldn't have an error for Analyser on %s, %s", collectionURI1, err)
		t.FailNow()
	}

	collectionURI2 := fmt.Sprintf("%s/testdata/blob_p.ds", GSBucket)
	keys1 := c1.Keys()
	err = c1.Clone(collectionURI2, keys1, verbose)
	if err != nil {
		t.Errorf("should be able to clone %q to %q, %s", collectionURI1, collectionURI2, err)
		t.FailNow()
	}

	c2, err := openCollection(collectionURI2)
	if err != nil {
		t.Errorf("expected err == nil, got %s for %s", err, collectionURI2)
	}
	keys2 := c2.Keys()
	if len(keys1) != len(keys2) {
		t.Errorf("expected %d keys1, got %d keys2", len(keys1), len(keys2))
	}

	err = analyzer(collectionURI2, verbose)
	if err != nil {
		t.Errorf("shouldn't have an error for Analyser on %s, %s", collectionURI2, err)
		t.FailNow()
	}
}

func setupBlobStoreTests(m *testing.M) {
	flag.StringVar(&S3Bucket, "s3", "", "Run S3 tests with bucketname")
	flag.StringVar(&GSBucket, "gs", "", "Run GS tests with bucketname")
	flag.Parse()
}
