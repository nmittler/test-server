FROM ubuntu:xenial
# Do not add more stuff to this list that isn't small or critically useful.
# If you occasionally need something on the container do
# sudo /usr/local/bin/proxy-redirection-clear
# sudo apt-get update && apt-get whichever
# sudo /usr/local/bin/proxy-redirection-restore
RUN apt-get update && \
    apt-get install --no-install-recommends -y \
      curl \
      iptables \
      iproute2 \
      iputils-ping \
      dnsutils \
      netcat \
      tcpdump \
      net-tools \
      libc6-dbg gdb \
      elvis-tiny \
      lsof \
      busybox \
      sudo && \
    rm -rf /var/lib/apt/lists/*


ADD test-server /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/test-server"]
