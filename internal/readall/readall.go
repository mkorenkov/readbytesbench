package readall

import (
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

// DoWork reads all bytes from the given reader and calculates their hash very naive, but extremely predictable way
func DoWork(src io.Reader) (int, error) {
	data, err := ioutil.ReadAll(src)
	if err != nil {
		return 0, errors.Wrap(err, "Error reading bytes")
	}
	var res int
	for i := range data {
		if data[i] == 0 {
			res++
		}
	}
	return res, nil
}
