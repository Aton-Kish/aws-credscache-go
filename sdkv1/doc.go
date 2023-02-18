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

// Package credscache provides credentials caching utilities for the AWS SDK for
// Go v1.
//
// # Inject the file cache provider
//
// By default, the file cache provider outputs cache files to the current
// directory.
//
//	sess, err := session.NewSessionWithOptions(session.Options{
//		SharedConfigState:       session.SharedConfigEnable,
//		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	injected, err := credscache.InjectFileCacheProvider(sess.Config)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if !injected {
//		log.Print("unable to inject file cache provider")
//	}
//
// You can share cache with the AWS CLI by specifying `$HOME/.aws/cli/cache`
// (experimental feature).
//
//	sess, err := session.NewSessionWithOptions(session.Options{
//		SharedConfigState:       session.SharedConfigEnable,
//		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	injected, err := credscache.InjectFileCacheProvider(sess.Config, func(o *credscache.FileCacheOptions) {
//		home, _ := os.UserHomeDir()
//		o.FileCacheDir = filepath.Join(home, ".aws/cli/cache")
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if !injected {
//		log.Print("unable to inject file cache provider")
//	}
package credscache
