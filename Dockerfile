FROM golang:1.22 AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cmd/apiserver/bin/main ./cmd/apiserver/

FROM alpine:latest
WORKDIR /apiserver

RUN mkdir /apiserver/logs
COPY --from=build /build/cmd/apiserver/bin/main .
COPY --from=build /build/migrations /apiserver/migrations
CMD ["/apiserver/main"]
