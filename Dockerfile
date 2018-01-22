FROM scratch
ADD test-server /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/test-server"]
