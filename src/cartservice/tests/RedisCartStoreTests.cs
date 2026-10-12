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

using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using cartservice.cartstore;
using Microsoft.Extensions.Caching.Distributed;
using Xunit;

namespace cartservice.tests
{
    public class RedisCartStoreTests
    {
        private readonly RedisCartStore _store = new RedisCartStore(new InMemoryDistributedCache());

        [Fact]
        public async Task GetCart_ReturnsEmptyCartWhenUserMissing()
        {
            var cart = await _store.GetCartAsync("new-user");
            Assert.NotNull(cart);
            Assert.Empty(cart.Items);
        }

        [Fact]
        public async Task AddItem_CreatesCartAndIncrementsQuantity()
        {
            const string userId = "redis-cart-user-1";
            await _store.AddItemAsync(userId, "sku-1", 2);
            await _store.AddItemAsync(userId, "sku-1", 3);

            var cart = await _store.GetCartAsync(userId);
            Assert.Single(cart.Items);
            Assert.Equal("sku-1", cart.Items[0].ProductId);
            Assert.Equal(5, cart.Items[0].Quantity);

            await _store.EmptyCartAsync(userId);
        }

        [Fact]
        public async Task AddItem_AppendsDistinctProducts()
        {
            const string userId = "redis-cart-user-2";
            await _store.AddItemAsync(userId, "sku-a", 1);
            await _store.AddItemAsync(userId, "sku-b", 4);

            var cart = await _store.GetCartAsync(userId);
            Assert.Equal(2, cart.Items.Count);

            await _store.EmptyCartAsync(userId);
            cart = await _store.GetCartAsync(userId);
            Assert.Empty(cart.Items);
        }

        [Fact]
        public void Ping_ReturnsTrue()
        {
            Assert.True(_store.Ping());
        }

        private sealed class InMemoryDistributedCache : IDistributedCache
        {
            private readonly Dictionary<string, byte[]> _data = new Dictionary<string, byte[]>();

            public byte[] Get(string key) => _data.TryGetValue(key, out var value) ? value : null;

            public Task<byte[]> GetAsync(string key, CancellationToken token = default) =>
                Task.FromResult(Get(key));

            public void Refresh(string key) { }

            public Task RefreshAsync(string key, CancellationToken token = default) => Task.CompletedTask;

            public void Remove(string key) => _data.Remove(key);

            public Task RemoveAsync(string key, CancellationToken token = default)
            {
                Remove(key);
                return Task.CompletedTask;
            }

            public void Set(string key, byte[] value, DistributedCacheEntryOptions options) =>
                _data[key] = value;

            public Task SetAsync(string key, byte[] value, DistributedCacheEntryOptions options, CancellationToken token = default)
            {
                Set(key, value, options);
                return Task.CompletedTask;
            }
        }
    }
}
