FROM golang:1.24 AS build

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.go-build \
    CGO_ENABLED=0 go build -o /go/bin/product-store

FROM gcr.io/distroless/base-debian12:debug
COPY --from=build /go/bin/product-store /usr/local/bin/product-store
ENTRYPOINT [ "/usr/local/bin/product-store" ]
