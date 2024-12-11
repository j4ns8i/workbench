from pydantic import BaseModel, Field

class NewEventRequest(BaseModel):
    topic: str = Field(min_length=1, max_length=127)
    message: str = Field(min_length=1, max_length=255)
