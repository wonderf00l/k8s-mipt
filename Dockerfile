FROM golang:1.23.6-alpine AS build

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . .

RUN go build -o application cmd/*.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/application .

ENTRYPOINT [ "./application" ]