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

// errUnknownRegion is returned when no region can be determined. AWS
// API calls require a region; failing fast here produces a clear error
// message rather than a cryptic SDK failure downstream.
var errUnknownRegion = errors.New("please specify an AWS region")

type (
	// NoOpRateLimit satisfies the SDK's rate-limiter interface but never
	// throttles. This is useful when the caller already manages
	// concurrency externally (e.g., worker pools with bounded
	// parallelism) and SDK-level rate limiting would add unnecessary
	// latency.
	NoOpRateLimit struct{}

	// AWSConfigOptions bundles the knobs callers can tweak when
	// building an AWS config. Grouping them in a struct keeps
	// GetAWSConfig's signature stable as new options are added over
	// time.
	AWSConfigOptions struct {
		Region  string
		Profile string
		Retries int
		Verbose bool
	}
)

// AddTokens is a no-op so the rate limiter never blocks on token
// replenishment.
func (NoOpRateLimit) AddTokens(uint) error {
	return nil
}

// GetToken always succeeds immediately, ensuring requests are never
// delayed by SDK-internal rate limiting.
func (NoOpRateLimit) GetToken(context.Context, uint) (func() error, error) {
	return noOpToken, nil
}

// noOpToken is the release callback returned by GetToken. It does
// nothing because there is no token budget to return to.
func noOpToken() error {
	return nil
}

// GetAWSConfig builds an aws.Config with opinionated defaults suitable
// for CLI tools: automatic region discovery from environment variables,
// bounded retries, and optional request/response logging for
// troubleshooting. The variadic opts parameter lets callers override
// defaults without changing every call site when a new option is added.
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

		// Fall back to the AWS_PROFILE environment variable so users
		// who set it globally don't have to pass --profile on every
		// invocation.
		if profile == "" {
			profile = os.Getenv("AWS_PROFILE")
		}

		// Region resolution follows the same precedence as the AWS CLI:
		// explicit option > AWS_REGION > AWS_DEFAULT_REGION. Failing
		// here is intentional — an API call without a region produces
		// confusing SDK errors.
		if region == "" {
			region, ok = os.LookupEnv("AWS_REGION")
			if !ok {
				region, ok = os.LookupEnv("AWS_DEFAULT_REGION")
				if !ok {
					return emptyConfig, errUnknownRegion
				}
			}
		}

		// Default to 3 retries as a safety net against transient
		// network errors without being so aggressive that a broken
		// endpoint stalls the CLI for too long.
		if retries == 0 {
			retries = 3
		}
	}

	// Pull AWS credentials from the environment. The SDK's default
	// credential chain handles ~/.aws/credentials, IAM roles, SSO
	// tokens, and container credentials — we just configure the
	// behavioral knobs on top.
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

	// Attach a stderr logger when verbose mode is on so AWS SDK debug
	// output goes to the diagnostic stream, not stdout (which may be
	// piped to jq or another tool).
	if verbose {
		conf.Logger = logging.NewStandardLogger(os.Stderr)
	}

	return conf, nil
}
