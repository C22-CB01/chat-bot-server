# Build stage
FROM golang:1.18-alpine3.15 as builder
WORKDIR /build

COPY go.mod . 
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o main .

# Run stage
FROM alpine:3.15 as prod
WORKDIR /projects
COPY --from=builder /build/main ./
EXPOSE 8080

CMD [ "./main" ]

FROM alpine:3.15 as dev
WORKDIR /projects
COPY --from=builder /build/main ./
EXPOSE 8080
COPY .env .
COPY firebase_key.json .

CMD [ "./main" ]
