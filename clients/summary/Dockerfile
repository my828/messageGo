# add the necessary instructions
# to create a Docker container image
# for your summary client using the
# NGINX web server image as the base
FROM nginx
COPY . /usr/share/ngix/html
ADD default.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
