FROM golang:1.19.0-alpine3.16 as builder
RUN apk --no-cache add tzdata
WORKDIR /go/src/app
COPY . .

ENV CGO_ENABLED=0 \
    TZ="Europe/Moscow" \
    GOOS=linux

# RUN adduser \
#     --disabled-password \
#     --gecos "" \
#     --home "/nohome" \
#     --shell "/sbin/nologin" \
#     --no-create-home \
#     --uid 65532 \
#     valhh

RUN mkdir -p ./bin && go build -o ./bin/ ./cmd/backup-agent/ ./internal/...

FROM scratch
ENV TZ="Europe/Moscow"
COPY --from=builder /go/src/app/bin/backup-agent ./
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# COPY --from=builder /etc/group /etc/group
# COPY --from=builder /etc/passwd /etc/passwd

# USER valhh:valhh

# probably not needed
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["./backup-agent"]

