FROM golang:alpine AS build-env

WORKDIR /go/src/github.com/martinplaner/mp

RUN apk --no-cache add git curl

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN go get -v -t -d ./...
RUN go get github.com/GeertJohan/go.rice/rice
RUN go generate -v .
RUN go build -ldflags="-s -w" -v .


FROM alpine

WORKDIR /app/

COPY --from=build-env /go/src/github.com/martinplaner/mp/mp /app/

USER nobody
ENTRYPOINT [ "/app/mp", "-file", "/data/words.txt" ]
