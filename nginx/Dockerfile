FROM nginx

COPY ./nginx.conf /etc/nginx/nginx.conf

# NOTE: Delete default.conf to not override...
RUN rm /etc/nginx/conf.d/default.conf

COPY ./server.conf /etc/nginx/conf.d/server.conf
COPY ./certs/cert.crt /etc/ssl/certs/cert.crt
COPY ./keys/key.pem /etc/ssl/private/key.pem
RUN chmod 600 /etc/ssl/private/key.pem
