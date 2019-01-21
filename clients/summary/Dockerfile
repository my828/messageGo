# add the necessary instructions
# to create a Docker container image
# for your summary client using the
# NGINX web server image as the base
FROM ngix
COPY summary /usr/share/ngix/html
ADD default.conf /etc/nginx/conf.d/default.conf
ENTRYPOINT [ "/usr/share/ngix/html" ]
