# syntax=docker/dockerfile:1.11

ARG GO_VERSION=1.25
ARG ALPINE_VERSION=3.22

FROM golang:${GO_VERSION}-alpine AS base
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download

FROM base AS build
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,target=. \
    go build -o /bin/api -ldflags "-s -w" main.go

FROM alpine:${ALPINE_VERSION} AS image
ENV GIN_MODE=release
COPY --from=build /bin/api /bin/api
EXPOSE 3000
ENTRYPOINT [ "/bin/api" ]
