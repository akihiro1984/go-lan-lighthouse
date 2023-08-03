#!/bin/bash


RES=$(sudo docker compose run --rm main)

TARGET=$(echo $RES | cut -d' ' -f2)

sed -e "s/^GATEWAY_RESOLVER_HOST=.*$/GATEWAY_RESOLVER_HOST=${TARGET}/g" .env.template > .env
