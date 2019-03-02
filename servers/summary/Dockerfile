# use `node` image as the base
FROM alpine
# copy your JavaScript source files
COPY summary /summary
# declare that your server will listen on ports 5001 through 5010
EXPOSE 80
# set the entrypoint command
ENTRYPOINT ["/summary"]
