# postgres
ALTER TABLE users ADD UNIQUE (chat_id, role);
ALTER TABLE user_event_visits ADD UNIQUE (user_id, quiz_id);

# helm
kubectl create secret generic quree-env-secrets --from-env-file=.env.prod

# nats port-forward
k port-forward svc/nats 4222:4222 -n nats

# pub-test
nats pub messages.tg123 '{"chat_id": "", "bot_token": "", "text": "hello {{.Count}}"}' --count=10
