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

//go:generate go run github.com/golang/mock/mockgen@latest -destination=../internal/mock/github.com/aws/aws-sdk-go-v2/aws/aws.go -package=mock_aws github.com/aws/aws-sdk-go-v2/aws CredentialsProvider

package credscache

import (
	"testing"

	mock "github.com/Aton-Kish/aws-credscache-go/internal/mock/github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInjectFileCacheProvider(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCredentialsProvider := mock.NewMockCredentialsProvider(ctrl)

	type args struct {
		cfg    *aws.Config
		optFns []func(o *FileCacheOptions)
	}

	type expected struct {
		res bool
		err error
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "positive case: succeeded to inject",
			args: args{
				cfg: &aws.Config{
					Credentials: aws.NewCredentialsCache(&stscreds.AssumeRoleProvider{}),
				},
				optFns: []func(o *FileCacheOptions){},
			},
			expected: expected{
				res: true,
				err: nil,
			},
		},
		{
			name: "positive case: failed to inject due to missing CredentialsCache",
			args: args{
				cfg: &aws.Config{
					Credentials: mockCredentialsProvider,
				},
				optFns: []func(o *FileCacheOptions){},
			},
			expected: expected{
				res: false,
				err: nil,
			},
		},
		{
			name: "positive case: failed to inject due to missing AssumeRoleProvider",
			args: args{
				cfg: &aws.Config{
					Credentials: aws.NewCredentialsCache(mockCredentialsProvider),
				},
				optFns: []func(o *FileCacheOptions){},
			},
			expected: expected{
				res: false,
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := InjectFileCacheProvider(tt.args.cfg, tt.args.optFns...)

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
