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
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/Aton-Kish/aws-credscache-go/internal/xfilepath"
	"github.com/aws/aws-sdk-go-v2/aws"
)

const FileCacheProviderName = "FileCacheProvider"

var (
	defaultFileCacheDir = ""
	defaultExpiryWindow = time.Duration(1) * time.Minute
)

type FileCacheProvider struct {
	provider aws.CredentialsProvider
	cacheKey string
	options  FileCacheOptions
}

type FileCacheOptions struct {
	FileCacheDir string
	ExpiryWindow time.Duration
}

var _ interface {
	aws.CredentialsProvider
} = &FileCacheProvider{}

func NewFileCacheProvider(provider aws.CredentialsProvider, cacheKey string, optFns ...func(o *FileCacheOptions)) *FileCacheProvider {
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

func (p *FileCacheProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	path := filepath.Join(p.options.FileCacheDir, fmt.Sprintf("%s.json", p.cacheKey))

	if xfilepath.Exists(path) {
		creds, err := LoadCredentials(path)
		if err != nil {
			return aws.Credentials{Source: FileCacheProviderName}, err
		}
		creds.Source = FileCacheProviderName

		if creds.Expires.After(time.Now().Add(p.options.ExpiryWindow)) {
			return *creds, nil
		}
	}

	creds, err := p.provider.Retrieve(ctx)
	if err != nil {
		return aws.Credentials{Source: FileCacheProviderName}, err
	}
	creds.Source = FileCacheProviderName

	if creds.CanExpire {
		if err := StoreCredentials(path, &creds); err != nil {
			return aws.Credentials{Source: FileCacheProviderName}, err
		}
	}

	return creds, nil
}
