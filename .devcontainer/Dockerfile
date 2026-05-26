FROM golang:1.26-alpine AS builder

WORKDIR /build

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o gostatus ./main.go

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /build/gostatus .
COPY --from=builder /build/assets ./assets

ENV PORT=8080
EXPOSE ${PORT}

ENTRYPOINT ["/app/gostatus"]
CMD        ["-config", "/app/config.toml", "-port", "8080"]
