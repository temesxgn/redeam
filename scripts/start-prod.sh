#!/bin/bash

echo ============= Composing Production API =============
docker-compose -f ../docker-compose.yml up -d

echo ============= Displaying Logs =============
docker-compose -f ../docker-compose.yml logs -f