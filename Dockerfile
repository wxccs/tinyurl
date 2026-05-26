FROM node:22-alpine AS frontend
COPY web/ /build/web/
WORKDIR /build/web
RUN npm ci && npm run build

FROM golang:1.26-alpine AS builder
ARG TARGETOS
ARG TARGETARCH
ARG VERSION="0.0.0"
ARG BUILD_TIME="1970-01-01T00:00:00Z"
ARG SHA="unknown"

COPY . /go/src/github.com/wxccs/tinyurl
COPY --from=frontend /build/web/dist /go/src/github.com/wxccs/tinyurl/web/dist
WORKDIR /go/src/github.com/wxccs/tinyurl

RUN go mod download && go mod verify

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -trimpath \
    -ldflags="-extldflags \"-static\" -X 'github.com/wxccs/tinyurl/global.Version=${VERSION}' -X 'github.com/wxccs/tinyurl/global.BuildTime=${BUILD_TIME}' -X 'github.com/wxccs/tinyurl/global.GitCommit=${SHA}'" \
    -o /bin/tinyurl

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /bin/tinyurl /usr/local/bin/tinyurl

RUN adduser -D -H appuser
USER appuser

EXPOSE 8080
ENTRYPOINT ["tinyurl"]
