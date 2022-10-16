FROM golang:1.19.2-alpine3.16 AS build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src/httpserver

COPY ./httpserver/* ./
RUN go build -o /bin/httpserver

FROM alpine:3.16
COPY --from=build /bin/httpserver /bin/httpserver
ENTRYPOINT ["/bin/httpserver"]


