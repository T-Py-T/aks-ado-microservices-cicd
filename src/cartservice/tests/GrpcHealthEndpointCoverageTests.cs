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

using System.Threading.Tasks;
using Grpc.Health.V1;
using Grpc.Net.Client;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.TestHost;
using Microsoft.Extensions.Hosting;
using Xunit;
using static Grpc.Health.V1.Health;

namespace cartservice.tests
{
    public class GrpcHealthEndpointCoverageTests
    {
        private readonly IHostBuilder _host;

        public GrpcHealthEndpointCoverageTests()
        {
            _host = new HostBuilder().ConfigureWebHost(webBuilder =>
            {
                webBuilder
                    .UseStartup<Startup>()
                    .UseTestServer();
            });
        }

        [Fact]
        public async Task HealthCheck_ReturnsServingForDefaultCartStore()
        {
            using var server = await _host.StartAsync();
            var httpClient = server.GetTestClient();

            var channel = GrpcChannel.ForAddress(httpClient.BaseAddress, new GrpcChannelOptions
            {
                HttpClient = httpClient
            });

            var healthClient = new HealthClient(channel);
            var response = await healthClient.CheckAsync(new HealthCheckRequest());

            Assert.Equal(HealthCheckResponse.Types.ServingStatus.Serving, response.Status);
        }
    }
}
