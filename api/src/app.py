from fastapi import FastAPI, Response
import fastapi
from loguru import logger
import redis.asyncio as _redis

app = FastAPI()
pool = _redis.ConnectionPool(host="workbench-redis-master", port=6379, protocol=3)
redis = _redis.Redis(connection_pool=pool)


@app.get("/healthz")
async def healthz(response: Response):
    try:
        await redis.ping()
        return {"status": "up"}
    except Exception as e:
        logger.error(f"Redis is down: {e}")
        response.status_code = fastapi.status.HTTP_503_SERVICE_UNAVAILABLE
        return {"status": "down"}
