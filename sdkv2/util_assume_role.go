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
	"fmt"

	"github.com/Aton-Kish/aws-credscache-go/credscacheutil"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
)

type AssumeRoleCacheKeyGenerator struct {
	credscacheutil.AssumeRoleCacheKeyGenerator
}

var _ interface {
	fmt.Stringer
	credscacheutil.CacheKeyer
	FromSTSCredsAssumeRoleOptions(options *stscreds.AssumeRoleOptions)
	ToSTSCredsAssumeRoleOptions() *stscreds.AssumeRoleOptions
} = &AssumeRoleCacheKeyGenerator{}

func NewAssumeRoleCacheKeyGenerator() *AssumeRoleCacheKeyGenerator {
	return new(AssumeRoleCacheKeyGenerator)
}

func NewAssumeRoleCacheKeyGeneratorFromSTSCredsAssumeRoleOptions(options *stscreds.AssumeRoleOptions) *AssumeRoleCacheKeyGenerator {
	g := NewAssumeRoleCacheKeyGenerator()
	g.FromSTSCredsAssumeRoleOptions(options)
	return g
}

func (g *AssumeRoleCacheKeyGenerator) FromSTSCredsAssumeRoleOptions(options *stscreds.AssumeRoleOptions) {
	g.RoleARN = options.RoleARN
	g.RoleSessionName = options.RoleSessionName
	g.ExternalID = options.ExternalID
	g.SerialNumber = options.SerialNumber
	g.Duration = options.Duration
}

func (g *AssumeRoleCacheKeyGenerator) ToSTSCredsAssumeRoleOptions() *stscreds.AssumeRoleOptions {
	return &stscreds.AssumeRoleOptions{
		RoleARN:         g.RoleARN,
		RoleSessionName: g.RoleSessionName,
		ExternalID:      g.ExternalID,
		SerialNumber:    g.SerialNumber,
		Duration:        g.Duration,
	}
}
