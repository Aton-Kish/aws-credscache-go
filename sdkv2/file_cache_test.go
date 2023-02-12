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

package credscache

import (
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

func TestLoadCredentials(t *testing.T) {
	tempDir := t.TempDir()
	creds := &aws.Credentials{
		AccessKeyID:     "AccessKeyID",
		SecretAccessKey: "SecretAccessKey",
		SessionToken:    "SessionToken",
		Source:          "TestProvider",
		CanExpire:       true,
		Expires:         time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
	}
	StoreCredentials(filepath.Join(tempDir, "cache.json"), creds)

	type args struct {
		path string
	}

	type expected struct {
		res *aws.Credentials
		err error
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "positive case: existing cache",
			args: args{
				path: filepath.Join(tempDir, "cache.json"),
			},
			expected: expected{
				res: &aws.Credentials{
					AccessKeyID:     "AccessKeyID",
					SecretAccessKey: "SecretAccessKey",
					SessionToken:    "SessionToken",
					Source:          "",
					CanExpire:       true,
					Expires:         time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
				},
				err: nil,
			},
		},
		{
			name: "negative case: no such file",
			args: args{
				path: filepath.Join(tempDir, "non-existing.json"),
			},
			expected: expected{
				res: nil,
				err: syscall.Errno(2),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := LoadCredentials(tt.args.path)

			if tt.expected.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.res, actual)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expected.err)
			}
		})
	}
}

func TestStoreCredentials(t *testing.T) {
	type args struct {
		path  string
		creds *aws.Credentials
	}

	type expected struct {
		err error
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "positive case: existing dir",
			args: args{
				path: filepath.Join(t.TempDir(), "cache.json"),
				creds: &aws.Credentials{
					AccessKeyID:     "AccessKeyID",
					SecretAccessKey: "SecretAccessKey",
					SessionToken:    "SessionToken",
					Source:          "TestProvider",
					CanExpire:       true,
					Expires:         time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
				},
			},
			expected: expected{
				err: nil,
			},
		},
		{
			name: "positive case: non-existing dir",
			args: args{
				path: filepath.Join(t.TempDir(), "non-existing/cache.json"),
				creds: &aws.Credentials{
					AccessKeyID:     "AccessKeyID",
					SecretAccessKey: "SecretAccessKey",
					SessionToken:    "SessionToken",
					Source:          "TestProvider",
					CanExpire:       true,
					Expires:         time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
				},
			},
			expected: expected{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StoreCredentials(tt.args.path, tt.args.creds)

			if tt.expected.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expected.err)
			}
		})
	}
}
