FROM golang:1.21 AS build

# Create a directory for the files
RUN mkdir -p /app

# Copy the project files to the /app directory
COPY . /app

# Change the current working directory to /app
WORKDIR /app

# Run go mod download
RUN go mod download

# Build main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o bicep-docs ./cmd/bicep-docs/main.go

# Final image
FROM alpine:3.18

# Install bash
RUN apk add --no-cache bash

# Copy the binary and entrypoint.sh from the build stage
COPY --from=build /app/bicep-docs /app/bicep-docs
COPY --from=build /app/entrypoint.sh /app/entrypoint.sh

# Make entrypoint.sh executable
RUN chmod +x /app/entrypoint.sh

# Set the entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]
