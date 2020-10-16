package reader

import (
	"io"

	"github.com/pkg/errors"
)

const initialBufSize = 64

// DoWork reads bytes from the given reader and calculates their hash very naive, but extremely predictable way
func DoWork(src io.Reader) (int, error) {
	var res int
	buf := make([]byte, initialBufSize)
	for {
		n, err := src.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, errors.Wrap(err, "reading error")
		}
		for i := 0; i < n; i++ {
			if buf[i] == 0 {
				res++
			}
		}
	}
	return res, nil
}
