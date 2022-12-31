FROM golang:1.19 AS build
WORKDIR /go/src/github.com/85yago/go_small_chat/app
COPY . .
RUN go mod download -x
RUN GOOS=linux CGO_ENABLED=0 go build -tags=release -ldflags="-s -w" -trimpath -o main .

FROM golang:1.19-bullseye AS development
RUN apt-get update && apt-get install -y \
    ca-certificates \
    gcc \
    git \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

FROM alpine:3.16 AS release
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /go/src/github.com/85yago/go_small_chat/app/ip_white_list .
COPY --from=build /go/src/github.com/85yago/go_small_chat/app/main .
CMD ["./main"]
