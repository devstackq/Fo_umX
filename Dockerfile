FROM golang:alpine as builder
# Install make and certificates
RUN apk --no-cache add tzdata zip ca-certificates make git
# Make repository path
RUN mkdir -p /go/src/github.com/devstackq/Forum-X
WORKDIR /go/src/github.com/devstackq/Forum-X
# Copy Makefile first, it will save time during development.
COPY ./Makefile ./Makefile
# Install deps
RUN make deps
# Copy all project files
ADD . .
# Generate a binary
RUN make bin

# Second (final) stage, base image is scratch
FROM scratch
# Copy statically linked binary
COPY --from=builder /go/src/github.com/devstackq/Forum-X/app-linux-amd64 /app
# Copy SSL certificates, eventhough we don't need it for this example
# but if you decide to talk to HTTPS sites, you'll need this, you'll thank me later.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Notice "CMD", we don't use "Entrypoint" because there is no OS
CMD [ "/app" ]