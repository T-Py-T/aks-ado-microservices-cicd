import re
import unittest
from pathlib import Path


REPOSITORY_ROOT = Path(__file__).resolve().parents[1]
SHA256_DIGEST = re.compile(r"@sha256:[0-9a-f]{64}$")
PINNED_ACTION = re.compile(r"\buses:\s+[^\s@]+@[0-9a-f]{40}(?:\s+#.*)?$")
DISALLOWED_GITHUB_EVENTS = re.compile(
    r"^  (?:push|schedule|workflow_dispatch|workflow_call|workflow_run|repository_dispatch):",
    re.MULTILINE,
)


def external_parent_image(from_line: str) -> str | None:
    tokens = from_line.split()
    if not tokens or tokens[0].upper() != "FROM":
        return None

    image_index = 2 if len(tokens) > 1 and tokens[1].startswith("--platform=") else 1
    image = tokens[image_index]
    if image.lower() in {"base", "build", "builder", "publish"}:
        return None
    return image


class RepositoryPolicyTests(unittest.TestCase):
    def test_github_workflows_only_run_for_pull_requests(self) -> None:
        workflow_paths = sorted((REPOSITORY_ROOT / ".github" / "workflows").glob("*.y*ml"))
        self.assertTrue(workflow_paths, "at least one GitHub workflow is required")

        for path in workflow_paths:
            contents = path.read_text()
            self.assertRegex(contents, r"(?m)^on:\s*$")
            self.assertRegex(contents, r"(?m)^  pull_request:\s*$")
            self.assertNotRegex(contents, DISALLOWED_GITHUB_EVENTS)

    def test_github_actions_are_pinned_to_commits(self) -> None:
        workflow_paths = sorted((REPOSITORY_ROOT / ".github" / "workflows").glob("*.y*ml"))
        action_lines = [
            line.strip()
            for path in workflow_paths
            for line in path.read_text().splitlines()
            if "uses:" in line
        ]
        self.assertTrue(action_lines, "at least one GitHub Action is required")
        for line in action_lines:
            self.assertRegex(line, PINNED_ACTION)

    def test_external_docker_parents_are_digest_pinned(self) -> None:
        dockerfiles = sorted(REPOSITORY_ROOT.glob("src/**/Dockerfile*"))
        self.assertTrue(dockerfiles, "service Dockerfiles are required")

        for path in dockerfiles:
            for line in path.read_text().splitlines():
                image = external_parent_image(line)
                if image is not None:
                    self.assertRegex(image, SHA256_DIGEST, f"unpinned parent in {path}: {image}")

    def test_azure_pipeline_has_no_automatic_trigger(self) -> None:
        pipeline = (REPOSITORY_ROOT / "azure-pipelines.yml").read_text()
        self.assertRegex(pipeline, r"(?m)^trigger: none$")
        self.assertRegex(pipeline, r"(?m)^pr: none$")

    def test_redis_image_is_versioned_and_digest_pinned(self) -> None:
        manifest = (REPOSITORY_ROOT / "deployment-service.yaml").read_text()
        self.assertRegex(manifest, r"image: redis:\d+\.\d+\.\d+-alpine@sha256:[0-9a-f]{64}")

    def test_manifest_uses_native_health_probes(self) -> None:
        manifest = (REPOSITORY_ROOT / "deployment-service.yaml").read_text()
        self.assertNotIn("/bin/grpc_health_probe", manifest)

    def test_security_scan_exceptions_are_narrow(self) -> None:
        config = (REPOSITORY_ROOT / ".grype.yaml").read_text()
        self.assertEqual(config.count("  - vulnerability:"), 4)
        self.assertEqual(config.count("      name: python"), 4)
        self.assertEqual(config.count("      version: 3.14.7"), 4)
        self.assertEqual(config.count("      type: binary"), 4)


if __name__ == "__main__":
    unittest.main()
