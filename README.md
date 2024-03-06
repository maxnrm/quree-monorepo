# postgres
create index messages_type on messages(type);
create index user_event_visits_user_id on user_event_visits(user_chat_id);
alter table user_event_visits add constraint fkey_users_users foreign key (user_chat_id) references users (chat_id);

psql "postgres://quree:qureequree@127.0.0.1:5432/quree

# helm
kubectl create secret generic quree-env-secrets --from-env-file=.env.prod -n quree

# nats port-forward
k port-forward svc/nats 4222:4222 -n nats

# pub-test
nats pub messages.tg123 '{"chat_id": "", "bot_token": "", "text": "hello {{.Count}}"}' --count=10

# ngrok

ngrok http --domain=cunning-communal-mudfish.ngrok-free.app 80
