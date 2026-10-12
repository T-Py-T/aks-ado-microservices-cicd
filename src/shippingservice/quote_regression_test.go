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
	"context"
	"testing"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/shippingservice/genproto"
)

func TestCreateQuoteFromCountZeroItems(t *testing.T) {
	quote := CreateQuoteFromCount(0)
	if quote.Dollars != 0 || quote.Cents != 0 {
		t.Fatalf("expected zero quote for empty cart, got %s", quote)
	}
}

func TestGetQuoteUsesItemQuantities(t *testing.T) {
	s := server{}
	res, err := s.GetQuote(context.Background(), &pb.GetQuoteRequest{
		Items: []*pb.CartItem{{ProductId: "x", Quantity: 2}},
	})
	if err != nil {
		t.Fatalf("GetQuote failed: %v", err)
	}
	if res.CostUsd.GetUnits() != 8 || res.CostUsd.GetNanos() != 990000000 {
		t.Fatalf("expected $8.99 for non-empty shipment, got %d.%09d",
			res.CostUsd.GetUnits(), res.CostUsd.GetNanos())
	}

	empty, err := s.GetQuote(context.Background(), &pb.GetQuoteRequest{})
	if err != nil {
		t.Fatalf("GetQuote failed: %v", err)
	}
	if empty.CostUsd.GetUnits() != 0 || empty.CostUsd.GetNanos() != 0 {
		t.Fatalf("expected zero quote for no items, got %d.%09d",
			empty.CostUsd.GetUnits(), empty.CostUsd.GetNanos())
	}
}
