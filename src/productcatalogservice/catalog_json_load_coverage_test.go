// Copyright 2026 Google LLC
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
	"os"
	"path/filepath"
	"testing"

	pb "github.com/GoogleCloudPlatform/microservices-demo/src/productcatalogservice/genproto"
)

func TestReadCatalogFileParsesBundledProducts(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(wd) != "productcatalogservice" {
		if err := os.Chdir(filepath.Join(wd, "src", "productcatalogservice")); err != nil {
			t.Fatalf("chdir to productcatalogservice: %v", err)
		}
		defer func() { _ = os.Chdir(wd) }()
	}

	catalog := &pb.ListProductsResponse{}
	if err := readCatalogFile(catalog); err != nil {
		t.Fatalf("readCatalogFile: %v", err)
	}
	if len(catalog.Products) == 0 {
		t.Fatal("expected products from products.json")
	}
	if catalog.Products[0].GetId() == "" {
		t.Error("first product should have an id")
	}
}
