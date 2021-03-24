FROM golang:1.13-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates busybox

WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app github.com/thlcodes/minrevpro/cmd/minrevpro
RUN stat -c "%s" app | awk '{printf "%.2f MB", $1/(1024*1024)}'
RUN chmod +x app

FROM busybox

ENV PORT 8080
ENV HOST ""
ENV TARGET ""
ENV DEBUG false
ENV SECRET ""

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/app /app

CMD ["sh", "-c", "/app --port \"$PORT\" --host \"$HOST\" --target \"$TARGET\" --secret \"$SECRET\" --debug=\"$DEBUG\""]