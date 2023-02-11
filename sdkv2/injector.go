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
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
)

func InjectFileCacheProvider(cfg *aws.Config, optFns ...func(o *FileCacheOptions)) (bool, error) {
	credsCache, ok := cfg.Credentials.(*aws.CredentialsCache)
	if !ok {
		return false, nil
	}

	credsCacheAccessor, err := NewCredentialsCacheUnsafeAccessor(credsCache)
	if err != nil {
		return false, err
	}

	provider := credsCacheAccessor.Provider()
	assumeRoleProvider, ok := provider.(*stscreds.AssumeRoleProvider)
	if !ok {
		return false, nil
	}

	assumeRoleProviderAccessor, err := NewAssumeRoleProviderUnsafeAccessor(assumeRoleProvider)
	if err != nil {
		return false, err
	}

	options := assumeRoleProviderAccessor.Options()
	key, err := AssumeRoleCacheKey(&options)
	if err != nil {
		return false, err
	}

	fileCacheProvider := NewFileCacheProvider(assumeRoleProvider, key, optFns...)
	credsCacheAccessor.SetProvider(fileCacheProvider)

	return true, nil
}
