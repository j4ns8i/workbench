from datetime import datetime
from fastapi import FastAPI, Response
import fastapi
from fastapi.responses import StreamingResponse
from loguru import logger
import redis.asyncio as _redis

from . import config
from .logger import logger
from .models import (
    EventEnvelope,
    MessageEventData,
    NewMessageRequest,
)
from .store import events_key

app = FastAPI()
pool = _redis.ConnectionPool(host=config.redis_host, port=config.redis_port, protocol=3)
redis = _redis.Redis(connection_pool=pool)


@app.post("/api/v0/messages")
async def post_messages(req: NewMessageRequest, response: Response):
    logctx = logger.bind(topic=req.topic, store_key=events_key)
    data = MessageEventData(
        topic=req.topic,
        message=req.message,
        timestamp=datetime.now(),
    )
    env = EventEnvelope(
        version="workbench/v0",
        kind="message",
        data=data.model_dump_json(),
    )
    logctx.info("storing new message")
    try:
        await redis.xadd(events_key, env.model_dump())  # type: ignore
    except Exception as e:
        logctx.bind(error=str(e)).error("error storing new message")
        response.status_code = fastapi.status.HTTP_503_SERVICE_UNAVAILABLE
        return {"error": "service unavailable"}
    response.status_code = fastapi.status.HTTP_204_NO_CONTENT


async def redis_stream_listener(kind: str | None = None, topic: str | None = None):
    logctx = logger.bind(store_key=events_key)
    logctx.info("listening for events")
    while True:
        try:
            resp = await redis.xread({events_key: "$"}, block=0)
        except Exception as e:
            # TODO: detect client disconnection through keep alive messages, treat as normal
            # TODO: send event: error types?
            logctx.bind(error=str(e)).error("error while listening for event")
            break

        try:
            env = EventEnvelope.from_redis(resp[events_key.encode()][0][0][1])
        except Exception as e:
            logctx.bind(error=str(e), resp=str(resp)).error("error parsing event data")
            break

        if kind is not None and env.kind != kind:
            continue
        elif topic is not None and env.kind != "message":
            continue

        message = MessageEventData.model_validate_json(env.data)
        if topic is not None and message.topic != topic:
            continue
        yield f"data: {env.model_dump_json()}\n\n"


@app.get("/api/v0/events")
async def post_events(kind: str | None = None, topic: str | None = None):
    return StreamingResponse(
        redis_stream_listener(kind=kind, topic=topic),
        media_type="text/event-stream",
    )


@app.get("/healthz")
async def healthz(response: Response):
    logger.info("checking health")
    try:
        await redis.ping()
        return {"status": "up"}
    except Exception as e:
        logger.bind(error=str(e)).error(f"error pinging redis")
        response.status_code = fastapi.status.HTTP_503_SERVICE_UNAVAILABLE
        return {"status": "down"}
