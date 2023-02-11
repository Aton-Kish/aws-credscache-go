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
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFileCache_StoreAndLoad(t *testing.T) {
	type args struct {
		path string
	}

	type expected struct {
		loadCache *FileCache
		loadErr   error
		storeErr  error
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
				path: filepath.Join(t.TempDir(), "existing-dir.json"),
			},
			expected: expected{
				loadCache: &FileCache{
					Credentials: CachedCredentials{
						AccessKeyID:     "AccessKeyID",
						SecretAccessKey: "SecretAccessKey",
						SessionToken:    "SessionToken",
						Expires:         time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
					},
				},
				loadErr:  nil,
				storeErr: nil,
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
				path: filepath.Join(t.TempDir(), "non-existing/non-existing-dir.json"),
			},
			expected: expected{
				loadCache: &FileCache{
					Credentials: CachedCredentials{
						AccessKeyID:     "AccessKeyID",
						SecretAccessKey: "SecretAccessKey",
						SessionToken:    "SessionToken",
						Expires:         time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
					},
				},
				loadErr:  nil,
				storeErr: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storeErr := tt.cache.Store(tt.args.path)
			loadCache := new(FileCache)
			loadErr := loadCache.Load(tt.args.path)

			if tt.expected.storeErr == nil {
				assert.NoError(t, storeErr)
			} else {
				assert.Error(t, storeErr)
				assert.Equal(t, tt.expected.storeErr, storeErr)
			}
			if tt.expected.loadErr == nil {
				assert.NoError(t, loadErr)
				assert.Equal(t, tt.expected.loadCache, loadCache)
			} else {
				assert.Error(t, loadErr)
				assert.Equal(t, tt.expected.loadErr, loadErr)
			}
		})
	}
}
