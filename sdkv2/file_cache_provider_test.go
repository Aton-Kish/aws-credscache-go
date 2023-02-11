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
	"reflect"
	"testing"
	"time"

	mock "github.com/Aton-Kish/aws-credscache-go/mock/github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFileCacheProvider_Retrieve(t *testing.T) {
	type fields struct {
		provider aws.CredentialsProvider
		cacheKey string
		options  FileCacheOptions
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    aws.Credentials
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &FileCacheProvider{
				provider: tt.fields.provider,
				cacheKey: tt.fields.cacheKey,
				options:  tt.fields.options,
			}
			got, err := p.Retrieve(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileCacheProvider.Retrieve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FileCacheProvider.Retrieve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileCacheProvider(t *testing.T) {
	cachedDir := t.TempDir()
	cachedCacheKey := "cached"
	cachedExpires := time.Now().UTC().Add(time.Duration(24) * time.Hour)
	assert.NoError(t, StoreCredentials(
		filepath.Join(cachedDir, fmt.Sprintf("%s.json", cachedCacheKey)),
		&aws.Credentials{
			AccessKeyID:     "CachedAccessKeyID",
			SecretAccessKey: "SecretAccessKey",
			SessionToken:    "SessionToken",
			Source:          "TestProvider",
			CanExpire:       true,
			Expires:         cachedExpires,
		}))

	expiredCacheKey := "expired"
	assert.NoError(t, StoreCredentials(
		filepath.Join(cachedDir, fmt.Sprintf("%s.json", expiredCacheKey)),
		&aws.Credentials{
			AccessKeyID:     "ExpiredAccessKeyID",
			SecretAccessKey: "SecretAccessKey",
			SessionToken:    "SessionToken",
			Source:          "TestProvider",
			CanExpire:       true,
			Expires:         time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
		}))

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
			name: "positive case: cached",
			fields: fields{
				cacheKey: cachedCacheKey,
				optFns: []func(o *FileCacheOptions){
					func(o *FileCacheOptions) {
						o.FileCacheDir = cachedDir
					},
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			mockCredentialsProviderRetrieve: mockCredentialsProviderRetrieve{
				res: aws.Credentials{
					AccessKeyID:     "NonCachedAccessKeyID",
					SecretAccessKey: "SecretAccessKey",
					SessionToken:    "SessionToken",
					Source:          "TestProvider",
					CanExpire:       true,
					Expires:         time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
				},
				err: nil,
			},
			expected: expected{
				res: aws.Credentials{
					AccessKeyID:     "CachedAccessKeyID",
					SecretAccessKey: "SecretAccessKey",
					SessionToken:    "SessionToken",
					Source:          "FileCacheProvider",
					CanExpire:       true,
					Expires:         cachedExpires,
				},
				err: nil,
			},
		},
		{
			name: "positive case: cached but expired",
			fields: fields{
				cacheKey: expiredCacheKey,
				optFns: []func(o *FileCacheOptions){
					func(o *FileCacheOptions) {
						o.FileCacheDir = cachedDir
					},
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			mockCredentialsProviderRetrieve: mockCredentialsProviderRetrieve{
				res: aws.Credentials{
					AccessKeyID:     "NonCachedAccessKeyID",
					SecretAccessKey: "SecretAccessKey",
					SessionToken:    "SessionToken",
					Source:          "TestProvider",
					CanExpire:       true,
					Expires:         time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
				},
				err: nil,
			},
			expected: expected{
				res: aws.Credentials{
					AccessKeyID:     "NonCachedAccessKeyID",
					SecretAccessKey: "SecretAccessKey",
					SessionToken:    "SessionToken",
					Source:          "FileCacheProvider",
					CanExpire:       true,
					Expires:         time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
				},
				err: nil,
			},
		},
		{
			name: "positive case: non-cached",
			fields: fields{
				cacheKey: "non-cached",
				optFns: []func(o *FileCacheOptions){
					func(o *FileCacheOptions) {
						o.FileCacheDir = t.TempDir()
					},
				},
			},
			args: args{
				ctx: context.TODO(),
			},
			mockCredentialsProviderRetrieve: mockCredentialsProviderRetrieve{
				res: aws.Credentials{
					AccessKeyID:     "NonCachedAccessKeyID",
					SecretAccessKey: "SecretAccessKey",
					SessionToken:    "SessionToken",
					Source:          "TestProvider",
					CanExpire:       true,
					Expires:         time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
				},
				err: nil,
			},
			expected: expected{
				res: aws.Credentials{
					AccessKeyID:     "NonCachedAccessKeyID",
					SecretAccessKey: "SecretAccessKey",
					SessionToken:    "SessionToken",
					Source:          "FileCacheProvider",
					CanExpire:       true,
					Expires:         time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
				},
				err: nil,
			},
		},
		{
			name: "negative case: failed to retrieve",
			fields: fields{
				cacheKey: "non-cached",
				optFns: []func(o *FileCacheOptions){
					func(o *FileCacheOptions) {
						o.FileCacheDir = t.TempDir()
					},
				},
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockCredentialsProvider := mock.NewMockCredentialsProvider(ctrl)
			mockCredentialsProvider.
				EXPECT().
				Retrieve(gomock.Any()).
				Return(tt.mockCredentialsProviderRetrieve.res, tt.mockCredentialsProviderRetrieve.err).
				AnyTimes()
			provider := NewFileCacheProvider(mockCredentialsProvider, tt.fields.cacheKey, tt.fields.optFns...)

			actual, err := provider.Retrieve(tt.args.ctx)

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
