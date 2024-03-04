docker stop test-miniapp
sleep 1
docker run --rm -d --network host --env-file=.env --name test-miniapp test-miniapp
