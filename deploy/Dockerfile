FROM debian:jessie

ADD bin/healthcheck /healthcheck

RUN apt-get -y update && \
    apt-get -y --no-install-recommends install ca-certificates && \
    rm -rf /var/lib/apt/lists/*

EXPOSE 80 443

ENTRYPOINT ["/healthcheck"]
CMD ["daemon"]
