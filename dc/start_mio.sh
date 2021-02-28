#!/bin/sh
docker run -p 9000:9000 \
-e "MINIO_ACCESS_KEY=liuz0063" \
-e "MINIO_SECRET_KEY=12345678" \
-v $(pwd)/data:/data \
minio/minio server /data
