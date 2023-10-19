# 1. Build executable binary

FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && \
    apk add git && \
    apk add build-base upx

WORKDIR /src/bookbase
COPY . .
# Fetch dependencies using `go get`.
# Build the binary and run it.
RUN go build -o /go/bin/bookbase
RUN upx /go/bin/bookbase

# 2. Build a small image
FROM alpine
# Copy out static executable.
RUN apk update && apk add --no-cache vips-dev
COPY --from=builder /go/bin/bookbase /go/bin/bookbase

# Run the binary.
ENTRYPOINT [ "/go/bin/bookbase" ]
