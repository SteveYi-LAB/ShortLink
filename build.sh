#!/bin/bash

echo "Starting for build!"

echo "Linux amd64"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/SteveYi-ShortLink_Linux_Intel64
echo "macOS amd64"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/SteveYi-ShortLink_macOS_Intel64
echo "Windows amd64"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/SteveYi-ShortLink_Windows_Intel64.exe