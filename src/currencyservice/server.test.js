'use strict';

const assert = require('node:assert/strict');
const { describe, it } = require('node:test');

const { convert, getSupportedCurrencies, check } = require('./server');

function invoke (handler, request) {
  return new Promise((resolve, reject) => {
    handler({ request }, (err, response) => {
      if (err) {
        reject(err);
        return;
      }
      resolve(response);
    });
  });
}

describe('currencyservice gRPC handlers', () => {
  it('getSupportedCurrencies returns ECB currency codes', async () => {
    const response = await invoke(getSupportedCurrencies, {});
    assert.ok(response.currency_codes.includes('EUR'));
  });

  it('convert returns converted money', async () => {
    const response = await invoke(convert, {
      from: { currency_code: 'EUR', units: 1, nanos: 0 },
      to_code: 'USD',
    });
    assert.equal(response.currency_code, 'USD');
    assert.equal(response.units, 1);
    assert.equal(response.nanos, 130500000);
  });

  it('check reports SERVING', async () => {
    const response = await invoke(check, {});
    assert.equal(response.status, 'SERVING');
  });
});
