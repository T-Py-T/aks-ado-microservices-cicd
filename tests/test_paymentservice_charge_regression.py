import subprocess
import unittest
from pathlib import Path

REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
PAYMENT_SERVICE_DIR = REPOSITORY_ROOT / "src" / "paymentservice"

CHARGE_NODE_SCRIPT = r"""
const charge = require('./charge');

const validRequest = {
  amount: { currency_code: 'USD', units: 42, nanos: 500000000 },
  credit_card: {
    credit_card_number: '4432-8015-6152-0454',
    credit_card_expiration_month: 1,
    credit_card_expiration_year: 2039,
  },
};

const result = charge(validRequest);
if (!result.transaction_id) {
  throw new Error('missing transaction_id');
}

let rejected = false;
try {
  charge({
    ...validRequest,
    credit_card: {
      ...validRequest.credit_card,
      credit_card_number: '3782-822463-10005',
    },
  });
} catch (error) {
  rejected = error.code === 400;
}
if (!rejected) {
  throw new Error('expected AMEX to be rejected with code 400');
}
"""


class PaymentServiceChargeRegressionTests(unittest.TestCase):
    def test_charge_accepts_visa_and_rejects_amex(self) -> None:
        completed = subprocess.run(
            ["node", "-e", CHARGE_NODE_SCRIPT],
            cwd=PAYMENT_SERVICE_DIR,
            check=False,
            capture_output=True,
            text=True,
        )
        self.assertEqual(
            completed.returncode,
            0,
            msg=completed.stderr or completed.stdout,
        )


if __name__ == "__main__":
    unittest.main()
