FROM golang:1.15 AS builder

# enable Go modules support
ENV GO111MODULE=on
WORKDIR /app

# manage dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM scratch
COPY --from=builder /app/ForumX /app/
EXPOSE 8181
ENTRYPOINT ["/app/ForumX"]
# Copy src code from the host and compile it
