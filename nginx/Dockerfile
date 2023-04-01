FROM nginx

COPY ./nginx.conf /etc/nginx/nginx.conf

# NOTE: Delete default.conf to not override...
RUN rm /etc/nginx/conf.d/default.conf

COPY ./server.conf /etc/nginx/conf.d/server.conf
