from datetime import datetime
from typing import Any
from pydantic import BaseModel, Field, field_serializer


class EventEnvelope(BaseModel):
    version: str
    kind: str
    data: Any


class NewMessageRequest(BaseModel):
    topic: str = Field(min_length=1, max_length=127)
    message: str = Field(min_length=1, max_length=255)


class MessageEventData(BaseModel):
    topic: str = Field(min_length=1, max_length=127)
    message: str = Field(min_length=1, max_length=255)
    timestamp: datetime

    @field_serializer("timestamp")
    def serialize_timestamp(self, timestamp: datetime) -> str:
        return timestamp.isoformat()
