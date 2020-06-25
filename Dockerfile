FROM golang:alpine as builder
WORKDIR /go/src/moxy
COPY . .
RUN go get -d -v ./...
RUN go build -o /moxy .


FROM alpine:latest
WORKDIR /moxy
COPY --from=builder /moxy .
ENV MOXY_PORT=8080
ENV MOXY_MOCK_FILE=""
CMD ["sh", "-c", "./moxy -backend=${MOXY_BACKEND} -port=${MOXY_PORT} -mock-file=${MOXY_MOCK_FILE}"]
