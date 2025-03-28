import json
import sys
from datetime import timezone
import loguru


def _formatter(record):
    serialized = record["extra"]
    serialized["time"] = record["time"].astimezone(timezone.utc).isoformat()
    serialized["level"] = record["level"].name
    serialized["message"] = record["message"]
    record["extra"]["serialized"] = json.dumps(serialized)
    return "{extra[serialized]}" + "\n"

logger = loguru.logger
logger.configure(
    handlers=[
        {
            "sink": sys.stderr,
            "format": _formatter,
        },
    ]
)
