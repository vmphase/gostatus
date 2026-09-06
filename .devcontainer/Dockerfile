FROM golang:1.26-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o gostatus ./cmd/gostatus

ENV PORT=8080
EXPOSE ${PORT}

ENTRYPOINT ["/app/gostatus"]
