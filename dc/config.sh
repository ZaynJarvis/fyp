#!/bin/sh
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"service":"abc","version":"123", "object_storage_path": "localhost:9000", "document_storage_path": "localhost:27017",  "rules": [{"field":"right_eye", "op": 2, "operand":"", "sample_rate": 1}]}' \
  http://localhost:8900
