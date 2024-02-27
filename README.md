# quree-monorepo

# nats port-forward

k port-forward svc/nats 4222:4222 -n nats

# pub-test
nats pub messages.tg123 '{"chat_id": "", "bot_token": "", "text": "hello {{.Count}}"}' --count=10
