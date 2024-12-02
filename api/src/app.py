from fastapi import FastAPI

app = FastAPI()


@app.get("/")
async def root():
    return {"message": "Hello World"}

@app.get("/users")
async def hello_world():
    return [{"name": "Alice"}, {"name": "Bob"}]
