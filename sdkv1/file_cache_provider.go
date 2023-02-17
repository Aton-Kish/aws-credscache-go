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

//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=../internal/mock/github.com/Aton-Kish/aws-credscache-go/sdkv1/$GOFILE

package credscache

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/Aton-Kish/aws-credscache-go/internal/xfilepath"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

const (
	FileCacheProviderName = "FileCacheProvider"
)

var (
	defaultFileCacheDir = ""
	defaultExpiryWindow = time.Duration(1) * time.Minute
)

type expireProviderWithContext interface {
	credentials.ProviderWithContext
	credentials.Expirer
}

type FileCacheProvider struct {
	credentials.Expiry
	provider credentials.ProviderWithContext
	cacheKey string
	options  FileCacheOptions
}

type FileCacheOptions struct {
	FileCacheDir string
	ExpiryWindow time.Duration
}

var _ interface {
	expireProviderWithContext
} = &FileCacheProvider{}

func NewFileCacheProvider(provider credentials.ProviderWithContext, cacheKey string, optFns ...func(o *FileCacheOptions)) *FileCacheProvider {
	o := FileCacheOptions{
		FileCacheDir: defaultFileCacheDir,
		ExpiryWindow: defaultExpiryWindow,
	}

	for _, fn := range optFns {
		fn(&o)
	}

	return &FileCacheProvider{
		provider: provider,
		cacheKey: cacheKey,
		options:  o,
	}
}

func (p *FileCacheProvider) Retrieve() (credentials.Value, error) {
	return p.RetrieveWithContext(context.Background())
}

func (p *FileCacheProvider) RetrieveWithContext(ctx context.Context) (credentials.Value, error) {
	path := filepath.Join(p.options.FileCacheDir, fmt.Sprintf("%s.json", p.cacheKey))

	if xfilepath.Exists(path) {
		creds, expires, err := LoadCredentials(path)
		if err != nil {
			err = &FileCacheProviderError{Err: err}
			return credentials.Value{ProviderName: FileCacheProviderName}, err
		}
		creds.ProviderName = FileCacheProviderName

		p.SetExpiration(expires, p.options.ExpiryWindow)

		if !p.IsExpired() {
			return *creds, nil
		}
	}

	creds, err := p.provider.RetrieveWithContext(ctx)
	if err != nil {
		err = &FileCacheProviderError{Err: err}
		return credentials.Value{ProviderName: FileCacheProviderName}, err
	}
	creds.ProviderName = FileCacheProviderName

	if expirer, ok := p.provider.(credentials.Expirer); ok {
		expires := expirer.ExpiresAt()

		p.SetExpiration(expires, p.options.ExpiryWindow)

		if err := StoreCredentials(path, &creds, expires); err != nil {
			err = &FileCacheProviderError{Err: err}
			return credentials.Value{ProviderName: FileCacheProviderName}, err
		}
	}

	return creds, nil
}

func (p *FileCacheProvider) IsExpired() bool {
	if _, ok := p.provider.(credentials.Expirer); ok {
		return p.Expiry.IsExpired()
	}

	return false
}
