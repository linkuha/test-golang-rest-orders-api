FROM nginx:1.23-alpine-perl

COPY ./docker/production/nginx-gateway/nginx.conf /etc/nginx/nginx.conf

COPY ./docker/production/nginx-gateway/sites-enabled /etc/nginx/sites-enabled/
COPY ./docker/production/nginx-gateway/snippets /etc/nginx/snippets/
COPY ./docker/production/nginx-gateway/templates /etc/nginx/templates/
