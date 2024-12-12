from datetime import datetime
from fastapi import FastAPI, Response
import fastapi
from loguru import logger
import redis.asyncio as _redis

from . import config
from .logger import logger
from .models import MessageEventData, NewMessageRequest
from .store import topic_key

app = FastAPI()
pool = _redis.ConnectionPool(host=config.redis_host, port=config.redis_port, protocol=3)
redis = _redis.Redis(connection_pool=pool)


@app.post("/api/v0/messages")
async def post_messages(req: NewMessageRequest, response: Response):
    key = topic_key(req.topic)
    logctx = logger.bind(key=key)
    data = MessageEventData(
        topic=req.topic, message=req.message, timestamp=datetime.now()
    )
    try:
        logctx.bind(data=data.model_dump()).info("storing new message")
        await redis.xadd(key, data.model_dump())  # type: ignore
    except Exception as e:
        logctx.bind(error=str(e)).error("error storing new message")
        response.status_code = fastapi.status.HTTP_503_SERVICE_UNAVAILABLE
        return {"error": "service unavailable"}
    response.status_code = fastapi.status.HTTP_204_NO_CONTENT


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
