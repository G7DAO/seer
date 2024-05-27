import os
import uuid

from bugout.app import Bugout


BUGOUT_BROOD_URL = os.environ.get("BUGOUT_BROOD_URL", "https://auth.bugout.dev")
BUGOUT_SPIRE_URL = os.environ.get("BUGOUT_SPIRE_URL", "https://spire.bugout.dev")
BUGOUT_SPIRE_EXTERNAL_URL = os.environ.get(
    "BUGOUT_SPIRE_EXTERNAL_URL", "https://spire.bugout.dev"
)


bugout_client = Bugout(
    brood_api_url=BUGOUT_BROOD_URL, spire_api_url=BUGOUT_SPIRE_EXTERNAL_URL
)

MOONSTREAM_MOONWORM_TASKS_JOURNAL_RAW = os.environ.get(
    "MOONSTREAM_MOONWORM_TASKS_JOURNAL_ID"
)

try:
    MOONSTREAM_MOONWORM_TASKS_JOURNAL_ID = uuid.UUID(
        MOONSTREAM_MOONWORM_TASKS_JOURNAL_RAW
    )
except ValueError:
    raise ValueError("MOONSTREAM_MOONWORM_TASKS_JOURNAL_ID must be a valid UUID")

if MOONSTREAM_MOONWORM_TASKS_JOURNAL_ID is None:
    raise ValueError(
        "MOONSTREAM_MOONWORM_TASKS_JOURNAL_ID environment variable must be set"
    )

MOONSTREAM_ADMIN_ACCESS_TOKEN = os.environ.get("MOONSTREAM_ADMIN_ACCESS_TOKEN")

if MOONSTREAM_ADMIN_ACCESS_TOKEN is None:
    raise ValueError("MOONSTREAM_ADMIN_ACCESS_TOKEN environment variable must be set")
