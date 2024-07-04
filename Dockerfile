FROM golang:1.22-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN cd cmd/api && go build -o /uilabs-api

FROM alpine:latest

COPY --from=build /uilabs-api /uilabs-api
WORKDIR /

ENTRYPOINT ["/uilabs-api"]
