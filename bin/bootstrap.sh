#!/bin/bash

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Docker not found. Installing Docker..."
    # Install Docker
    # Add your installation command here
    # Check for and install docker
    echo "Docker installed successfully."
else
    echo "Docker is already installed."
fi

# Check if Terraform is installed
if ! command -v terraform &> /dev/null; then
    echo "Terraform not found. Installing Terraform..."
    # Install Terraform
    # Add your installation command here
    echo "Terraform installed successfully."
else
    echo "Terraform is already installed."
fi

# Check if Golang is installed
if ! command -v go &> /dev/null; then
    echo "Golang not found. Installing Golang..."
    # Install Golang
    # Add your installation command here
    echo "Golang installed successfully."
else
    echo "Golang is already installed."
fi
