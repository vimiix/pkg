// Copyright (c) 2024 vimiix
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package file

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

// TailN returns last n lines of file
func TailN(path string, n int) (rs []string, err error) {
	var f *os.File
	f, err = os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	var (
		fileSize      = stat.Size()
		bufferSize    = int64(1024)
		offset        = bufferSize
		readSize      int64
		totalReadSize int64
		finalBytes    []byte
		lines         []string
	)

	for {
		if offset > fileSize {
			offset = fileSize
		}
		readSize = offset - totalReadSize

		if _, err := f.Seek(-offset, io.SeekEnd); err != nil {
			return nil, err
		}

		buffer := make([]byte, readSize)
		readBytes, err := f.Read(buffer)
		if err != nil {
			return nil, err
		}
		totalReadSize += int64(readBytes)
		finalBytes = append(buffer[:readBytes], finalBytes...)
		lines = splitLines(finalBytes)
		if len(lines) > n || offset == fileSize {
			break
		}

		offset += bufferSize
	}

	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return lines, nil
}

func splitLines(data []byte) []string {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
