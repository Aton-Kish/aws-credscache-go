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
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	credscache "github.com/Aton-Kish/aws-credscache-go/sdkv2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/cobra"
)

var sdkv2Cmd = &cobra.Command{
	Use:   "sdkv2",
	Short: "An example CLI for AWS SDK v2",
}

var sdkv2CacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Call AWS API with the cached credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		cfg, err := sdkv2LoadConfig(ctx, profile)
		if err != nil {
			return err
		}

		if err := sdkv2InjectFileCacheProvider(ctx, &cfg); err != nil {
			return err
		}

		if err := sdkv2PrintCallerIdentity(ctx, sts.NewFromConfig(cfg)); err != nil {
			return err
		}

		return nil
	},
}

var sdkv2NocacheCmd = &cobra.Command{
	Use:   "nocache",
	Short: "Call AWS API without the cached credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		cfg, err := sdkv2LoadConfig(ctx, profile)
		if err != nil {
			return err
		}

		if err := sdkv2PrintCallerIdentity(ctx, sts.NewFromConfig(cfg)); err != nil {
			return err
		}

		return nil
	},
}

func sdkv2LoadConfig(ctx context.Context, profile string) (aws.Config, error) {
	optFns := []func(o *config.LoadOptions) error{
		config.WithAssumeRoleCredentialOptions(func(o *stscreds.AssumeRoleOptions) {
			o.TokenProvider = stscreds.StdinTokenProvider
		}),
	}

	if profile != "" {
		optFns = append(optFns, config.WithSharedConfigProfile(profile))
	}

	return config.LoadDefaultConfig(ctx, optFns...)
}

func sdkv2InjectFileCacheProvider(ctx context.Context, cfg *aws.Config) error {
	optFns := []func(o *credscache.FileCacheOptions){
		func(o *credscache.FileCacheOptions) {
			home, _ := os.UserHomeDir()
			o.FileCacheDir = filepath.Join(home, ".aws/cli/cache")
		},
	}
	_, err := credscache.InjectFileCacheProvider(cfg, optFns...)

	return err
}

func sdkv2PrintCallerIdentity(ctx context.Context, client *sts.Client) error {
	output, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(output, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}

func init() {
	sdkv2Cmd.AddCommand(sdkv2CacheCmd, sdkv2NocacheCmd)
	rootCmd.AddCommand(sdkv2Cmd)
}
