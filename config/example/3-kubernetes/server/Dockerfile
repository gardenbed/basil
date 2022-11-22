FROM golang:1.19-alpine as builder
RUN apk add --no-cache git
WORKDIR /repo
COPY . .
RUN go build

FROM alpine:3.17
EXPOSE 8080
COPY --from=builder /repo/server /usr/local/bin/
RUN chown -R nobody:nogroup /usr/local/bin/server
USER nobody
CMD [ "server" ]
