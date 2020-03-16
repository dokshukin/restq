# Start from golang v1.12 base image
FROM golang:1.12 as builder

# Add Maintainer Info
LABEL maintainer="Ilya Dokshukin <dokshukin@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /go/src/github.com/dokshukin/restq
RUN apt-get update && \
  apt-get install -y unzip upx && \
  wget https://github.com/dokshukin/restq/archive/master.zip && \
  unzip master.zip -d /tmp/ && \
  mv /tmp/restq-master/* .  && \
  rm master.zip

# Download dependencies
RUN go get -d -v ./...

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build --ldflags "-s -w" -a -installsuffix cgo -o /go/bin/restq .
RUN upx /go/bin/restq


######## Start a new stage from scratch #######
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/bin/restq .

EXPOSE 8080

CMD ["./restq"]
