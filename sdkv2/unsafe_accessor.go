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
	"reflect"
	"unsafe"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
)

type CredentialsCacheUnsafeAccessor struct {
	ptr *aws.CredentialsCache
}

func NewCredentialsCacheUnsafeAccessor(ptr *aws.CredentialsCache) (*CredentialsCacheUnsafeAccessor, error) {
	if ptr == nil {
		return nil, ErrNilPointer
	}

	a := &CredentialsCacheUnsafeAccessor{
		ptr: ptr,
	}

	return a, nil
}

func (a *CredentialsCacheUnsafeAccessor) provider() *aws.CredentialsProvider {
	v := reflect.ValueOf(a.ptr).Elem()
	f := v.FieldByName("provider")
	ptr := (*aws.CredentialsProvider)(unsafe.Pointer(f.UnsafeAddr()))
	return ptr
}

func (a *CredentialsCacheUnsafeAccessor) Provider() aws.CredentialsProvider {
	ptr := a.provider()
	return *ptr
}

func (a *CredentialsCacheUnsafeAccessor) SetProvider(provider aws.CredentialsProvider) {
	ptr := a.provider()
	*ptr = provider
}

type AssumeRoleProviderUnsafeAccessor struct {
	ptr *stscreds.AssumeRoleProvider
}

func NewAssumeRoleProviderUnsafeAccessor(ptr *stscreds.AssumeRoleProvider) (*AssumeRoleProviderUnsafeAccessor, error) {
	if ptr == nil {
		return nil, ErrNilPointer
	}

	a := &AssumeRoleProviderUnsafeAccessor{
		ptr: ptr,
	}

	return a, nil
}

func (a *AssumeRoleProviderUnsafeAccessor) options() *stscreds.AssumeRoleOptions {
	v := reflect.ValueOf(a.ptr).Elem()
	f := v.FieldByName("options")
	ptr := (*stscreds.AssumeRoleOptions)(unsafe.Pointer(f.UnsafeAddr()))
	return ptr
}

func (a *AssumeRoleProviderUnsafeAccessor) Options() stscreds.AssumeRoleOptions {
	ptr := a.options()
	return *ptr
}
