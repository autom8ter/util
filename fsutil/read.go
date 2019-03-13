package fsutil

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

// consistentReadSync is the main functionality of ConsistentRead but
// introduces a sync callback that can be used by the tests to mutate the file
// from which the test data is being read
func ConsistentRead(filename string, attempts int, sync func(int)) ([]byte, error) {
	oldContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	for i := 0; i < attempts; i++ {
		if sync != nil {
			sync(i)
		}
		newContent, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		if bytes.Compare(oldContent, newContent) == 0 {
			return newContent, nil
		}
		// Files are different, continue reading
		oldContent = newContent
	}
	return nil, fmt.Errorf("could not get consistent content of %s after %d attempts", filename, attempts)
}
