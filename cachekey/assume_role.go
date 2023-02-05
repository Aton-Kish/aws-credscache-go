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

package cachekey

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"
	"time"
)

type AssumeRoleOptions struct {
	RoleARN         string
	RoleSessionName string
	ExternalID      *string
	SerialNumber    *string
	Duration        time.Duration
}

var _ interface {
	fmt.Stringer
	cachekey
} = &AssumeRoleOptions{}

func (o AssumeRoleOptions) String() string {
	opts := []string{}

	opts = append(opts, fmt.Sprintf(`"RoleArn": "%s"`, o.RoleARN))

	if o.RoleSessionName != "" {
		opts = append(opts, fmt.Sprintf(`"RoleSessionName": "%s"`, o.RoleSessionName))
	}

	if o.ExternalID != nil {
		opts = append(opts, fmt.Sprintf(`"ExternalId": "%s"`, *o.ExternalID))
	}

	if o.SerialNumber != nil {
		opts = append(opts, fmt.Sprintf(`"SerialNumber": "%s"`, *o.SerialNumber))
	}

	if o.Duration != 0 {
		opts = append(opts, fmt.Sprintf(`"DurationSeconds": %d`, int(o.Duration.Seconds())))
	}

	sort.Slice(opts, func(i, j int) bool { return opts[i] < opts[j] })

	return fmt.Sprintf("{%s}", strings.Join(opts, ", "))
}

func (o *AssumeRoleOptions) CacheKey() string {
	h := sha1.New()
	h.Write([]byte(o.String()))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}
