// Copyright 2023 Google LLC
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

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/productcatalogservice/genproto"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func catalogWithProducts() *productCatalog {
	return &productCatalog{
		catalog: pb.ListProductsResponse{
			Products: []*pb.Product{
				{Id: "p1", Name: "Wool Scarf", Description: "Hand-knit winter accessory"},
				{Id: "p2", Name: "Desk Lamp", Description: "Adjustable LED light"},
				{Id: "p3", Name: "Travel Mug", Description: "Insulated stainless steel"},
			},
		},
	}
}

func TestProductCatalogHealthCheck(t *testing.T) {
	catalog := catalogWithProducts()
	resp, err := catalog.Check(context.Background(), &healthpb.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if resp.Status != healthpb.HealthCheckResponse_SERVING {
		t.Fatalf("status = %v, want SERVING", resp.Status)
	}
}

func TestProductCatalogWatchUnimplemented(t *testing.T) {
	catalog := catalogWithProducts()
	err := catalog.Watch(&healthpb.HealthCheckRequest{}, nil)
	if status.Code(err) != codes.Unimplemented {
		t.Fatalf("Watch error code = %v, want Unimplemented", status.Code(err))
	}
}

func TestSearchProductsMatchesDescriptionCaseInsensitive(t *testing.T) {
	catalog := catalogWithProducts()
	resp, err := catalog.SearchProducts(context.Background(),
		&pb.SearchProductsRequest{Query: "LED"},
	)
	if err != nil {
		t.Fatalf("SearchProducts failed: %v", err)
	}
	if len(resp.Results) != 1 || resp.Results[0].GetId() != "p2" {
		t.Fatalf("unexpected results: %#v", resp.Results)
	}

	resp, err = catalog.SearchProducts(context.Background(),
		&pb.SearchProductsRequest{Query: "winter"},
	)
	if err != nil {
		t.Fatalf("SearchProducts failed: %v", err)
	}
	if len(resp.Results) != 1 || resp.Results[0].GetId() != "p1" {
		t.Fatalf("description search returned %#v", resp.Results)
	}
}

func TestSearchProductsNoMatches(t *testing.T) {
	catalog := catalogWithProducts()
	resp, err := catalog.SearchProducts(context.Background(),
		&pb.SearchProductsRequest{Query: "nonexistent-sku"},
	)
	if err != nil {
		t.Fatalf("SearchProducts failed: %v", err)
	}
	if len(resp.Results) != 0 {
		t.Fatalf("expected empty results, got %#v", resp.Results)
	}
}
