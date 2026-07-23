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

// Package aws provides shared AWS SDK configuration helpers. Multiple
// CLI tools and services in this organization need authenticated AWS
// clients with consistent retry behavior, region resolution, and
// optional debug logging. Centralizing that setup here avoids
// duplicating boilerplate and ensures uniform defaults (retry counts,
// environment variable fallback order, log verbosity) across all
// consumers.
package aws
