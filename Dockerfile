# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.19 as builder

# Create and change to the app directory.
WORKDIR /app

# Install wkhtmltopdf in the builder stage.
RUN apt-get update && apt-get install -y wkhtmltopdf && apt-get clean

# Copy go.sum/go.mod and warm up the module cache (so that this
# rather long step can be cached if go.mod/go.sum don't change)
COPY go.* ./
RUN go mod download

# Set the environment variable for Gin in release mode.
ENV GIN_MODE release

# Now copy the rest
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server github.com/savannahghi/clinical

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3 as production

# Install ca-certificates for SSL.
RUN apk add --no-cache ca-certificates

# Copy the wkhtmltopdf binary from the builder stage to the production image.
# The path /usr/bin/wkhtmltopdf is typical for Debian-based installations; adjust if necessary.
COPY --from=builder /usr/bin/wkhtmltopdf /usr/local/bin/

# Copy the Go binary to the production image from the builder stage.
COPY --from=builder /app/server /server

# Run the web service on container startup.
CMD ["/server"]
