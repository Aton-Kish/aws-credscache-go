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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	credscache "github.com/Aton-Kish/aws-credscache-go/sdkv1"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

func ExampleAssumeRoleCacheKey() {
	key, err := credscache.AssumeRoleCacheKey(&stscreds.AssumeRoleProvider{
		RoleARN: "role_arn",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(key)
	// Output:
	// de1969e7a880d858c9bef3ba110acf78869d4527
}

func ExampleAssumeRoleCacheKey_withRoleSessionNameAndMFASerial() {
	key, err := credscache.AssumeRoleCacheKey(&stscreds.AssumeRoleProvider{
		RoleARN:         "role_arn",
		RoleSessionName: "role_session_name",
		SerialNumber:    aws.String("mfa_serial"),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(key)
	// Output:
	// cda918cacd9e1d1c71d510d187e90c5817e04b97
}

func ExampleLoadCredentials() {
	path := "/home/gopher/.aws/cli/cache/de1969e7a880d858c9bef3ba110acf78869d4527.json"
	creds, expires, err := credscache.LoadCredentials(path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(creds)
	// &credentials.Value{
	// 	AccessKeyID:     "AccessKeyID",
	// 	SecretAccessKey: "SecretAccessKey",
	// 	SessionToken:    "SessionToken",
	// 	ProviderName:    "",
	// }
	fmt.Println(expires)
	// time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)
}

func ExampleStoreCredentials() {
	path := "/home/gopher/.aws/cli/cache/de1969e7a880d858c9bef3ba110acf78869d4527.json"
	creds := &credentials.Value{
		AccessKeyID:     "AccessKeyID",
		SecretAccessKey: "SecretAccessKey",
		SessionToken:    "SessionToken",
		ProviderName:    "TestProvider",
	}
	expires := time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)
	if err := credscache.StoreCredentials(path, creds, expires); err != nil {
		log.Fatal(err)
	}
}

func ExampleInjectFileCacheProvider() {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	})
	if err != nil {
		log.Fatal(err)
	}

	injected, err := credscache.InjectFileCacheProvider(sess.Config)
	if err != nil {
		log.Fatal(err)
	}

	if !injected {
		log.Print("unable to inject file cache provider")
	}
}

func ExampleInjectFileCacheProvider_specifiedFileCacheDir() {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	})
	if err != nil {
		log.Fatal(err)
	}

	injected, err := credscache.InjectFileCacheProvider(sess.Config, func(o *credscache.FileCacheOptions) {
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
