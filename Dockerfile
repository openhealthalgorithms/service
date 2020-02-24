FROM alpine:edge AS builder
RUN apk update
RUN apk upgrade
RUN apk add --update go=1.13.4-r2 gcc=9.2.0-r5 g++=9.2.0-r5

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -a -tags netgo -installsuffix cgo -o ohas .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/ohas .

# Build a small image
FROM alpine:edge

COPY ./contents/. /contents

COPY ohas.toml /etc/ohas/ohas.toml

RUN mkdir /ohaslogs

RUN touch /ohaslogs/ohas-logs.db

COPY --from=builder /dist/ohas /

# Command to run
ENTRYPOINT ["/ohas"]
