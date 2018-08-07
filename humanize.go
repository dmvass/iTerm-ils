package main

import (
	"fmt"
	"math"
	"os"
)

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanateBytes(s int64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f%s"
	if val < 10 {
		f = "%.1f%s"
	}

	return fmt.Sprintf(f, val, suffix)
}

// bytes produces a human readable representation of an EC size.
// bytes(82854982) -> 79 MiB
func bytes(s int64) string {
	sizes := []string{"B", "K", "M", "G", "T", "P", "E"}
	return humanateBytes(s, 1000, sizes)
}

// permissions produces a human readable representation of file mode
// permissions(420) -> rw-r--r--
func permissions(m os.FileMode) string {
	var buf [32]byte // Mode is uint32.
	w := 0

	const rwx = "rwxrwxrwx"
	for i, c := range rwx {
		if m&(1<<uint(9-1-i)) != 0 {
			buf[w] = byte(c)
		} else {
			buf[w] = '-'
		}
		w++
	}
	return string(buf[:w])
}
