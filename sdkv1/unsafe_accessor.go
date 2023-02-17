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

	"github.com/aws/aws-sdk-go/aws/credentials"
)

type CredentialsUnsafeAccessor struct {
	ptr *credentials.Credentials
}

func NewCredentialsUnsafeAccessor(ptr *credentials.Credentials) (*CredentialsUnsafeAccessor, error) {
	if ptr == nil {
		return nil, ErrNilPointer
	}

	a := &CredentialsUnsafeAccessor{
		ptr: ptr,
	}

	return a, nil
}

func (a *CredentialsUnsafeAccessor) provider() *credentials.Provider {
	v := reflect.ValueOf(a.ptr).Elem()
	f := v.FieldByName("provider")
	ptr := (*credentials.Provider)(unsafe.Pointer(f.UnsafeAddr()))
	return ptr
}

func (a *CredentialsUnsafeAccessor) Provider() credentials.Provider {
	ptr := a.provider()
	return *ptr
}

func (a *CredentialsUnsafeAccessor) SetProvider(provider credentials.Provider) {
	ptr := a.provider()
	*ptr = provider
}
