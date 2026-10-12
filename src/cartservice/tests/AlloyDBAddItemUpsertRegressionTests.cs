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

using System;
using System.IO;
using Xunit;

namespace cartservice.tests
{
    /// <summary>
    /// Regression guard: AlloyDB AddItem must upsert a single row per product, not insert duplicates.
    /// </summary>
    public class AlloyDBAddItemUpsertRegressionTests
    {
        [Fact]
        public void AddItemAsync_UsesUpsertOnUserAndProductConflict()
        {
            var storePath = Path.GetFullPath(
                Path.Combine(
                    AppContext.BaseDirectory,
                    "..",
                    "..",
                    "..",
                    "..",
                    "src",
                    "cartstore",
                    "AlloyDBCartStore.cs"));

            var source = File.ReadAllText(storePath);

            Assert.Contains("ON CONFLICT (userId, productId)", source);
            Assert.Contains("DO UPDATE SET quantity = EXCLUDED.quantity", source);
        }
    }
}
