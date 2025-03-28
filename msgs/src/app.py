import asyncio
from datetime import datetime
from fastapi import FastAPI, Response
import fastapi
from fastapi.responses import StreamingResponse
from loguru import logger
from redis.asyncio import Redis, ConnectionPool

from . import config
from .logger import logger
from .models import (
    EventEnvelope,
    MessageEventData,
    NewMessageRequest,
    StreamEvent,
)
from .store import events_key

app = FastAPI()
pool = ConnectionPool(host=config.redis_host, port=config.redis_port, password=config.redis_password, protocol=3)
redis = Redis(connection_pool=pool)


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


async def read_stream_event(
    redis: Redis,
    key: str,
    from_id: str,
    kind: str | None = None,
    topic: str | None = None,
) -> StreamEvent:
    while True:
        # TODO: when another topic is published to, this no longer returns
        #       messages even for the desired topic
        resp = await redis.xread({key: from_id}, block=0)

        (id, data) = resp[events_key.encode()][0][0]
        id = id.decode()
        from_id = id

        envelope = EventEnvelope.from_redis(data)

        if kind is not None and envelope.kind != kind:
            continue
        elif topic is not None and envelope.kind != "message":
            continue

        message = MessageEventData.model_validate_json(envelope.data)
        if topic is not None and message.topic != topic:
            continue
        return StreamEvent(id=id, data=envelope.model_dump_json())


async def stream_redis_events(
    keep_alive_interval: int,
    kind: str | None = None,
    topic: str | None = None,
):
    logctx = logger.bind(store_key=events_key)
    logctx.info("streaming events")

    def new_keep_alive_task(interval: int):
        return asyncio.create_task(asyncio.sleep(interval))

    def new_read_stream_event_task(from_id: str):
        return asyncio.create_task(
            read_stream_event(redis, events_key, from_id, kind, topic)
        )

    read_stream_event_task = new_read_stream_event_task("$")
    while True:
        keep_alive_task = new_keep_alive_task(keep_alive_interval)
        done, _ = await asyncio.wait(
            [keep_alive_task, read_stream_event_task],
            return_when=asyncio.FIRST_COMPLETED,
        )

        if read_stream_event_task in done:
            if e := read_stream_event_task.exception():
                logctx.bind(error=str(e)).error("error reading stream event")
                return
            stream_event = read_stream_event_task.result()
            yield f"data: {stream_event.data}\n\n"
            read_stream_event_task = new_read_stream_event_task(stream_event.id)

        elif keep_alive_task in done:
            yield ": keep-alive\n\n"
        else:
            keep_alive_task.cancel()


@app.get("/api/v0/events")
async def post_events(kind: str | None = None, topic: str | None = None):
    return StreamingResponse(
        stream_redis_events(
            config.events_keep_alive_interval,
            kind=kind,
            topic=topic,
        ),
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
