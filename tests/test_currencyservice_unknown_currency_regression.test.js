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

const { test } = require('node:test');
const assert = require('node:assert/strict');
const path = require('node:path');

process.env.ENABLE_TRACING = '0';

const grpc = require(path.join(__dirname, '../src/currencyservice/node_modules/@grpc/grpc-js'));
const { convert } = require(path.join(__dirname, '../src/currencyservice/server.js'));

function convertAsync (request) {
  return new Promise((resolve, reject) => {
    convert({ request }, (err, result) => {
      if (err) {
        resolve({ err });
        return;
      }
      resolve({ result });
    });
  });
}

test('unknown from currency is rejected', async () => {
  const { err, result } = await convertAsync({
    from: { currency_code: 'XXX', units: 10, nanos: 0 },
    to_code: 'USD'
  });
  assert.ok(err, 'expected gRPC error for unknown from currency');
  assert.equal(err.code, grpc.status.INVALID_ARGUMENT);
  assert.equal(result, undefined);
});

test('unknown to currency is rejected', async () => {
  const { err, result } = await convertAsync({
    from: { currency_code: 'USD', units: 10, nanos: 0 },
    to_code: 'XXX'
  });
  assert.ok(err, 'expected gRPC error for unknown to currency');
  assert.equal(err.code, grpc.status.INVALID_ARGUMENT);
  assert.equal(result, undefined);
});
