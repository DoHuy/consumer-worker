FROM golang:1.14.1 as builder
ARG PROJECT=vtp-consumer
ARG COMMIT_HASH_SHORT
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o main \
    -ldflags="-X '$PROJECT/common.ServiceCode=vtp-consumer'\
              -X '$PROJECT/common.CommitHashShort=$COMMIT_HASH_SHORT'"

FROM golang:1.14.1
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]