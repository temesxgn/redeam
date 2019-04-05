#!/bin/bash

echo ============= Building Production API =============
docker build -f ../Dockerfile -t redeam/book-api ..

echo ============= Composing Production API =============
docker-compose -f ../docker-compose.yml up -d

echo ============= Displaying Logs =============
docker-compose -f ../docker-compose.yml logs -f