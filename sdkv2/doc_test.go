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

package credscache_test

import (
	"context"
	"log"
	"os"
	"path/filepath"

	credscache "github.com/Aton-Kish/aws-credscache-go/sdkv2"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
)

func ExampleInjectFileCacheProvider() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithAssumeRoleCredentialOptions(func(options *stscreds.AssumeRoleOptions) {
		options.TokenProvider = stscreds.StdinTokenProvider
	}))
	if err != nil {
		log.Fatal(err)
	}

	injected, err := credscache.InjectFileCacheProvider(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if !injected {
		log.Print("unable to inject file cache provider")
	}
}

func ExampleInjectFileCacheProvider_specifiedFileCacheDir() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithAssumeRoleCredentialOptions(func(options *stscreds.AssumeRoleOptions) {
		options.TokenProvider = stscreds.StdinTokenProvider
	}))
	if err != nil {
		log.Fatal(err)
	}

	injected, err := credscache.InjectFileCacheProvider(&cfg, func(o *credscache.FileCacheOptions) {
		home, _ := os.UserHomeDir()
		o.FileCacheDir = filepath.Join(home, ".aws/cli/cache")
	})
	if err != nil {
		log.Fatal(err)
	}

	if !injected {
		log.Print("unable to inject file cache provider")
	}
}
