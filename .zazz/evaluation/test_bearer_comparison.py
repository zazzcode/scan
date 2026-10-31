"""Offline regression checks for invalid benchmark execution handling."""

import contextlib
import importlib.util
import io
import json
import os
from pathlib import Path
import sys
import tempfile
import unittest


spec = importlib.util.spec_from_file_location(
    "comparison", Path(__file__).with_name("run-bearer-comparison.py")
)
comparison = importlib.util.module_from_spec(spec)
spec.loader.exec_module(comparison)


class ExecutionTests(unittest.TestCase):
    def test_execution_validity(self):
        cases = [
            ("bearer-findings", '{"critical": []}', "", 1, True),
            ("bearer-worker", "{}", "error processing server.ts: failed worker", 0, False),
            ("bearer-empty", "{}", "no files to scan", 0, False),
            ("bearer-bad-json", "", "invalid configuration", 1, False),
            ("bearer-bad-shape", "[]", "", 0, False),
            ("scan-offline", '{"schema_version": 1}', "", 0, True),
            ("scan-failed", '{"schema_version": 1}', "", 1, False),
        ]
        for name, report, log, exit_code, valid in cases:
            with self.subTest(name=name), tempfile.TemporaryDirectory() as directory:
                output = Path(directory)
                code = "import sys; print(%r); print(%r, file=sys.stderr); sys.exit(%d)" % (
                    report, log, exit_code
                )
                with contextlib.redirect_stdout(io.StringIO()):
                    if valid:
                        comparison.run(name, [sys.executable, "-c", code], output, {"PATH": os.defpath})
                    else:
                        with self.assertRaises(RuntimeError):
                            comparison.run(name, [sys.executable, "-c", code], output, {"PATH": os.defpath})
                record = json.loads((output / (name + ".execution.json")).read_text())
                self.assertEqual(record["valid_execution"], valid)
                self.assertEqual(record["exit_code"], exit_code)


if __name__ == "__main__":
    unittest.main()
