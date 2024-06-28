# build app
FROM golang:1.18-alpine3.15 as builder
RUN apk --no-cache add git

ENV \
    GO111MODULE=on

ADD backend /build/backend
ADD .git /build/.git
WORKDIR /build/backend

COPY . .

RUN go build -o startuplink-web -ldflags="-s -w" .

# todo: add versioning for image

# pack app into working container
FROM alpine:3.11

RUN apk upgrade && apk add --no-cache curl

COPY --from=builder /build/backend/startuplink-web /app/
COPY --from=builder /build/backend/static /app/static
COPY --from=builder /build/backend/template /app/template

EXPOSE $PORT

WORKDIR /app
CMD ["./startuplink-web"]
