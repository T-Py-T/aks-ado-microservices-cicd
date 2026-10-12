// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"math/rand"
	"testing"
)

// getRandomLetterCode must be able to emit every uppercase ASCII letter A–Z.
func TestGetRandomLetterCodeIncludesZ(t *testing.T) {
	seen := make(map[rune]bool)
	for seed := int64(0); seed < 10000; seed++ {
		rand.Seed(seed)
		seen[rune(getRandomLetterCode())] = true
	}
	if !seen['Z'] {
		t.Fatal("getRandomLetterCode never produced 'Z'; expected full A–Z range")
	}
}
