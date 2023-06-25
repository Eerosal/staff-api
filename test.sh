#!/bin/sh
docker compose -f docker-compose.test.yml down && \
docker compose -f docker-compose.test.yml build --no-cache && \
docker compose -f docker-compose.test.yml run --rm --build staff-api-tests