#!/bin/bash

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    read -p "protoc-gen-go is not installed. Do you want to install it? (y/n) " choice
    case "$choice" in
        y|Y ) go install google.golang.org/protobuf/cmd/protoc-gen-go;;
        n|N ) echo "Please manually install protoc-gen-go by running 'go install google.golang.org/protobuf/cmd/protoc-gen-go'";;
        * ) echo "Invalid choice. Please enter y or n.";;
    esac
fi

# Check if protoc-gen-go-grpc is installed
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    read -p "protoc-gen-go-grpc is not installed. Do you want to install it? (y/n) " choice
    case "$choice" in
        y|Y ) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc;;
        n|N ) echo "Please manually install protoc-gen-go-grpc by running 'go install google.golang.org/grpc/cmd/protoc-gen-go-grpc'";;
        * ) echo "Invalid choice. Please enter y or n.";;
    esac
fi

# Check if protoc-gen-ts is installed
if ! command -v protoc-gen-ts &> /dev/null; then
    read -p "protoc-gen-ts is not installed. Do you want to install it? (y/n) " choice
    case "$choice" in
        y|Y ) npm i -g @protobuf-ts/plugin;;
        n|N ) echo "Please manually install protoc-gen-ts by running 'npm i -g @protobuf-ts/plugin'";;
        * ) echo "Invalid choice. Please enter y or n.";;
    esac
fi

# Check if all dependencies are installed, and if not, list the ones that aren't
if ! command -v protoc-gen-go &> /dev/null || ! command -v protoc-gen-go-grpc &> /dev/null || ! command -v protoc-gen-ts &> /dev/null; then
    echo "Some dependencies are not installed. Please install them and try again."
    exit 1
else
    echo "All dependencies are installed."
fi