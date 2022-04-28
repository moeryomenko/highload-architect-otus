FROM haproxytech/haproxy-ubuntu:latest

RUN apt update && apt install mysql-client -y

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["haproxy", "-f", "/usr/local/etc/haproxy/haproxy.cfg"]
