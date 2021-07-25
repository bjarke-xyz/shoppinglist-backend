#!/usr/bin/env bash
docker-compose build && docker-compose push
docker-compose run migrate
docker stack deploy --compose-file=docker-compose.yml --with-registry-auth shoppinglistv4-backend
