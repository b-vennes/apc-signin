#!/usr/bin/env bash

git clone https://github.com/b-vennes/apc-signin.git

cd ./apc-signin || exit

docker compose -f compose.live.yaml up -d
