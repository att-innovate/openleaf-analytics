FROM debian:jessie

ENV DOWNLOAD_URL=https://grafanarel.s3.amazonaws.com/builds/grafana_4.1.2-1486989747_amd64.deb

RUN apt-get update && \
    apt-get -y --no-install-recommends install libfontconfig curl ca-certificates && \
    apt-get clean && \
    curl ${DOWNLOAD_URL} > /tmp/grafana.deb && \
    dpkg -i /tmp/grafana.deb && \
    rm /tmp/grafana.deb && \
    curl -L https://github.com/tianon/gosu/releases/download/1.7/gosu-amd64 > /usr/sbin/gosu && \
    chmod +x /usr/sbin/gosu && \
    apt-get remove -y curl && \
    apt-get autoremove -y && \
    rm -rf /var/lib/apt/lists/*

VOLUME ["/var/lib/grafana", "/var/log/grafana"]

COPY ./grafana.ini /usr/share/grafana/conf/defaults.ini
COPY ./run.sh /run.sh

ENTRYPOINT ["/run.sh"]
