#!/bin/bash

docker build -f build/tg-miniapp/Dockerfile . -t test-miniapp
sleep 0.5
docker stop test-miniapp
sleep 0.5
docker run --rm -d --network host --env-file=.env --name test-miniapp test-miniapp
