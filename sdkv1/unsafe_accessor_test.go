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

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/stretchr/testify/assert"
)

func TestNewCredentialsUnsafeAccessor(t *testing.T) {
	type args struct {
		ptr *credentials.Credentials
	}

	type expected struct {
		res *CredentialsUnsafeAccessor
		err error
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "positive case: Credentials with valid provider",
			args: args{
				ptr: credentials.NewCredentials(&stscreds.AssumeRoleProvider{}),
			},
			expected: expected{
				res: &CredentialsUnsafeAccessor{ptr: credentials.NewCredentials(&stscreds.AssumeRoleProvider{})},
				err: nil,
			},
		},
		{
			name: "positive case: Credentials with nil provider",
			args: args{
				ptr: credentials.NewCredentials(nil),
			},
			expected: expected{
				res: &CredentialsUnsafeAccessor{ptr: credentials.NewCredentials(nil)},
				err: nil,
			},
		},
		{
			name: "negative case: nil Credentials",
			args: args{
				ptr: nil,
			},
			expected: expected{
				res: nil,
				err: ErrNilPointer,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := NewCredentialsUnsafeAccessor(tt.args.ptr)

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

func TestCredentialsUnsafeAccessor_Provider(t *testing.T) {
	type expected struct {
		res credentials.Provider
	}

	tests := []struct {
		name     string
		accessor *CredentialsUnsafeAccessor
		expected expected
	}{
		{
			name:     "positive case: get valid provider",
			accessor: &CredentialsUnsafeAccessor{ptr: credentials.NewCredentials(&stscreds.AssumeRoleProvider{})},
			expected: expected{
				res: &stscreds.AssumeRoleProvider{},
			},
		},
		{
			name:     "positive case: get nil provider",
			accessor: &CredentialsUnsafeAccessor{ptr: credentials.NewCredentials(nil)},
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

func TestCredentialsUnsafeAccessor_SetProvider(t *testing.T) {
	tests := []struct {
		name     string
		accessor *CredentialsUnsafeAccessor
		provider credentials.Provider
	}{
		{
			name:     "positive case: set valid provider",
			accessor: &CredentialsUnsafeAccessor{ptr: credentials.NewCredentials(&stscreds.AssumeRoleProvider{})},
			provider: &credentials.StaticProvider{},
		},
		{
			name:     "positive case: set nil provider",
			accessor: &CredentialsUnsafeAccessor{ptr: credentials.NewCredentials(&stscreds.AssumeRoleProvider{})},
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
