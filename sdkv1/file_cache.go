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
	"time"

	"github.com/Aton-Kish/aws-credscache-go/credscacheutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

func LoadCredentials(path string) (*credentials.Value, time.Time, error) {
	cache := new(credscacheutil.FileCache)
	if err := cache.Load(path); err != nil {
		return nil, time.Time{}, err
	}

	creds := &credentials.Value{
		AccessKeyID:     cache.Credentials.AccessKeyID,
		SecretAccessKey: cache.Credentials.SecretAccessKey,
		SessionToken:    cache.Credentials.SessionToken,
	}

	return creds, cache.Credentials.Expires, nil
}

func StoreCredentials(path string, creds *credentials.Value, expires time.Time) error {
	cache := &credscacheutil.FileCache{
		Credentials: credscacheutil.CachedCredentials{
			AccessKeyID:     creds.AccessKeyID,
			SecretAccessKey: creds.SecretAccessKey,
			SessionToken:    creds.SessionToken,
			Expires:         expires,
		},
	}

	if err := cache.Store(path); err != nil {
		return err
	}

	return nil
}
