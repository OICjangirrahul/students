FROM golang:1.22-alpine

WORKDIR /app

# Install required system packages
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main cmd/main.go

# Expose port 8082
EXPOSE 8082

# Command to run the executable
CMD ["./main"] 