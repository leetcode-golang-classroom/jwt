FROM grafana/xk6 AS builder
RUN xk6 build --with github.com/leetcode-golang-classroom/jwt=./jwt --output k6 
RUN cp k6 $GOPATH/bin/k6
FROM alpine:3.13
RUN  adduser -D -u 12345 -g 12345 k6
COPY --from=builder /go/bin/k6 /usr/bin/k6
USER 12345
CMD ["k6", "--version"]

