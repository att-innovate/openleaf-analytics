#! /bin/sh

export GRAFANA_CONFIG_DIR=/root/openleaf-analytics/grafana

echo --------- configure datasource ---------
curl -X POST --header "Content-Type: application/json" -d @$GRAFANA_CONFIG_DIR/datasource.json http://admin:admin@localhost:8082/api/datasources

echo --------- configure dashboards ---------
curl -X POST --header "Content-Type: application/json" -d @$GRAFANA_CONFIG_DIR/health.json http://admin:admin@localhost:8082/api/dashboards/db
curl -X POST --header "Content-Type: application/json" -d @$GRAFANA_CONFIG_DIR/network.json http://admin:admin@localhost:8082/api/dashboards/db
curl -X POST --header "Content-Type: application/json" -d @$GRAFANA_CONFIG_DIR/docker.json http://admin:admin@localhost:8082/api/dashboards/db
