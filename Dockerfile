FROM golang:alpine AS build-env

WORKDIR /go/src/github.com/martinplaner/mp

RUN apk --no-cache add git curl

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN go build -ldflags="-s -w"


FROM alpine

WORKDIR /app/

ENV GIN_MODE=release

COPY --from=build-env /go/src/github.com/martinplaner/mp/mp /app/
COPY --from=build-env /go/src/github.com/martinplaner/mp/templates/*.* /app/templates/

USER nobody
ENTRYPOINT [ "/app/mp", "-file", "/data/words.txt" ]
