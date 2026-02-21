#!/bin/bash
# chmod +x build.sh

echo "Building image:"
docker image build -f Dockerfile -t forum .
echo ""
echo "Running container:"
docker container run -p 8080:8080 -d --name web-container forum

