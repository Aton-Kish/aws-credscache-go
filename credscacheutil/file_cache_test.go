// Copyright (c) 2023 Aton-Kish
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package credscacheutil

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFileCache_Load(t *testing.T) {
	tempDir := t.TempDir()
	cache := &FileCache{
		Credentials: CachedCredentials{
			AccessKeyID:     "AccessKeyID",
			SecretAccessKey: "SecretAccessKey",
			SessionToken:    "SessionToken",
			Expires:         time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
		},
	}
	cache.Store(filepath.Join(tempDir, "cache.json"))

	type args struct {
		path string
	}

	type expected struct {
		cache *FileCache
		err   error
	}

	tests := []struct {
		name     string
		cache    *FileCache
		args     args
		expected expected
	}{
		{
			name:  "positive case: existing cache",
			cache: new(FileCache),
			args: args{
				path: filepath.Join(tempDir, "cache.json"),
			},
			expected: expected{
				cache: &FileCache{
					Credentials: CachedCredentials{
						AccessKeyID:     "AccessKeyID",
						SecretAccessKey: "SecretAccessKey",
						SessionToken:    "SessionToken",
						Expires:         time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
					},
				},
				err: nil,
			},
		},
		{
			name:  "negative case: no such file",
			cache: new(FileCache),
			args: args{
				path: filepath.Join(tempDir, "non-existing.json"),
			},
			expected: expected{
				cache: nil,
				err:   &os.PathError{Op: "open", Path: filepath.Join(tempDir, "non-existing.json"), Err: syscall.Errno(2)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cache.Load(tt.args.path)

			if tt.expected.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.cache, tt.cache)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expected.err, err)
			}
		})
	}
}

func TestFileCache_Store(t *testing.T) {
	type args struct {
		path string
	}

	type expected struct {
		err error
	}

	tests := []struct {
		name     string
		cache    *FileCache
		args     args
		expected expected
	}{
		{
			name: "positive case: existing dir",
			cache: &FileCache{
				Credentials: CachedCredentials{
					AccessKeyID:     "AccessKeyID",
					SecretAccessKey: "SecretAccessKey",
					SessionToken:    "SessionToken",
					Expires:         time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
				},
			},
			args: args{
				path: filepath.Join(t.TempDir(), "cache.json"),
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "positive case: non-existing dir",
			cache: &FileCache{
				Credentials: CachedCredentials{
					AccessKeyID:     "AccessKeyID",
					SecretAccessKey: "SecretAccessKey",
					SessionToken:    "SessionToken",
					Expires:         time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
				},
			},
			args: args{
				path: filepath.Join(t.TempDir(), "non-existing/cache.json"),
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cache.Store(tt.args.path)

			if tt.expected.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expected.err, err)
			}
		})
	}
}
