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

//go:generate go run github.com/golang/mock/mockgen@latest -destination=../internal/mock/github.com/aws/aws-sdk-go/aws/credentials/credentials.go -package=mock_credentials github.com/aws/aws-sdk-go/aws/credentials ProviderWithContext,Expirer

package credscache

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	mock_credscache "github.com/Aton-Kish/aws-credscache-go/internal/mock/github.com/Aton-Kish/aws-credscache-go/sdkv1"
	mock_credentials "github.com/Aton-Kish/aws-credscache-go/internal/mock/github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFileCacheProvider_RetrieveWithExpireProvider(t *testing.T) {
	expiresIn15Minutes := time.Now().UTC().Add(time.Duration(15) * time.Minute)
	expired15MinutesAgo := time.Now().UTC().Add(-time.Duration(15) * time.Minute)
	cachedCreds := &credentials.Value{
		AccessKeyID:     "CachedAccessKeyID",
		SecretAccessKey: "CachedSecretAccessKey",
		SessionToken:    "CachedSessionToken",
		ProviderName:    "TestProvider",
	}

	errRetrieveFailure := errors.New("failed to retrieve")

	type fields struct {
		cacheKey string
		optFns   []func(o *FileCacheOptions)
	}

	type mockProviderWithContextRetrieveWithContext struct {
		times int
		res   credentials.Value
		err   error
	}

	type mockProviderWithContextExpiresAt struct {
		times int
		res   time.Time
	}

	type expected struct {
		res credentials.Value
		err error
	}

	tests := []struct {
		name                                       string
		fields                                     fields
		mockProviderWithContextRetrieveWithContext mockProviderWithContextRetrieveWithContext
		mockProviderWithContextExpiresAt           mockProviderWithContextExpiresAt
		expected                                   expected
	}{
		{
			name: "positive case: cached credentials",
			fields: fields{
				cacheKey: "cached",
				optFns:   []func(o *FileCacheOptions){},
			},
			mockProviderWithContextRetrieveWithContext: mockProviderWithContextRetrieveWithContext{
				times: 0,
				res: credentials.Value{
					AccessKeyID:     "NonCachedAccessKeyID",
					SecretAccessKey: "NonCachedSecretAccessKey",
					SessionToken:    "NonCachedSessionToken",
					ProviderName:    "TestProvider",
				},
				err: nil,
			},
			mockProviderWithContextExpiresAt: mockProviderWithContextExpiresAt{
				times: 0,
				res:   expiresIn15Minutes,
			},
			expected: expected{
				res: credentials.Value{
					AccessKeyID:     "CachedAccessKeyID",
					SecretAccessKey: "CachedSecretAccessKey",
					SessionToken:    "CachedSessionToken",
					ProviderName:    "FileCacheProvider",
				},
				err: nil,
			},
		},
		{
			name: "positive case: expired credentials",
			fields: fields{
				cacheKey: "expired",
				optFns:   []func(o *FileCacheOptions){},
			},
			mockProviderWithContextRetrieveWithContext: mockProviderWithContextRetrieveWithContext{
				times: 1,
				res: credentials.Value{
					AccessKeyID:     "NonCachedAccessKeyID",
					SecretAccessKey: "NonCachedSecretAccessKey",
					SessionToken:    "NonCachedSessionToken",
					ProviderName:    "TestProvider",
				},
				err: nil,
			},
			mockProviderWithContextExpiresAt: mockProviderWithContextExpiresAt{
				times: 1,
				res:   expiresIn15Minutes,
			},
			expected: expected{
				res: credentials.Value{
					AccessKeyID:     "NonCachedAccessKeyID",
					SecretAccessKey: "NonCachedSecretAccessKey",
					SessionToken:    "NonCachedSessionToken",
					ProviderName:    "FileCacheProvider",
				},
				err: nil,
			},
		},
		{
			name: "positive case: non-cached credentials",
			fields: fields{
				cacheKey: "non-cached",
				optFns:   []func(o *FileCacheOptions){},
			},
			mockProviderWithContextRetrieveWithContext: mockProviderWithContextRetrieveWithContext{
				times: 1,
				res: credentials.Value{
					AccessKeyID:     "NonCachedAccessKeyID",
					SecretAccessKey: "NonCachedSecretAccessKey",
					SessionToken:    "NonCachedSessionToken",
					ProviderName:    "TestProvider",
				},
				err: nil,
			},
			mockProviderWithContextExpiresAt: mockProviderWithContextExpiresAt{
				times: 1,
				res:   expiresIn15Minutes,
			},
			expected: expected{
				res: credentials.Value{
					AccessKeyID:     "NonCachedAccessKeyID",
					SecretAccessKey: "NonCachedSecretAccessKey",
					SessionToken:    "NonCachedSessionToken",
					ProviderName:    "FileCacheProvider",
				},
				err: nil,
			},
		},
		{
			name: "negative case: failed to retrieve (expired credentials)",
			fields: fields{
				cacheKey: "expired",
				optFns:   []func(o *FileCacheOptions){},
			},
			mockProviderWithContextRetrieveWithContext: mockProviderWithContextRetrieveWithContext{
				times: 1,
				res:   credentials.Value{ProviderName: "TestProvider"},
				err:   errRetrieveFailure,
			},
			mockProviderWithContextExpiresAt: mockProviderWithContextExpiresAt{
				times: 0,
				res:   time.Time{},
			},
			expected: expected{
				res: credentials.Value{ProviderName: "FileCacheProvider"},
				err: errRetrieveFailure,
			},
		},
		{
			name: "negative case: failed to retrieve (non-cached credentials)",
			fields: fields{
				cacheKey: "non-cached",
				optFns:   []func(o *FileCacheOptions){},
			},
			mockProviderWithContextRetrieveWithContext: mockProviderWithContextRetrieveWithContext{
				times: 1,
				res:   credentials.Value{ProviderName: "TestProvider"},
				err:   errRetrieveFailure,
			},
			mockProviderWithContextExpiresAt: mockProviderWithContextExpiresAt{
				times: 0,
				res:   time.Time{},
			},
			expected: expected{
				res: credentials.Value{ProviderName: "FileCacheProvider"},
				err: errRetrieveFailure,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			cachedDir := t.TempDir()
			StoreCredentials(filepath.Join(cachedDir, fmt.Sprintf("%s.json", "cached")), cachedCreds, expiresIn15Minutes)
			StoreCredentials(filepath.Join(cachedDir, fmt.Sprintf("%s.json", "expired")), cachedCreds, expired15MinutesAgo)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProviderWithContext := mock_credscache.NewMockexpireProviderWithContext(ctrl)
			mockProviderWithContext.
				EXPECT().
				RetrieveWithContext(gomock.Any()).
				Return(tt.mockProviderWithContextRetrieveWithContext.res, tt.mockProviderWithContextRetrieveWithContext.err).
				Times(tt.mockProviderWithContextRetrieveWithContext.times)
			mockProviderWithContext.
				EXPECT().
				ExpiresAt().
				Return(tt.mockProviderWithContextExpiresAt.res).
				Times(tt.mockProviderWithContextExpiresAt.times)

			tt.fields.optFns = append(tt.fields.optFns, func(o *FileCacheOptions) { o.FileCacheDir = cachedDir })

			provider := NewFileCacheProvider(mockProviderWithContext, tt.fields.cacheKey, tt.fields.optFns...)

			// Act
			actual, err := provider.Retrieve()

			// Assert
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

func TestFileCacheProvider_RetrieveWithNonExpireProvider(t *testing.T) {
	errRetrieveFailure := errors.New("failed to retrieve")

	type fields struct {
		cacheKey string
		optFns   []func(o *FileCacheOptions)
	}

	type mockProviderWithContextRetrieveWithContext struct {
		times int
		res   credentials.Value
		err   error
	}

	type expected struct {
		res credentials.Value
		err error
	}

	tests := []struct {
		name                                       string
		fields                                     fields
		mockProviderWithContextRetrieveWithContext mockProviderWithContextRetrieveWithContext
		expected                                   expected
	}{
		{
			name: "positive case: constant credentials",
			fields: fields{
				cacheKey: "constant",
				optFns:   []func(o *FileCacheOptions){},
			},
			mockProviderWithContextRetrieveWithContext: mockProviderWithContextRetrieveWithContext{
				times: 1,
				res: credentials.Value{
					AccessKeyID:     "ConstantAccessKeyID",
					SecretAccessKey: "ConstantSecretAccessKey",
					SessionToken:    "ConstantSessionToken",
					ProviderName:    "TestProvider",
				},
				err: nil,
			},
			expected: expected{
				res: credentials.Value{
					AccessKeyID:     "ConstantAccessKeyID",
					SecretAccessKey: "ConstantSecretAccessKey",
					SessionToken:    "ConstantSessionToken",
					ProviderName:    "FileCacheProvider",
				},
				err: nil,
			},
		},
		{
			name: "negative case: failed to retrieve",
			fields: fields{
				cacheKey: "constant",
				optFns:   []func(o *FileCacheOptions){},
			},
			mockProviderWithContextRetrieveWithContext: mockProviderWithContextRetrieveWithContext{
				times: 1,
				res:   credentials.Value{ProviderName: "TestProvider"},
				err:   errRetrieveFailure,
			},
			expected: expected{
				res: credentials.Value{ProviderName: "FileCacheProvider"},
				err: errRetrieveFailure,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProviderWithContext := mock_credentials.NewMockProviderWithContext(ctrl)
			mockProviderWithContext.
				EXPECT().
				RetrieveWithContext(gomock.Any()).
				Return(tt.mockProviderWithContextRetrieveWithContext.res, tt.mockProviderWithContextRetrieveWithContext.err).
				Times(tt.mockProviderWithContextRetrieveWithContext.times)

			provider := NewFileCacheProvider(mockProviderWithContext, tt.fields.cacheKey, tt.fields.optFns...)

			// Act
			actual, err := provider.Retrieve()

			// Assert
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

func TestFileCacheProvider_IsExpiredWithExpireProvider(t *testing.T) {
	expiresIn15Minutes := time.Now().UTC().Add(time.Duration(15) * time.Minute)
	expired15MinutesAgo := time.Now().UTC().Add(-time.Duration(15) * time.Minute)

	type fields struct {
		expires time.Time
		optFns  []func(o *FileCacheOptions)
	}

	type expected struct {
		res bool
	}

	tests := []struct {
		name     string
		fields   fields
		expected expected
	}{
		{
			name: "positive case: not expired",
			fields: fields{
				expires: expiresIn15Minutes,
				optFns:  []func(o *FileCacheOptions){},
			},
			expected: expected{
				res: false,
			},
		},
		{
			name: "positive case: expired",
			fields: fields{
				expires: expired15MinutesAgo,
				optFns:  []func(o *FileCacheOptions){},
			},
			expected: expected{
				res: true,
			},
		},
		{
			name: "positive case: expired by ExpiryWindow",
			fields: fields{
				expires: expiresIn15Minutes,
				optFns:  []func(o *FileCacheOptions){func(o *FileCacheOptions) { o.ExpiryWindow = time.Duration(15) * time.Minute }},
			},
			expected: expected{
				res: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProviderWithContext := mock_credscache.NewMockexpireProviderWithContext(ctrl)

			provider := NewFileCacheProvider(mockProviderWithContext, "key", tt.fields.optFns...)
			provider.SetExpiration(tt.fields.expires, provider.options.ExpiryWindow)

			// Act
			actual := provider.IsExpired()

			// Assert
			assert.Equal(t, tt.expected.res, actual)
		})
	}
}

func TestFileCacheProvider_IsExpiredWithNonExpireProvider(t *testing.T) {
	type expected struct {
		res bool
	}

	tests := []struct {
		name     string
		expected expected
	}{
		{
			name: "positive case: not expired",
			expected: expected{
				res: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProviderWithContext := mock_credentials.NewMockProviderWithContext(ctrl)

			provider := NewFileCacheProvider(mockProviderWithContext, "key")

			// Act
			actual := provider.IsExpired()

			// Assert
			assert.Equal(t, tt.expected.res, actual)
		})
	}
}
