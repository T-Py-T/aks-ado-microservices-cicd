'use strict';

const assert = require('node:assert/strict');
const { describe, it } = require('node:test');

const charge = require('./charge');

const VALID_VISA = {
  credit_card_number: '4242424242424242',
  credit_card_cvv: 123,
  credit_card_expiration_year: 2099,
  credit_card_expiration_month: 12,
};

function chargeRequest (overrides = {}) {
  return {
    amount: { currency_code: 'USD', units: 99, nanos: 990000000 },
    credit_card: { ...VALID_VISA, ...overrides },
  };
}

describe('charge', () => {
  it('returns a transaction id for valid Visa cards', () => {
    const response = charge(chargeRequest());
    assert.match(response.transaction_id, /^[0-9a-f-]{36}$/);
  });

  it('accepts valid Mastercard numbers', () => {
    const response = charge(chargeRequest({
      credit_card_number: '5555555555554444',
    }));
    assert.ok(response.transaction_id);
  });

  it('rejects invalid card numbers', () => {
    assert.throws(
      () => charge(chargeRequest({ credit_card_number: '1234567890123456' })),
      /invalid/i
    );
  });

  it('rejects unsupported card brands', () => {
    assert.throws(
      () => charge(chargeRequest({ credit_card_number: '378282246310005' })),
      /cannot process/i
    );
  });

  it('rejects expired cards', () => {
    assert.throws(
      () => charge(chargeRequest({
        credit_card_expiration_year: 2000,
        credit_card_expiration_month: 1,
      })),
      /expired/i
    );
  });

  it('strips dashes from masked card numbers in expiry errors', () => {
    try {
      charge(chargeRequest({
        credit_card_number: '4242-4242-4242-4242',
        credit_card_expiration_year: 2000,
        credit_card_expiration_month: 1,
      }));
      assert.fail('expected expired card error');
    } catch (err) {
      assert.match(err.message, /ending 4242\)/);
    }
  });
});
