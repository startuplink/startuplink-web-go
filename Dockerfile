# build app
FROM golang:1.13.5-alpine3.11 as builder

ENV \
    GO111MODULE=on \
    GOFLAGS="-mod=vendor"

ADD backend /build/backend
ADD .git /build/.git
WORKDIR /build/backend

COPY . .

RUN go build -o startuplink-web -ldflags="-s -w" .

# todo: add versioning for image

# pack app into working container
FROM alpine:3.11

COPY --from=builder /build/backend/startuplink-web /app/
COPY --from=builder /build/backend/static /app/static
COPY --from=builder /build/backend/template /app/template

EXPOSE $PORT
HEALTHCHECK --interval=30s --timeout=3s CMD curl --fail http://localhost:8080/ping || exit 1

WORKDIR /app
CMD ["./startuplink-web"]
