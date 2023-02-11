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

//go:generate go run github.com/golang/mock/mockgen@latest -destination=../mock/github.com/aws/aws-sdk-go-v2/credentials/stscreds/stscreds.go -package=mock_stscreds github.com/aws/aws-sdk-go-v2/credentials/stscreds AssumeRoleAPIClient

package credscache

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	mock "github.com/Aton-Kish/aws-credscache-go/mock/github.com/aws/aws-sdk-go-v2/credentials/stscreds"
)

func TestNewCredentialsCacheUnsafeAccessor(t *testing.T) {
	type args struct {
		credsCache *aws.CredentialsCache
	}

	type expected struct {
		res CredentialsCacheUnsafeAccessor
		err error
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "positive case: CredentialsCache with valid provider",
			args: args{
				credsCache: aws.NewCredentialsCache(&stscreds.AssumeRoleProvider{}),
			},
			expected: expected{
				res: &credentialsCacheUnsafeAccessor{ptr: aws.NewCredentialsCache(&stscreds.AssumeRoleProvider{})},
				err: nil,
			},
		},
		{
			name: "positive case: CredentialsCache with nil provider",
			args: args{
				credsCache: aws.NewCredentialsCache(nil),
			},
			expected: expected{
				res: &credentialsCacheUnsafeAccessor{ptr: aws.NewCredentialsCache(nil)},
				err: nil,
			},
		},
		{
			name: "negative case: nil CredentialsCache",
			args: args{
				credsCache: nil,
			},
			expected: expected{
				res: nil,
				err: ErrNilPointer,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := NewCredentialsCacheUnsafeAccessor(tt.args.credsCache)

			if tt.expected.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.res, actual)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expected.err, err)
			}
		})
	}
}

func Test_credentialsCacheUnsafeAccessor_Provider(t *testing.T) {
	type expected struct {
		res aws.CredentialsProvider
	}

	tests := []struct {
		name     string
		accessor *credentialsCacheUnsafeAccessor
		expected expected
	}{
		{
			name:     "positive case: get valid provider",
			accessor: &credentialsCacheUnsafeAccessor{ptr: aws.NewCredentialsCache(&stscreds.AssumeRoleProvider{})},
			expected: expected{
				res: &stscreds.AssumeRoleProvider{},
			},
		},
		{
			name:     "positive case: get nil provider",
			accessor: &credentialsCacheUnsafeAccessor{ptr: aws.NewCredentialsCache(nil)},
			expected: expected{
				res: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.accessor.Provider()

			assert.Equal(t, tt.expected.res, actual)
		})
	}
}

func Test_credentialsCacheUnsafeAccessor_SetProvider(t *testing.T) {
	tests := []struct {
		name     string
		accessor *credentialsCacheUnsafeAccessor
		provider aws.CredentialsProvider
	}{
		{
			name:     "positive case: set valid provider",
			accessor: &credentialsCacheUnsafeAccessor{ptr: aws.NewCredentialsCache(&stscreds.AssumeRoleProvider{})},
			provider: &credentials.StaticCredentialsProvider{},
		},
		{
			name:     "positive case: set nil provider",
			accessor: &credentialsCacheUnsafeAccessor{ptr: aws.NewCredentialsCache(&stscreds.AssumeRoleProvider{})},
			provider: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.accessor.SetProvider(tt.provider)

			assert.Equal(t, tt.provider, tt.accessor.Provider())
		})
	}
}

func TestNewAssumeRoleProviderUnsafeAccessor(t *testing.T) {
	type args struct {
		provider *stscreds.AssumeRoleProvider
	}

	type expected struct {
		res AssumeRoleProviderUnsafeAccessor
		err error
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "positive case: AssumeRoleProvider",
			args: args{
				provider: &stscreds.AssumeRoleProvider{},
			},
			expected: expected{
				res: &assumeRoleProviderUnsafeAccessor{ptr: &stscreds.AssumeRoleProvider{}},
				err: nil,
			},
		},
		{
			name: "negative case: nil AssumeRoleProvider",
			args: args{
				provider: nil,
			},
			expected: expected{
				res: nil,
				err: ErrNilPointer,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := NewAssumeRoleProviderUnsafeAccessor(tt.args.provider)

			if tt.expected.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.res, actual)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expected.err, err)
			}
		})
	}
}

func Test_assumeRoleProviderUnsafeAccessor_Options(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAssumeRoleAPIClient := mock.NewMockAssumeRoleAPIClient(ctrl)

	type expected struct {
		res stscreds.AssumeRoleOptions
	}

	tests := []struct {
		name     string
		accessor *assumeRoleProviderUnsafeAccessor
		expected expected
	}{
		{
			name:     "positive case: get basic option",
			accessor: &assumeRoleProviderUnsafeAccessor{ptr: stscreds.NewAssumeRoleProvider(mockAssumeRoleAPIClient, "role_arn")},
			expected: expected{
				res: stscreds.AssumeRoleOptions{
					Client:  mockAssumeRoleAPIClient,
					RoleARN: "role_arn",
				},
			},
		},
		{
			name: "positive case: get option with role session name",
			accessor: &assumeRoleProviderUnsafeAccessor{ptr: stscreds.NewAssumeRoleProvider(mockAssumeRoleAPIClient, "role_arn", func(o *stscreds.AssumeRoleOptions) {
				o.RoleSessionName = "role_session_name"
			})},
			expected: expected{
				res: stscreds.AssumeRoleOptions{
					Client:          mockAssumeRoleAPIClient,
					RoleARN:         "role_arn",
					RoleSessionName: "role_session_name",
				},
			},
		},
		{
			name:     "positive case: get empty option",
			accessor: &assumeRoleProviderUnsafeAccessor{ptr: &stscreds.AssumeRoleProvider{}},
			expected: expected{
				res: stscreds.AssumeRoleOptions{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.accessor.Options()

			assert.Equal(t, tt.expected.res, actual)
		})
	}
}
