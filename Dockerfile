# Build stage
FROM golang:1.18-alpine3.15 as builder
WORKDIR /build
COPY go.mod . 
COPY go.sum .
RUN go mod download

COPY . .
RUN go mod tidy &&\
    go build -o main .

# Run stage
FROM alpine:3.15
WORKDIR /projects
COPY --from=builder /build/main ./
EXPOSE 8000

ENTRYPOINT [ "./main" ]
