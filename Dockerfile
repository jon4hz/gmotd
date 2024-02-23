# Use a Debian base image
FROM debian:latest

# Avoid prompts from apt
ENV DEBIAN_FRONTEND=noninteractive

# Install necessary packages
RUN apt-get update && apt-get install -y \
    build-essential \
    libpam0g-dev \
    golang-go \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory in the container
WORKDIR /app

# Copy your module source code into the container
COPY . /app

# Build the Go part of your module (adjust the package name as needed)
RUN go build -buildmode=c-shared -o gmotd.so

# Use this command to run a shell when the container starts, for testing.
# You may want to adjust this to automatically run tests or a specific command.
CMD ["/bin/bash"]
