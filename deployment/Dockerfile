# Start from the official Golang image to build our source code.
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Copy the local modules and the modified go.mod file
COPY ./deployment/docker-build-context /app

# List contents after copy
RUN echo "Contents of /app after copy:" && ls -la

# If go.mod exists, print its contents
RUN if [ -f go.mod ]; then cat go.mod; else echo "go.mod not found"; fi

# Check if memstore directory exists and list its contents
RUN if [ -d memstore ]; then echo "memstore directory contents:" && ls -la memstore; else echo "memstore directory not found"; fi

# Print current module path
RUN go list -m

# Download all the dependencies
RUN go mod download

# List all modules
RUN go list -m all

# Print go.mod contents again
RUN cat go.mod

# List contents after copying everything
RUN echo "Contents of /app after copying everything:" && ls -la

# Check memstore directory again
RUN if [ -d memstore ]; then echo "memstore directory contents after full copy:" && ls -la memstore; else echo "memstore directory not found after full copy"; fi

# Attempt to build with verbose output
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o memorianexus ./cmd

# Start a new stage from scratch for a lightweight final image
FROM alpine:latest

WORKDIR /app/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/memorianexus /app/

RUN chmod +x /app/memorianexus

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD /app/memorianexus