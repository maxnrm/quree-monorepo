FROM nginx:stable-alpine
COPY ./cmd/tg-miniapp/index.html /usr/share/nginx/html/
COPY ./cmd/tg-miniapp/default.conf.template /etc/nginx/templates/
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]