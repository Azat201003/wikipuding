FROM golang:latest AS builder

WORKDIR /build/
COPY src .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main ./main.go

FROM alpine

WORKDIR /app/src
COPY --from=builder /build/main .
RUN chmod +x main

EXPOSE 8080

CMD /app/src/main
