#BUILD
FROM golang:1.22-alpine AS build_base

#RUN apk add --no-cache git
WORKDIR /tmp/build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./out/tg-update ./cmd/tg-update

#DEPLOY
FROM alpine:3.19 
RUN apk add ca-certificates

COPY --from=build_base /tmp/build/out/tg-update /app/tg-update

EXPOSE 3000

WORKDIR /app
ENTRYPOINT ["/app/tg-update"]