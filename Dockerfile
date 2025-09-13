FROM golang:1.17 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /xauth-build

FROM alpine

WORKDIR /app

RUN apk add --no-cache bash

COPY --from=build /xauth-build .

COPY static ./static

COPY wait-for-it.sh /usr/local/bin/wait-for-it.sh

RUN chmod +x /usr/local/bin/wait-for-it.sh

CMD ["/xauth-build"]
