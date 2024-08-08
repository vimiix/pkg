// Copyright (c) 2024 vimiix
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package file

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestTailN(t *testing.T) {
	tmpdir := t.TempDir()
	tests := []struct {
		content string
		n       int
		expect  []string
	}{
		{"1\n2\n3\n4\n5", 2, []string{"4", "5"}},
		{"1\n2", 2, []string{"1", "2"}},
		{"1\n2", 3, []string{"1", "2"}},
		{"1\n2222\n3\n4444\n5", 2, []string{"4444", "5"}},
	}
	for idx, tt := range tests {
		name := fmt.Sprintf("case_%d", idx)
		t.Run(name, func(t *testing.T) {
			f, err := os.CreateTemp(tmpdir, name)
			if err != nil {
				t.Fatal(err)
			}
			_, err = f.WriteString(tt.content)
			if err != nil {
				t.Fatal(err)
			}
			rs, err := TailN(f.Name(), tt.n)
			if err != nil {
				t.Error(err)
				return
			}
			if !reflect.DeepEqual(rs, tt.expect) {
				t.Errorf("expect %v, actual %v", tt.expect, rs)
			}
		})
	}
}
