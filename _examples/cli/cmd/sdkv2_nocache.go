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

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/cobra"
)

var sdkv2NocacheCmd = &cobra.Command{
	Use:   "nocache",
	Short: "Call AWS API without the cached credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		loadOptFns := []func(o *config.LoadOptions) error{
			config.WithAssumeRoleCredentialOptions(func(o *stscreds.AssumeRoleOptions) {
				o.TokenProvider = stscreds.StdinTokenProvider
			}),
		}

		if profile != "" {
			loadOptFns = append(loadOptFns, config.WithSharedConfigProfile(profile))
		}

		cfg, err := config.LoadDefaultConfig(ctx, loadOptFns...)
		if err != nil {
			return err
		}

		stsClient := sts.NewFromConfig(cfg)
		output, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(data))

		return nil
	},
}

func init() {
	sdkv2Cmd.AddCommand(sdkv2NocacheCmd)
}
