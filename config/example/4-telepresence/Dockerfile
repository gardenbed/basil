FROM golang:1.22-alpine as builder
RUN apk add --no-cache git
WORKDIR /repo
COPY . .
RUN go build -o app

FROM alpine:3.19
EXPOSE 8080
COPY --from=builder /repo/app /usr/local/bin/
RUN chown -R nobody:nogroup /usr/local/bin/app
USER nobody
CMD [ "app" ]
