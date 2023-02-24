# Start from golang base image
FROM golang:1.12.0-alpine3.9

# Install git.
# Git is required for fetching the dependencies.
RUN mkdir /app

# Setup folders
ADD . /app
WORKDIR /app

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 7000

# Run the executable
CMD [ "/app/cmd/api/main" ]