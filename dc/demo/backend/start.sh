#!/bin/sh
docker run --name=frc -p 8000:8000 -v $(pwd)/demo/backend:/root/face_recognition/backend fr
