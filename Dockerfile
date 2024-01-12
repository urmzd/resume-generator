# Start the first stage using the Go image for building the Go application
FROM golang:1.22-rc-bookworm as builder

# Set the working directory in the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the Go source code and other necessary files into the container at /app
# This includes all necessary files from the 'pkg', 'cmd' directories, and the 'main.go' file
COPY pkg/ pkg/
COPY cmd/ cmd/
COPY main.go .

# Compile the Go application
RUN go build -o generate-resumes main.go

# Start the second stage using your custom base image
FROM urmzd/generate-resumes-base:24.01.11 as base

# Set the working directory in the container
WORKDIR /app

# Copy the built Go binary and any other necessary files from the builder stage
COPY --from=builder /app/generate-resumes .

# Define the container's entrypoint as the application
ENTRYPOINT [ "./generate-resumes" ]
