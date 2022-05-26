package dataset

import (
	"fmt"
)

func (c *Collection) Clone(cloneName string, keys []string, verbose bool) error {
	return fmt.Errorf("Clone() not implemented")
}

func (c *Collection) CloneSample(trainingName string, testName string, keys []string, trainingSetSize int, verbose bool) error {
	return fmt.Errorf("CloneSample() not implemented")
}
