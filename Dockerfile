FROM grafana/xk6 AS builder
# COPY . .
RUN xk6 build  --with github.com/leetcode-golang-classroom/jwt=./jwt --output k6 && cp k6 $GOPATH/bin/k6
FROM alpine:latest
RUN RUN apk add --no-cache ca-certificates && \
    adduser -D -u 12345 -g 12345 k6
COPY --from=builder /go/bin/k6 /usr/bin/k6
USER 123456
CMD ["k6", "--version"]

