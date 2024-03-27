# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.22-bullseye as builder

# Create and change to the app directory.
WORKDIR /app

# Copy go.sum/go.mod and warm up the module cache.
COPY go.* ./
RUN go mod download

# Set the environment variable for Gin in release mode.
ENV GIN_MODE release

# Now copy the rest of the application's source code
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server github.com/savannahghi/clinical

FROM debian:bullseye-slim as production

# Install necessary libraries for wkhtmltopdf.
RUN apt-get update && apt-get install -y --no-install-recommends \
    wkhtmltopdf \
    libstdc++6 \
    libx11-6 \
    libxrender1 \
    libxext6 \
    libfontconfig1 \
    fonts-dejavu \
    fonts-droid-fallback \
    fonts-freefont-ttf \
    fonts-liberation \
    libqt5webkit5 \
    libqt5widgets5 \
    libqt5gui5 \
    libqt5core5a \
    libqt5network5 \
    ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*


# Copy the Go binary to the production image from the builder stage.
COPY --from=builder /app/server /server

# Ensure your templates directory is correctly copied into the Docker image.
COPY --from=builder /app/templates /app/templates

# Set the working directory to where your binary and templates are.
WORKDIR /app

# Run the web service on container startup.
CMD ["/server"]
