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
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

type AssumeRoleCacheKeyGenerator struct {
	RoleARN         string
	RoleSessionName string
	ExternalID      *string
	SerialNumber    *string
	Duration        time.Duration
}

var _ interface {
	fmt.Stringer
	CacheKeyer
} = &AssumeRoleCacheKeyGenerator{}

func (g AssumeRoleCacheKeyGenerator) String() string {
	opts := []string{}

	opts = append(opts, fmt.Sprintf(`"RoleArn": "%s"`, g.RoleARN))

	if g.RoleSessionName != "" {
		opts = append(opts, fmt.Sprintf(`"RoleSessionName": "%s"`, g.RoleSessionName))
	}

	if g.ExternalID != nil {
		opts = append(opts, fmt.Sprintf(`"ExternalId": "%s"`, *g.ExternalID))
	}

	if g.SerialNumber != nil {
		opts = append(opts, fmt.Sprintf(`"SerialNumber": "%s"`, *g.SerialNumber))
	}

	if g.Duration != 0 {
		opts = append(opts, fmt.Sprintf(`"DurationSeconds": %d`, int(g.Duration.Seconds())))
	}

	sort.Slice(opts, func(i, j int) bool { return opts[i] < opts[j] })

	return fmt.Sprintf("{%s}", strings.Join(opts, ", "))
}

func (g *AssumeRoleCacheKeyGenerator) CacheKey() (string, error) {
	hash := sha1.New()
	if _, err := hash.Write([]byte(g.String())); err != nil {
		return "", err
	}

	key := strings.ToLower(hex.EncodeToString(hash.Sum(nil)))

	return key, nil
}
