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

var trackingIDPattern = regexp.MustCompile(`^[A-Z]{2}-\d+-\d+$`)

func trackingIDLength(salt string) int {
	middle := len(fmt.Sprintf("%d", len(salt))) + 3
	final := len(fmt.Sprintf("%d", len(salt)/2)) + 7
	return 2 + 1 + middle + 1 + final
}

func TestCreateTrackingIdFormat(t *testing.T) {
	for _, salt := range []string{"", "ab", "hipstershop", "Muffin Man, London, "} {
		id := CreateTrackingId(salt)
		wantLen := trackingIDLength(salt)
		if len(id) != wantLen {
			t.Fatalf("CreateTrackingId(%q) length %d, want %d (%q)", salt, len(id), wantLen, id)
		}
		if !trackingIDPattern.MatchString(id) {
			t.Fatalf("CreateTrackingId(%q) = %q, does not match expected pattern", salt, id)
		}
	}
}

func TestCreateTrackingIdEmbedsSaltLengths(t *testing.T) {
	salt := "abcd"
	id := CreateTrackingId(salt)
	parts := regexp.MustCompile(`-`).Split(id, 3)
	if len(parts) != 3 {
		t.Fatalf("unexpected segments in %q", id)
	}
	if parts[1][0] != '4' {
		t.Fatalf("middle segment %q should start with salt length 4", parts[1])
	}
	if parts[2][0] != '2' {
		t.Fatalf("final segment %q should start with half salt length 2", parts[2])
	}
}
