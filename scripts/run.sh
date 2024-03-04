docker stop test-miniapp
sleep 1
docker run --rm -d -p 127.0.0.1:7777:80 --env-file=.env --name test-miniapp test-miniapp
