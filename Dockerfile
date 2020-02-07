FROM golang:1.12.4-alpine3.9 as builder

WORKDIR /go/src/
ADD header.go .
ADD table.html .

RUN go build -v -o header

FROM alpine:3.9
COPY --from=builder /go/src/header /header
COPY --from=builder /go/src/table.html /table.html
ENTRYPOINT ["/header"]