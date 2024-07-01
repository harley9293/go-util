package crypto

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// Md5SumFile calculates the MD5 checksum by reading the file in chunks, supporting the calculation of MD5 for large files.
func Md5SumFile(file string) (value string, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}

	d := md5.New()
	r := bufio.NewReader(f)
	buf := make([]byte, 512)
	for {
		var n int
		n, err = r.Read(buf)
		if err != nil && err != io.EOF {
			return
		}

		if n == 0 {
			break
		}

		_, err = d.Write(buf[:n])
		if err != nil {
			return
		}
	}

	if err == io.EOF {
		err = nil
	}

	value = fmt.Sprintf("%x", d.Sum([]byte{}))
	return
}
