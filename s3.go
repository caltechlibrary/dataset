package dataset

// From https://github.com/aws/aws-sdk-go README.md
import (
	"fmt"
	"time"
)

/*
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
*/

// List send a request to bring back a list of bucket/object names
func List(bucketName, objectPath string, timeout *time.Duration) ([]string, error) {
	return nil, fmt.Errorf("List(%s, %s, %v) not implemented", bucketName, objectName, timeout)
}

// Upload send the contents of buf to AWS S3 bucketed with Object Name provided
func Upload(bucketName, objectName string, timeout *time.Duration, buf []byte) error {
	//NOTE: need to handle streaming multi-part upload/download
	return fmt.Errorf("Upload(%s, %s, %v, %+v) not implemented", bucketName, objectName, timeout, buf)
}

// Download request the contents of a bucket from AWS S3 returning buffer and error
func Download(bucketName, objectName string, timeout *time.Duration) ([]byte, error) {
	//NOTE: need to handle streaming multi-part upload/download
	return nil, fmt.Errorf("Download(%s, %s, %v) not implemented", bucketName, objectName, timeout)
}

// Delete send a delete request to AWS S3
func Delete(bucketName, objectName string, timeout *time.Duration) error {
	return fmt.Errorf("Delete(%s, %s, %v) not implemented", bucketName, objectName, timeout)
}
