# Build stage
FROM golang:1.18 as builder
WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o server

FROM alpine:3 as prod
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server /server
EXPOSE 8080

CMD [ "/server" ]

FROM alpine:3 as dev
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server /server
EXPOSE 8080
COPY firebase_key.json ./

CMD [ "/server" ]
