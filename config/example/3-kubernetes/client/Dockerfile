FROM golang:1.18-alpine as builder
RUN apk add --no-cache git
WORKDIR /repo
COPY . .
RUN go build

FROM alpine:3.16
COPY --from=builder /repo/client /usr/local/bin/
RUN chown -R nobody:nogroup /usr/local/bin/client
USER nobody
CMD [ "client" ]
