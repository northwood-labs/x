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

package grammar

// Pluralize selects between the singular and plural form of a word
// based on the count. CLI status messages frequently display item
// counts, and returning the grammatically correct form here keeps
// that concern out of the command logic. Only U.S. English is
// supported because all current consumers produce English output.
func Pluralize(amount int, singular, plural string) string {
	if amount == 1 {
		return singular
	}

	return plural
}
