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
	"fmt"
	"regexp"
	"testing"
)

// Two letters, salt length and random digits, half-length and random digits.
var trackingIDPattern = regexp.MustCompile(`^[A-Z]{2}-\d+\d{3}-\d+\d{7}$`)

func TestCreateTrackingIdShape(t *testing.T) {
	for _, salt := range []string{"Muffin Man", "1600 Amphitheatre Parkway", ""} {
		id := CreateTrackingId(salt)
		if !trackingIDPattern.MatchString(id) {
			t.Fatalf("CreateTrackingId(%q) = %q; does not match expected pattern", salt, id)
		}
	}
}

func TestCreateTrackingIdEmbedsSaltLengths(t *testing.T) {
	salt := "abcde"
	id := CreateTrackingId(salt)
	// Format: LL-<len(salt)><3digits>-<len(salt)/2><7digits>
	wantPrefix := fmt.Sprintf("%c%c-%d", id[0], id[1], len(salt))
	if id[:len(wantPrefix)] != wantPrefix {
		t.Fatalf("tracking id %q should start with %q", id, wantPrefix)
	}
}
