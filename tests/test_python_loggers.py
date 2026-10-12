import json
import logging
import os
import sys
import unittest

EMAIL_SERVICE_ROOT = os.path.join(os.path.dirname(__file__), "..", "src", "emailservice")
sys.path.insert(0, EMAIL_SERVICE_ROOT)

from logger import CustomJsonFormatter, getJSONLogger  # noqa: E402


class PythonLoggerTests(unittest.TestCase):
    def test_custom_json_formatter_emits_severity_and_timestamp(self) -> None:
        formatter = CustomJsonFormatter("%(timestamp)s %(severity)s %(name)s %(message)s")
        record = logging.LogRecord(
            name="test",
            level=logging.INFO,
            pathname=__file__,
            lineno=1,
            msg="hello",
            args=(),
            exc_info=None,
        )
        payload = json.loads(formatter.format(record))
        self.assertEqual(payload["severity"], "INFO")
        self.assertIn("timestamp", payload)
        self.assertEqual(payload["message"], "hello")

    def test_get_json_logger_does_not_propagate(self) -> None:
        logger = getJSONLogger("emailservice-test-logger")
        self.assertFalse(logger.propagate)


if __name__ == "__main__":
    unittest.main()
