FROM golang:1.21 AS builder

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o irc-logger .


FROM alpine:3.18
# For timezone data
RUN apk --no-cache add tzdata

WORKDIR /app
COPY --from=builder /src/irc-logger /app/irc-logger
EXPOSE 8080

CMD ["/app/irc-logger"]
