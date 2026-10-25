/**
 * Unit tests for paymentservice charge validation paths.
 */

const { describe, it } = require('node:test');
const assert = require('node:assert/strict');
const path = require('node:path');

const charge = require(path.join(__dirname, '../src/paymentservice/charge.js'));

const validRequest = {
  amount: { currency_code: 'USD', units: 10, nanos: 0 },
  credit_card: {
    credit_card_number: '4432-8015-6152-0454',
    credit_card_expiration_month: 1,
    credit_card_expiration_year: 2039,
    credit_card_cvv: 672,
  },
};

describe('paymentservice charge', () => {
  it('returns a transaction id for valid Visa cards', () => {
    const result = charge(validRequest);
    assert.ok(result.transaction_id);
    assert.match(result.transaction_id, /^[0-9a-f-]{36}$/i);
  });

  it('rejects invalid card numbers', () => {
    assert.throws(
      () =>
        charge({
          ...validRequest,
          credit_card: {
            ...validRequest.credit_card,
            credit_card_number: '0000-0000-0000-0000',
          },
        }),
      /invalid/i,
    );
  });

  it('rejects expired cards', () => {
    assert.throws(
      () =>
        charge({
          ...validRequest,
          credit_card: {
            ...validRequest.credit_card,
            credit_card_expiration_year: 2000,
          },
        }),
      /expired/i,
    );
  });
});
