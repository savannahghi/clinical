# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.21 as builder

# Create and change to the app directory.
WORKDIR /app

# Install wkhtmltopdf in the builder stage.
RUN apt-get update && apt-get install -y wkhtmltopdf && apt-get clean

# Copy go.sum/go.mod and warm up the module cache.
COPY go.* ./
RUN go mod download

# Set the environment variable for Gin in release mode.
ENV GIN_MODE release

# Now copy the rest of the application's source code
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server github.com/savannahghi/clinical

# Use the official Alpine image for a lean production container.
FROM alpine:3 as production

# Install ca-certificates and other runtime dependencies for wkhtmltopdf.
RUN apk add --no-cache ca-certificates

# Copy the wkhtmltopdf binary and related files from the builder stage to the production image.
COPY --from=builder /usr/bin/wkhtmltopdf /usr/local/bin/

# Copy the Go binary to the production image from the builder stage.
COPY --from=builder /app/server /server

# Ensure your templates directory is correctly copied into the Docker image.
COPY --from=builder /app/templates /app/templates

# Set the working directory to where your binary and templates are
WORKDIR /app

# Run the web service on container startup.
CMD ["/server"]
