# readbytesbench

If you are writing a code for a production system, use `io.Reader` unless you can prove why you shouldn’t
```
BenchmarkRead/ioutil.ReadAll-nonbuf-20M-throttle-16         	       1	2187802238 ns/op	   9.14 MB/s	        64.0 TotalAlloc(MiB)	67109128 B/op	      30 allocs/op
BenchmarkRead/ioutil.ReadAll-buf-20M-throttle-16            	       1	2114013952 ns/op	   9.46 MB/s	        64.0 TotalAlloc(MiB)	67109872 B/op	      36 allocs/op
BenchmarkRead/io.Reader-nonbuf-20M-throttle-16              	       1	2008726980 ns/op	   9.96 MB/s	       0.000404 TotalAlloc(MiB)	     424 B/op	       9 allocs/op
BenchmarkRead/io.Reader-buf-20M-throttle-16                 	       1	2020003524 ns/op	   9.90 MB/s	       0.000969 TotalAlloc(MiB)	    1016 B/op	      10 allocs/op
```
I implemented something akin 10MB/s network connection, that acts as a bottleneck and used 20M from /dev/zero as a payload

However, if you want to win a pissing contest, there is a use case for `ioutil.ReadAll` :
```
BenchmarkRead/ioutil.ReadAll-nonbuf-20M-nothrottle-16       	1000000000	         0.0770 ns/op	259635916448.90 MB/s	        64.0 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
BenchmarkRead/ioutil.ReadAll-buf-20M-nothrottle-16          	1000000000	         0.0763 ns/op	262160171764.20 MB/s	        64.0 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
BenchmarkRead/io.Reader-nonbuf-20M-nothrottle-16            	1000000000	         0.486 ns/op	41123292776.27 MB/s	          0.000206 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
BenchmarkRead/io.Reader-buf-20M-nothrottle-16               	1000000000	         0.115 ns/op	174173742191.91 MB/s	      0.000786 TotalAlloc(MiB)	       0 B/op	       0 allocs/op
```
(this is how the same code behaves with no io bottleneck)
