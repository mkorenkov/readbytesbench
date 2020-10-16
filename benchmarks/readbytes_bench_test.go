package bench

import (
	"bufio"
	"io"
	"log"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/aybabtme/iocontrol"
	"github.com/mkorenkov/readbytesbench/internal/readall"
	"github.com/mkorenkov/readbytesbench/internal/reader"
	"github.com/pkg/errors"
)

const K = 1000
const M = 1000 * K
const blobSource = "/dev/urandom"
const bufSize = 512

type testFunction func(bytesToRead int64, throttler func(io.Reader) io.Reader, worker func(io.Reader) (int, error)) (int, error)

func throttle(r io.Reader) io.Reader {
	return iocontrol.ThrottledReader(r, 10*M, 10*time.Millisecond)
}

func nothrottle(r io.Reader) io.Reader {
	return r
}

func buffered(bytesToRead int64, throttler func(io.Reader) io.Reader, worker func(io.Reader) (int, error)) (int, error) {
	f, err := os.OpenFile(blobSource, os.O_RDONLY, 0755)
	if err != nil {
		return 0, errors.Wrapf(err, "Error opening %s", blobSource)
	}
	defer func() {
		if dErr := f.Close(); dErr != nil {
			panic(dErr)
		}
	}()
	bufReader := bufio.NewReaderSize(throttler(f), bufSize)
	limitReader := io.LimitReader(bufReader, bytesToRead)
	return worker(limitReader)
}

func nonBuffered(bytesToRead int64, throttler func(io.Reader) io.Reader, worker func(io.Reader) (int, error)) (int, error) {
	f, err := os.OpenFile(blobSource, os.O_RDONLY, 0755)
	if err != nil {
		return 0, errors.Wrapf(err, "Error opening %s", blobSource)
	}
	defer func() {
		if dErr := f.Close(); dErr != nil {
			panic(dErr)
		}
	}()
	limitReader := io.LimitReader(throttler(f), bytesToRead)
	return worker(limitReader)
}

func BenchmarkRead(b *testing.B) {
	testCases := []struct {
		name                string
		size                int64
		testFunction        testFunction
		implemenationDetail func(io.Reader) (int, error)
		throttler           func(io.Reader) io.Reader
	}{
		{"ioutil.ReadAll-nonbuf-20M-throttle", 20 * M, nonBuffered, readall.DoWork, throttle},
		{"ioutil.ReadAll-buf-20M-throttle", 20 * M, buffered, readall.DoWork, throttle},
		{"io.Reader-nonbuf-20M-throttle", 20 * M, nonBuffered, reader.DoWork, throttle},
		{"io.Reader-buf-20M-throttle", 20 * M, buffered, reader.DoWork, throttle},
		{"ioutil.ReadAll-nonbuf-20M-nothrottle", 20 * M, nonBuffered, readall.DoWork, nothrottle},
		{"ioutil.ReadAll-buf-20M-nothrottle", 20 * M, buffered, readall.DoWork, nothrottle},
		{"io.Reader-nonbuf-20M-nothrottle", 20 * M, nonBuffered, reader.DoWork, nothrottle},
		{"io.Reader-buf-20M-nothrottle", 20 * M, buffered, reader.DoWork, nothrottle},
	}

	for _, testCase := range testCases {
		b.Run(testCase.name, func(b *testing.B) {
			var start, end runtime.MemStats
			runtime.GC()

			runtime.ReadMemStats(&start)
			b.SetBytes(testCase.size)
			b.ResetTimer()

			if _, err := testCase.testFunction(testCase.size, testCase.throttler, testCase.implemenationDetail); err != nil {
				log.Fatal(err)
			}

			b.ReportAllocs()
			runtime.ReadMemStats(&end)
			b.ReportMetric(float64(end.TotalAlloc-start.TotalAlloc)/1024/1024, "TotalAlloc(MiB)")
		})
		b.StopTimer()
	}
}
