## User Guide

### Grafana Dashboards
The project comes with three pre-configured Grafana dashboards for System-Health, Docker, and Network metrics data.

Dashboards are listed under the “Home” pulldown in the header of the main page.

Grafana is listening on port 8082 on the main interface of the embedded micro-server: http://ip-address:8082/

Username and password are *admin/admin*.

The dashboards show a subset of all the collected telemetry data. To show additional data or change the dashboards check out to [Grafana documentation](http://docs.grafana.org/guides/getting_started/).

#### System-Health Dashboard
![health dashboard](https://github.com/att-innovate/openleaf-analytics/blob/master/docs/grafana-health.png)

Simple visualization of system data from the embedded micro-server. It shows percentage of CPU, Memory, and Disk space in use.

The “Health Index” chart shows an overall index score of the health of the system (10 == healthiest state). The index gets calculated every 5 minutes by Kapacitor and sent to InfluxDB. For this demo setup Kapacitor checks the Disk space only. More details can be found in the section about [Kapacitor](#kapacitor).

#### Network Dashboard
![network dashboard](https://github.com/att-innovate/openleaf-analytics/blob/master/docs/grafana-network.png)

Visualizes telemetry data from the switching silicon. “Port Traffic” shows “packets per seconds sent” for the individual active ports.

“Port Data In” and “Port Data Out” show “bytes per seconds sent and received”.

“Multicast/Broadcast” shows “max packets per seconds sent” of multicast and broadcast packets overall.

“Discarded Packages” shows “count of packages” discarded since last measurement for ingress and egress for the individual active ports.

#### Docker Dashboard
![docker dashboard](https://github.com/att-innovate/openleaf-analytics/blob/master/docs/grafana-docker.png)

Dashboard shows Memory and CPU data for each individual Docker container.
  
### Telegraf
The configuration file for Telegraf can be found at [/docker/telegraf/telegraf.conf](../docker/telegraf/telegraf.conf).

By default we only use the docker, snaproute, and system input plugins. The plugins are listed in [/telegraf-snaproute/all/all.go](../telegraf-snaproute/all/all.go).

Changes of any of those configurations requires a re-build of the analytics stack, see [Installation Guide](install.md).

Detailed information about Telegraf can be found [online](https://docs.influxdata.com/telegraf/v1.3/).

### InfluxDB and Query Interface
The InfluxDB configuration file can be found at [/docker/influxdb/influxdb.conf](../docker/influxdb/influxdb.conf).

The database with all the other persistent data from the different components gets stored under `/root/openleaf-data`.

For the demo setup we don’t have a retention policy in place. Metrics data will accumulate slowly. Further information about InfluxDB and the configuration options can be found [online](https://docs.influxdata.com/influxdb/v1.3/).

There are two query interfaces available for InfluxDB, a web-based UI and a command-line tool.

#### Web-based Admin UI
![admin ui](https://github.com/att-innovate/openleaf-analytics/blob/master/docs/influxdb-web.png)

The web-based admin and query UI is listening on port 8083 on the main interface of the embedded micro-server: http://ip-address:8083/

To work with the time-series data select “Database: telegraf” from the database pulldown in the right part of the header.

#### Command-line Query Tool
![query tool](https://github.com/att-innovate/openleaf-analytics/blob/master/docs/influxdb-cmdline.png)

On the switch the tool can be started from the shell by calling `.scripts/run_client_influxdb.sh`.

Detailed information about InfluxDB and its query language can be found [online](https://docs.influxdata.com/influxdb/v1.3/).

### Kapacitor
For the demo scenario we use Kapacitor to calculate a simple “Health Index” every 5 minute. The result gets inserted in to InfluxDB. The script for the demo can be found at [/docker/kapacitor/tasks/health.tick](../docker/kapacitor/tasks/health.tick).

Kapacitor itself offers a REST API to upload, update, delete tasks. Detailed information can be found [online](https://docs.influxdata.com/kapacitor/v1.3/api/api/).

### Tools
The [/scripts](../scripts) folder contains some additional operational tools. It is assumed that all those scripts get started from the root of the project directory, for example `./scripts/build_it.sh`.

- `build_it.sh`: Builds all the components needed for the Analytics Stack. It also creates directories need to store the persistent state for all the components.
- `clean_it.sh`: Can be used to clean out old container snapshots created by Docker Compose.
- `configure-grafana.sh`: Uploads our default set of Dashboards to Grafana.
- `configure-kapacitor.sh`: Uploads and activates health script for Kapacitor.
- `log_it.sh`: Prints out log statements for all the components. Some error messages are expected, for example an error about InfluxDB not being accessible. Those errors can occur at startup because all containers get started at the same time and InfluxDB will take a bit to become ready.
- `run_client_influxdb.sh`: Access to the InfluxDB query tool from.
- `run_it.sh`: Starts all the Analytics Stack components.
- `stop_it.sh`: Stops all the Analytics Stack components.
