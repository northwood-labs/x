// Copyright 2026, Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aws

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/smithy-go/logging"
)

var errUnknownRegion = errors.New("please specify an AWS region")

type (
	// NoOpRateLimit to prevent limiting of queries to AWS.
	NoOpRateLimit struct{}

	// AWSConfigOptions to manage AWS configuration options.
	AWSConfigOptions struct {
		Region  string
		Profile string
		Retries int
		Verbose bool
	}
)

// AddTokens to return nil for NoOpRateLimit.
func (NoOpRateLimit) AddTokens(uint) error {
	return nil
}

// GetToken will return nil so that there will be no rate limiting.
func (NoOpRateLimit) GetToken(context.Context, uint) (func() error, error) {
	return noOpToken, nil
}

func noOpToken() error {
	return nil
}

// GetAWSConfig returns a standard AWS config object pre-configured for use with
// regions, retries, and verbosity.
//
// If region is empty, we will attempt to read AWS_REGION then
// AWS_DEFAULT_REGION.
func GetAWSConfig(ctx context.Context, opts ...AWSConfigOptions) (aws.Config, error) {
	var (
		emptyConfig = aws.Config{}

		ok      bool
		region  string
		profile string
		retries int
		verbose bool
	)

	if len(opts) > 0 {
		opt := opts[0]

		region = opt.Region
		profile = opt.Profile
		retries = opt.Retries
		verbose = opt.Verbose

		if profile == "" {
			profile = os.Getenv("AWS_PROFILE")
		}

		if region == "" {
			region, ok = os.LookupEnv("AWS_REGION")
			if !ok {
				region, ok = os.LookupEnv("AWS_DEFAULT_REGION")
				if !ok {
					return emptyConfig, errUnknownRegion
				}
			}
		}

		if retries == 0 {
			retries = 3
		}
	}

	// Pull AWS credentials from the environment.
	conf, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
		config.WithRetryer(func() aws.Retryer {
			// https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/configure-retries-timeouts.html
			retryLogic := retry.NewStandard()
			retry.AddWithMaxAttempts(retryLogic, retries)

			return retryLogic
		}),
		config.WithSharedConfigProfile(profile),
		func(verbose bool) config.LoadOptionsFunc {
			// https://docs.aws.amazon.com/sdk-for-go/v2/developer-guide/configure-logging.html
			if !verbose {
				return config.WithClientLogMode(0)
			}

			// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws#ClientLogMode
			return config.WithClientLogMode(
				aws.LogRetries |
					aws.LogRequestWithBody |
					aws.LogResponseWithBody |
					aws.LogDeprecatedUsage |
					aws.LogRequestEventMessage |
					aws.LogResponseEventMessage,
			)
		}(verbose),
	)
	if err != nil {
		return emptyConfig, fmt.Errorf("AWS configuration error: %w", err)
	}

	if verbose {
		conf.Logger = logging.NewStandardLogger(os.Stderr)
	}

	return conf, nil
}
