# Build stage
FROM golang:latest AS build


# Set the working directory
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port on which the API will run
EXPOSE 8080

# Run the application
CMD ["./main"]