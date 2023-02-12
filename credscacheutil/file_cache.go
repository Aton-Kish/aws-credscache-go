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

package credscacheutil

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Aton-Kish/aws-credscache-go/internal/xfilepath"
)

type Loader interface {
	Load(path string) error
}

type Storer interface {
	Store(path string) error
}

type FileCache struct {
	Credentials CachedCredentials `json:"Credentials"`
}

type CachedCredentials struct {
	AccessKeyID     string    `json:"AccessKeyId"`
	SecretAccessKey string    `json:"SecretAccessKey"`
	SessionToken    string    `json:"SessionToken"`
	Expires         time.Time `json:"Expiration"`
}

var _ interface {
	Loader
	Storer
} = &FileCache{}

func (c *FileCache) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("failed to read cache file, %w", err)
		return err
	}

	cache := new(FileCache)
	if err := json.Unmarshal(data, cache); err != nil {
		err = fmt.Errorf("failed to decode cache json, %w", err)
		return err
	}

	*c = *cache

	return nil
}

func (c *FileCache) Store(path string) error {
	data, err := json.Marshal(c)
	if err != nil {
		err = fmt.Errorf("failed to encode cache json, %w", err)
		return err
	}

	dir := filepath.Dir(path)
	if !xfilepath.Exists(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			err = fmt.Errorf("failed to make directories, %w", err)
			return err
		}
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		err = fmt.Errorf("failed to write cache file, %w", err)
		return err
	}

	return nil
}
