from datetime import datetime
from typing import Any
from pydantic import BaseModel, Field, field_serializer


class EventEnvelope(BaseModel):
    version: str
    kind: str
    data: str

    @classmethod
    def from_redis(cls, env: dict[bytes, bytes]) -> "EventEnvelope":
        return cls(
            version=env[b"version"].decode(),
            kind=env[b"kind"].decode(),
            data=env[b"data"].decode(),
        )


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


class StreamEvent(BaseModel):
    id: str
    data: str
