import os


redis_host = os.getenv("REDIS_HOST")
redis_port = os.getenv("REDIS_PORT")

events_keep_alive_interval = int(os.getenv("EVENTS_KEEP_ALIVE_INTERVAL", 10))
