FROM alpine:latest
COPY chisel /usr/bin/chisel
COPY ./auth.json /tmp/auth.json
CMD /usr/bin/chisel server --authfile /tmp/auth.json --reverse
