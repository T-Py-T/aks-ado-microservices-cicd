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

using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using cartservice.cartstore;
using Microsoft.Extensions.Caching.Distributed;
using Xunit;

namespace cartservice.tests
{
    public class EmptyCartPreservesUserIdRegressionTests
    {
        [Fact]
        public async Task EmptyCartAsync_PreservesUserIdOnStoredCart()
        {
            const string userId = "empty-cart-user-regression";
            var store = new RedisCartStore(new InMemoryDistributedCache());
            await store.AddItemAsync(userId, "sku-1", 1);

            await store.EmptyCartAsync(userId);

            var cart = await store.GetCartAsync(userId);
            Assert.Equal(userId, cart.UserId);
            Assert.Empty(cart.Items);
        }

        private sealed class InMemoryDistributedCache : IDistributedCache
        {
            private readonly Dictionary<string, byte[]> _data = new Dictionary<string, byte[]>();

            public byte[] Get(string key) => _data.TryGetValue(key, out var value) ? value : null;

            public Task<byte[]> GetAsync(string key, CancellationToken token = default) =>
                Task.FromResult(Get(key));

            public void Set(string key, byte[] value) => _data[key] = value;

            public void Set(string key, byte[] value, DistributedCacheEntryOptions options) =>
                Set(key, value);

            public Task SetAsync(string key, byte[] value, DistributedCacheEntryOptions options,
                CancellationToken token = default)
            {
                Set(key, value, options);
                return Task.CompletedTask;
            }

            public void Remove(string key) => _data.Remove(key);

            public Task RemoveAsync(string key, CancellationToken token = default)
            {
                Remove(key);
                return Task.CompletedTask;
            }

            public void Refresh(string key) { }

            public Task RefreshAsync(string key, CancellationToken token = default) => Task.CompletedTask;
        }
    }
}
