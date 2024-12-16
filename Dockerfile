# Dockerfile
# Build stage
FROM golang:1.23.2 AS build

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code to the image
COPY . ./

# Compile the binary
RUN go build -o main ./cmd

# Runtime stage
FROM golang:1.23.2

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/main .

# Copy the source code files from the build stage (including cmd/load_csv/load.go)
COPY --from=build /app /app

# Expose the application's port
EXPOSE 8080

# Command to run the application
CMD ["./main"]
