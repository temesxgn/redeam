#!/bin/bash

echo ============= Tearing Down Containers =============
docker-compose -f ../docker-compose.local.yml down

echo ============= Cleaning Up =============
rm -rf bin
docker system prune -f
docker volume prune -f
