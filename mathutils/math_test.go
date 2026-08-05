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

package mathutils

import (
	"testing"

	"github.com/go-openapi/testify/assert"
)

func TestClamp(t *testing.T) {
	assert.Equal(t, 5, Clamp(5, 0, 10))
	assert.Equal(t, 0, Clamp(-5, 0, 10))
	assert.Equal(t, 10, Clamp(15, 0, 10))
	assert.Equal(t, 2.5, Clamp(2.5, 0.0, 5.0))
}

func TestInRange(t *testing.T) {
	assert.True(t, InRange(5, 0, 10))
	assert.True(t, InRange(0, 0, 10))
	assert.True(t, InRange(10, 0, 10))
	assert.False(t, InRange(-1, 0, 10))
	assert.False(t, InRange(11, 0, 10))
}
