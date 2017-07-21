#! /bin/sh
echo --------- grafana ---------
docker logs docker_grafana_1

echo --------- kapacitor ---------
docker logs docker_kapacitor_1

echo --------- telegraf ---------
docker logs docker_telegraf_1
