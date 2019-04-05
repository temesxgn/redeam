#!/bin/bash

echo ============= Building Development API =============
docker build -f ../Dockerfile.local -t redeam/book-api ..

echo ============= Composing Development API =============
docker-compose -f ../docker-compose.local.yml up -d

echo ============= Displaying Logs =============
docker-compose -f ../docker-compose.local.yml logs -f