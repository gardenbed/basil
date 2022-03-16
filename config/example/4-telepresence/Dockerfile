FROM golang:1.18-alpine as builder
RUN apk add --no-cache git
WORKDIR /repo
COPY . .
RUN go build -o app

FROM alpine:3.15
EXPOSE 8080
COPY --from=builder /repo/app /usr/local/bin/
RUN chown -R nobody:nogroup /usr/local/bin/app
USER nobody
CMD [ "app" ]
