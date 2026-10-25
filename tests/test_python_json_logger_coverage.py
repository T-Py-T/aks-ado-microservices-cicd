"""Unit tests for JSON logging helpers in Python microservices."""

import json
import logging
import sys
import unittest
from pathlib import Path
from unittest.mock import patch

REPO_ROOT = Path(__file__).resolve().parents[1]


def _import_logger_module(service_name: str):
    service_dir = REPO_ROOT / "src" / service_name
    sys.path.insert(0, str(service_dir))
    try:
        import logger as logger_module  # noqa: WPS433 — dynamic import per service
    finally:
        sys.path.pop(0)
    return logger_module


class PythonJsonLoggerCoverageTests(unittest.TestCase):
    def test_recommendationservice_logger_emits_json_with_severity(self) -> None:
        logger_module = _import_logger_module("recommendationservice")
        logger = logger_module.getJSONLogger("recommendationservice-test")

        with patch.object(sys.stdout, "write") as mock_write:
            logger.info("hello from recommendation test")

        output = "".join(call.args[0] for call in mock_write.call_args_list)
        payload = json.loads(output.strip())
        self.assertEqual(payload["severity"], "INFO")
        self.assertEqual(payload["name"], "recommendationservice-test")
        self.assertIn("hello from recommendation test", payload["message"])

    def test_emailservice_formatter_uppercases_existing_severity(self) -> None:
        logger_module = _import_logger_module("emailservice")
        formatter = logger_module.CustomJsonFormatter(
            "%(timestamp)s %(severity)s %(name)s %(message)s"
        )
        record = logging.LogRecord(
            name="emailservice-test",
            level=logging.WARNING,
            pathname=__file__,
            lineno=1,
            msg="warn path",
            args=(),
            exc_info=None,
        )
        log_record: dict = {}
        formatter.add_fields(log_record, record, {"severity": "warning"})
        self.assertEqual(log_record["severity"], "WARNING")


if __name__ == "__main__":
    unittest.main()
