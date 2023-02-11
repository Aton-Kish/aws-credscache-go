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
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

func TestStoreCredentialsAndLoadCredentials(t *testing.T) {
	type args struct {
		path  string
		creds *aws.Credentials
	}

	type expected struct {
		loadCreds *aws.Credentials
		loadErr   error
		storeErr  error
	}

	tests := []struct {
		name string
		args args
		expected
	}{
		{
			name: "positive case: existing dir",
			args: args{
				path: filepath.Join(t.TempDir(), "existing-dir.json"),
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
				loadCreds: &aws.Credentials{
					AccessKeyID:     "AccessKeyID",
					SecretAccessKey: "SecretAccessKey",
					SessionToken:    "SessionToken",
					Source:          "",
					CanExpire:       true,
					Expires:         time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
				},
				loadErr:  nil,
				storeErr: nil,
			},
		},
		{
			name: "positive case: non-existing dir",
			args: args{
				path: filepath.Join(t.TempDir(), "non-existing/non-existing-dir.json"),
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
				loadCreds: &aws.Credentials{
					AccessKeyID:     "AccessKeyID",
					SecretAccessKey: "SecretAccessKey",
					SessionToken:    "SessionToken",
					Source:          "",
					CanExpire:       true,
					Expires:         time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
				},
				loadErr:  nil,
				storeErr: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storeErr := StoreCredentials(tt.args.path, tt.args.creds)
			loadCreds, loadErr := LoadCredentials(tt.args.path)

			if tt.expected.storeErr == nil {
				assert.NoError(t, storeErr)
			} else {
				assert.Error(t, storeErr)
				assert.Equal(t, tt.expected.storeErr, storeErr)
			}
			if tt.expected.loadErr == nil {
				assert.NoError(t, loadErr)
				assert.Equal(t, tt.expected.loadCreds, loadCreds)
			} else {
				assert.Error(t, loadErr)
				assert.Equal(t, tt.expected.loadErr, loadErr)
			}
		})
	}
}
