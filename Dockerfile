FROM   scratch

ADD  zoneinfo.tar.gz /
COPY html /html
COPY raad071cal  /raad071cal

EXPOSE 80
WORKDIR /
CMD ["/raad071cal"]
