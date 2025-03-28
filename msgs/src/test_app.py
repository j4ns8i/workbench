from fastapi.testclient import TestClient
from .app import app

client = TestClient(app)


def test_ping_redis_healthy(mocker):
    async def mock_ping():
        return True

    mocker.patch("redis.asyncio.Redis.ping", side_effect=mock_ping)
    response = client.get("/healthz")
    assert response.status_code == 200
    assert response.json() == {"status": "up"}


def test_ping_redis_unhealthy(mocker):
    async def mock_ping():
        raise Exception("mock error")

    mocker.patch("redis.asyncio.Redis.ping", side_effect=mock_ping)
    response = client.get("/healthz")
    assert response.status_code == 503
    assert response.json() == {"status": "down"}
