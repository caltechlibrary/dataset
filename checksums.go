package dataset

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func calcChecksum(fName string) (string, error) {
	f, err := os.Open(fName)
	if err != nil {
		return "", err
	}
	defer f.Close()
	hasher := md5.New()
	_, err = io.Copy(hasher, f)
	if err != nil {
		return "", err
	}
	checksum := hasher.Sum(nil)
	return fmt.Sprintf("%x", checksum), nil
}
