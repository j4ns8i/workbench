FROM debian:12-slim AS python
RUN apt-get update && \
    apt-get install --no-install-suggests --no-install-recommends --yes python3.11

FROM python AS build-app
COPY --from=ghcr.io/astral-sh/uv:debian-slim /usr/local/bin/uv /usr/local/bin/uv
ENV UV_LINK_MODE=copy \
    UV_COMPILE_BYTECODE=1 \
    UV_PYTHON_DOWNLOADS=never \
    UV_PYTHON=python3.11 \
    UV_PROJECT_ENVIRONMENT=/app
COPY pyproject.toml /build/pyproject.toml
COPY uv.lock /build/uv.lock
RUN cd /build && uv sync --locked --no-dev --no-install-project

# Copy the virtualenv into a distroless image
FROM gcr.io/distroless/python3-debian12:debug
ENV PATH=/app/bin:$PATH
COPY --from=build-app --chown=nonroot:nonroot /app /app
COPY --chown=nonroot:nonroot src /app/src
USER nonroot
WORKDIR /app
STOPSIGNAL SIGINT
ENTRYPOINT ["fastapi", "run", "src/app.py"]
