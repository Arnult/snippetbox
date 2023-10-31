FROM golang:1.21-alpine AS build-stage

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /app
COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o snippetbox ./cmd/snippetbox

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app

COPY --from=build-stage /app/snippetbox /app/snippetbox
COPY --from=build-stage /app/ui/ /app/ui/

ENTRYPOINT ["/app/snippetbox"]