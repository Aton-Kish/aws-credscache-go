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

//go:generate go run github.com/golang/mock/mockgen@latest -destination=../mock/github.com/aws/aws-sdk-go-v2/aws/aws.go -package=mock_aws github.com/aws/aws-sdk-go-v2/aws CredentialsProvider

package credscache

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	mock "github.com/Aton-Kish/aws-credscache-go/mock/github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFileCacheProvider_RetrieveCachedCredentials(t *testing.T) {
	expiresIn15Minutes := time.Now().UTC().Add(time.Duration(15) * time.Minute)
	cachedCreds := &aws.Credentials{
		AccessKeyID:     "CachedAccessKeyID",
		SecretAccessKey: "CachedSecretAccessKey",
		SessionToken:    "CachedSessionToken",
		Source:          "TestProvider",
		CanExpire:       true,
		Expires:         expiresIn15Minutes,
	}

	type fields struct {
		cacheKey string
		optFns   []func(o *FileCacheOptions)
	}

	type args struct {
		ctx context.Context
	}

	type mockCredentialsProviderRetrieve struct {
		res aws.Credentials
		err error
	}

	type expected struct {
		res aws.Credentials
		err error
	}

	tests := []struct {
		name                            string
		fields                          fields
		args                            args
		mockCredentialsProviderRetrieve mockCredentialsProviderRetrieve
		expected                        expected
	}{
		{
			name: "positive case",
			fields: fields{
				cacheKey: "cached",
				optFns:   []func(o *FileCacheOptions){},
			},
			mockCredentialsProviderRetrieve: mockCredentialsProviderRetrieve{
				res: aws.Credentials{
					AccessKeyID:     "NonCachedAccessKeyID",
					SecretAccessKey: "NonCachedSecretAccessKey",
					SessionToken:    "NonCachedSessionToken",
					Source:          "TestProvider",
					CanExpire:       true,
					Expires:         expiresIn15Minutes,
				},
				err: nil,
			},
			args: args{
				ctx: context.TODO(),
			},
			expected: expected{
				res: aws.Credentials{
					AccessKeyID:     "CachedAccessKeyID",
					SecretAccessKey: "CachedSecretAccessKey",
					SessionToken:    "CachedSessionToken",
					Source:          "FileCacheProvider",
					CanExpire:       true,
					Expires:         expiresIn15Minutes,
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cachedDir := t.TempDir()
			StoreCredentials(filepath.Join(cachedDir, fmt.Sprintf("%s.json", tt.fields.cacheKey)), cachedCreds)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCredentialsProvider := mock.NewMockCredentialsProvider(ctrl)
			mockCredentialsProvider.
				EXPECT().
				Retrieve(gomock.Any()).
				Return(tt.mockCredentialsProviderRetrieve.res, tt.mockCredentialsProviderRetrieve.err).
				Times(0)

			tt.fields.optFns = append(tt.fields.optFns, func(o *FileCacheOptions) { o.FileCacheDir = cachedDir })

			provider := NewFileCacheProvider(mockCredentialsProvider, tt.fields.cacheKey, tt.fields.optFns...)

			// Act
			actual, err := provider.Retrieve(tt.args.ctx)

			// Assert
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

func TestFileCacheProvider_RetrieveExpiredCredentials(t *testing.T) {
	expiresIn15Minutes := time.Now().UTC().Add(time.Duration(15) * time.Minute)
	expired15MinutesAgo := time.Now().UTC().Add(-time.Duration(15) * time.Minute)
	expiredCreds := &aws.Credentials{
		AccessKeyID:     "CachedAccessKeyID",
		SecretAccessKey: "CachedSecretAccessKey",
		SessionToken:    "CachedSessionToken",
		Source:          "TestProvider",
		CanExpire:       true,
		Expires:         expired15MinutesAgo,
	}

	type fields struct {
		cacheKey string
		optFns   []func(o *FileCacheOptions)
	}

	type args struct {
		ctx context.Context
	}

	type mockCredentialsProviderRetrieve struct {
		res aws.Credentials
		err error
	}

	type expected struct {
		res aws.Credentials
		err error
	}

	tests := []struct {
		name                            string
		fields                          fields
		args                            args
		mockCredentialsProviderRetrieve mockCredentialsProviderRetrieve
		expected                        expected
	}{
		{
			name: "positive case",
			fields: fields{
				cacheKey: "expired",
				optFns:   []func(o *FileCacheOptions){},
			},
			mockCredentialsProviderRetrieve: mockCredentialsProviderRetrieve{
				res: aws.Credentials{
					AccessKeyID:     "NonCachedAccessKeyID",
					SecretAccessKey: "NonCachedSecretAccessKey",
					SessionToken:    "NonCachedSessionToken",
					Source:          "TestProvider",
					CanExpire:       true,
					Expires:         expiresIn15Minutes,
				},
				err: nil,
			},
			args: args{
				ctx: context.TODO(),
			},
			expected: expected{
				res: aws.Credentials{
					AccessKeyID:     "NonCachedAccessKeyID",
					SecretAccessKey: "NonCachedSecretAccessKey",
					SessionToken:    "NonCachedSessionToken",
					Source:          "FileCacheProvider",
					CanExpire:       true,
					Expires:         expiresIn15Minutes,
				},
				err: nil,
			},
		},
		{
			name: "negative case: failed to retrieve",
			fields: fields{
				cacheKey: "expired",
				optFns:   []func(o *FileCacheOptions){},
			},
			args: args{
				ctx: context.TODO(),
			},
			mockCredentialsProviderRetrieve: mockCredentialsProviderRetrieve{
				res: aws.Credentials{Source: "TestProvider"},
				err: errors.New("failed to retrieve"),
			},
			expected: expected{
				res: aws.Credentials{Source: "FileCacheProvider"},
				err: errors.New("failed to retrieve"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cachedDir := t.TempDir()
			StoreCredentials(filepath.Join(cachedDir, fmt.Sprintf("%s.json", tt.fields.cacheKey)), expiredCreds)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCredentialsProvider := mock.NewMockCredentialsProvider(ctrl)
			mockCredentialsProvider.
				EXPECT().
				Retrieve(gomock.Any()).
				Return(tt.mockCredentialsProviderRetrieve.res, tt.mockCredentialsProviderRetrieve.err).
				Times(1)

			tt.fields.optFns = append(tt.fields.optFns, func(o *FileCacheOptions) { o.FileCacheDir = cachedDir })

			provider := NewFileCacheProvider(mockCredentialsProvider, tt.fields.cacheKey, tt.fields.optFns...)

			// Act
			actual, err := provider.Retrieve(tt.args.ctx)

			// Assert
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

func TestFileCacheProvider_RetrieveNonCachedCredentials(t *testing.T) {
	expiresIn15Minutes := time.Now().UTC().Add(time.Duration(15) * time.Minute)

	type fields struct {
		cacheKey string
		optFns   []func(o *FileCacheOptions)
	}

	type args struct {
		ctx context.Context
	}

	type mockCredentialsProviderRetrieve struct {
		res aws.Credentials
		err error
	}

	type expected struct {
		res aws.Credentials
		err error
	}

	tests := []struct {
		name                            string
		fields                          fields
		args                            args
		mockCredentialsProviderRetrieve mockCredentialsProviderRetrieve
		expected                        expected
	}{
		{
			name: "positive case",
			fields: fields{
				cacheKey: "non-cached",
				optFns:   []func(o *FileCacheOptions){},
			},
			mockCredentialsProviderRetrieve: mockCredentialsProviderRetrieve{
				res: aws.Credentials{
					AccessKeyID:     "NonCachedAccessKeyID",
					SecretAccessKey: "NonCachedSecretAccessKey",
					SessionToken:    "NonCachedSessionToken",
					Source:          "TestProvider",
					CanExpire:       true,
					Expires:         expiresIn15Minutes,
				},
				err: nil,
			},
			args: args{
				ctx: context.TODO(),
			},
			expected: expected{
				res: aws.Credentials{
					AccessKeyID:     "NonCachedAccessKeyID",
					SecretAccessKey: "NonCachedSecretAccessKey",
					SessionToken:    "NonCachedSessionToken",
					Source:          "FileCacheProvider",
					CanExpire:       true,
					Expires:         expiresIn15Minutes,
				},
				err: nil,
			},
		},
		{
			name: "negative case: failed to retrieve",
			fields: fields{
				cacheKey: "non-cached",
				optFns:   []func(o *FileCacheOptions){},
			},
			args: args{
				ctx: context.TODO(),
			},
			mockCredentialsProviderRetrieve: mockCredentialsProviderRetrieve{
				res: aws.Credentials{Source: "TestProvider"},
				err: errors.New("failed to retrieve"),
			},
			expected: expected{
				res: aws.Credentials{Source: "FileCacheProvider"},
				err: errors.New("failed to retrieve"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cachedDir := t.TempDir()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCredentialsProvider := mock.NewMockCredentialsProvider(ctrl)
			mockCredentialsProvider.
				EXPECT().
				Retrieve(gomock.Any()).
				Return(tt.mockCredentialsProviderRetrieve.res, tt.mockCredentialsProviderRetrieve.err).
				Times(1)

			tt.fields.optFns = append(tt.fields.optFns, func(o *FileCacheOptions) { o.FileCacheDir = cachedDir })

			provider := NewFileCacheProvider(mockCredentialsProvider, tt.fields.cacheKey, tt.fields.optFns...)

			// Act
			actual, err := provider.Retrieve(tt.args.ctx)

			// Assert
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
