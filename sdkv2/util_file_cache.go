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
	"github.com/Aton-Kish/aws-credscache-go/credscacheutil"
	"github.com/aws/aws-sdk-go-v2/aws"
)

type FileCache struct {
	credscacheutil.FileCache
}

var _ interface {
	credscacheutil.Loader
	credscacheutil.Storer
	FromAWSCredentials(creds *aws.Credentials)
	ToAWSCredentials(source string) *aws.Credentials
} = &FileCache{}

func NewFileCache() *FileCache {
	return new(FileCache)
}

func NewFileCacheFromAWSCredentials(creds *aws.Credentials) *FileCache {
	c := NewFileCache()
	c.FromAWSCredentials(creds)
	return c
}

func (c *FileCache) FromAWSCredentials(creds *aws.Credentials) {
	c.Credentials.AccessKeyID = creds.AccessKeyID
	c.Credentials.SecretAccessKey = creds.SecretAccessKey
	c.Credentials.SessionToken = creds.SessionToken
	c.Credentials.Expires = creds.Expires
}

func (c *FileCache) ToAWSCredentials(source string) *aws.Credentials {
	return &aws.Credentials{
		AccessKeyID:     c.Credentials.AccessKeyID,
		SecretAccessKey: c.Credentials.SecretAccessKey,
		SessionToken:    c.Credentials.SessionToken,
		Source:          source,
		CanExpire:       true,
		Expires:         c.Credentials.Expires,
	}
}
