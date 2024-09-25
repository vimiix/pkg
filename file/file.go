// Copyright (c) 2024 vimiix
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package file

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"os/user"
	"path/filepath"
)

// ExpandHomePath returns a path with a leading ~ replaced with the
// user's home directory. If the path does not start with a ~, or if
// the user's home directory cannot be determined, the original path
// is returned unchanged.
func ExpandHomePath(s string) string {
	if s == "" || s[0] != '~' {
		return s
	}
	u, err := user.Current()
	if err != nil {
		return s
	}
	return u.HomeDir + s[1:]
}

// Exists checks whether a file exists.
// It returns false if the given path is empty.
func Exists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// IsDir checks whether a path is a directory.
// It returns false if an error occurs.
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsSymLink checks whether a path is a symbolic link.
// It returns false if an error occurs.
func IsSymLink(path string) bool {
	s, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return s.Mode()&os.ModeSymlink == os.ModeSymlink
}

// EnsureDirExists ensures that the parent directory of the given path exists
// and is writable, by creating it if necessary.
func EnsureDirExists(path string) error {
	parent := filepath.Dir(ExpandHomePath(path))
	return os.MkdirAll(parent, os.ModePerm)
}

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
