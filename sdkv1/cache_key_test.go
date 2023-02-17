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
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/stretchr/testify/assert"
)

func TestAssumeRoleCacheKey(t *testing.T) {
	type args struct {
		provider *stscreds.AssumeRoleProvider
	}

	type expected struct {
		res string
		err error
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "positive case: with RoleARN",
			args: args{
				provider: &stscreds.AssumeRoleProvider{
					RoleARN: "role_arn",
				},
			},
			expected: expected{
				res: "de1969e7a880d858c9bef3ba110acf78869d4527",
				err: nil,
			},
		},
		{
			name: "positive case: with Duration 15 minutes or less",
			args: args{
				provider: &stscreds.AssumeRoleProvider{
					RoleARN:  "role_arn",
					Duration: time.Duration(15) * time.Minute,
				},
			},
			expected: expected{
				res: "de1969e7a880d858c9bef3ba110acf78869d4527",
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := AssumeRoleCacheKey(tt.args.provider)

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