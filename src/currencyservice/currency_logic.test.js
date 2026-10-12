'use strict';

const assert = require('node:assert/strict');
const { describe, it } = require('node:test');

const {
  carry,
  convertAmount,
  getSupportedCurrencyCodes,
} = require('./currency_logic');

describe('currency_logic', () => {
  it('carry normalizes fractional units into nanos', () => {
    const amount = carry({ units: 1.25, nanos: 0 });
    assert.equal(amount.units, 1);
    assert.equal(amount.nanos, 250000000);
  });

  it('lists supported currency codes from the ECB table', () => {
    const codes = getSupportedCurrencyCodes();
    assert.ok(codes.includes('EUR'));
    assert.ok(codes.includes('USD'));
  });

  it('converts EUR to USD using stored rates', () => {
    const result = convertAmount(
      { currency_code: 'EUR', units: 10, nanos: 0 },
      'USD'
    );
    assert.equal(result.currency_code, 'USD');
    assert.equal(result.units, 11);
    assert.ok(Math.abs(result.nanos - 305000000) <= 1);
  });

  it('rejects unknown currency codes', () => {
    assert.throws(
      () => convertAmount({ currency_code: 'ZZZ', units: 1, nanos: 0 }, 'USD'),
      /unsupported currency conversion/
    );
  });
});
