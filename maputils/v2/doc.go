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

// Package maputils provides helpers for transforming Go maps into
// formats required by other subsystems. The primary consumer is
// structured logging (slog), which accepts key-value pairs as
// variadic []any arguments. Converting a map into that layout is
// mechanical but easy to get wrong (forgetting to interleave keys
// and values), so centralizing the conversion prevents subtle
// logging bugs across all commands.
package maputils
