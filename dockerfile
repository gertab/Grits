# Use an official Golang runtime as a parent image
FROM golang:alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go application inside the container
RUN go build -o grits

# Define the command to run your application
ENTRYPOINT ["./grits"]
