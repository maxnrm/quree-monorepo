FROM nginx:stable-alpine
COPY ./cmd/tg-miniapp/admin.html /usr/share/nginx/html/webapp/admin/
COPY ./cmd/tg-miniapp/admin_pass_check.html /usr/share/nginx/html/webapp/admin_pass_check/
COPY ./cmd/tg-miniapp/user.html /usr/share/nginx/html/webapp/user/
COPY ./cmd/tg-miniapp/default.conf.template /etc/nginx/templates/
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]