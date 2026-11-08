FROM golang:alpine AS build-env

WORKDIR /src

# Download dependencies first to make better use of the layer cache.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN go build -trimpath -ldflags="-s -w" -v -o /mp .


# The binary is fully static and needs no OS packages, so ship it on an empty
# base image. This avoids inheriting (unfixable) CVEs from alpine packages.
FROM scratch

WORKDIR /app/

COPY --from=build-env /mp /app/

USER 65534:65534
ENTRYPOINT [ "/app/mp", "-file", "/data/words.txt" ]
