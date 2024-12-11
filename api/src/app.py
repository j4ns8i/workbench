from fastapi import FastAPI, Response
import fastapi
from loguru import logger
import redis.asyncio as _redis
import xxhash

from . import config
from .logger import logger
from .models import NewEventRequest

app = FastAPI()
pool = _redis.ConnectionPool(host=config.redis_host, port=config.redis_port, protocol=3)
redis = _redis.Redis(connection_pool=pool)


@app.post("/api/v0/events")
async def counter(req: NewEventRequest, response: Response):
    key = xxhash.xxh64_hexdigest(req.topic)
    logctx = logger.bind(key=key)
    try:
        await redis.hset(key, mapping={"message": req.message})  # type: ignore
    except Exception as e:
        logctx.bind(error=str(e)).error("error storing new event")
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
