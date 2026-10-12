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

func TestGetQuoteEmptyCartRegression(t *testing.T) {
	s := server{}
	res, err := s.GetQuote(context.Background(), &pb.GetQuoteRequest{})
	if err != nil {
		t.Fatalf("GetQuote: %v", err)
	}
	if res.CostUsd.GetUnits() != 0 || res.CostUsd.GetNanos() != 0 {
		t.Fatalf("empty cart should have $0 shipping, got %d.%09d",
			res.CostUsd.GetUnits(), res.CostUsd.GetNanos())
	}
}

func TestGetQuoteSumsItemQuantitiesRegression(t *testing.T) {
	s := server{}
	res, err := s.GetQuote(context.Background(), &pb.GetQuoteRequest{
		Items: []*pb.CartItem{
			{ProductId: "a", Quantity: 2},
			{ProductId: "b", Quantity: 3},
		},
	})
	if err != nil {
		t.Fatalf("GetQuote: %v", err)
	}
	if res.CostUsd.GetUnits() != 8 || res.CostUsd.GetNanos() != 990000000 {
		t.Fatalf("non-empty cart should be flat $8.99, got %d.%09d",
			res.CostUsd.GetUnits(), res.CostUsd.GetNanos())
	}
}
