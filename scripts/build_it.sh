#! /bin/sh

if [[ $UID != "0" ]]
then
    echo build_it script needs to be run as root user
    exit 1
fi

export ROOT=/root
export ANALYTICS_DIR=$ROOT/openleaf-analytics
export DOCKER_DIR=$ANALYTICS_DIR/docker
export INFLUXDB_DIR=$ROOT/openleaf-data/influxdb
export KAPACITOR_DIR=$ROOT/openleaf-data/kapacitor
export GRAFANA_DIR=$ROOT/openleaf-data/grafana

echo -- Directories used by build_it --
echo $ANALYTICS_DIR
echo $DOCKER_DIR

echo -- Directories created for persistent data storage --
echo $INFLUXDB_DIR
echo $KAPACITOR_DIR
echo $GRAFANA_DIR

echo -- Building ... --

# Create folder for InfluxDB data
mkdir -p $INFLUXDB_DIR

# Create folder for Kapacitor data
mkdir -p $KAPACITOR_DIR

# Create folder for Grafana data
mkdir -p $GRAFANA_DIR

# Build Docker Containers
cd $DOCKER_DIR/grafana
docker build -t grafana:openleaf .

cd $DOCKER_DIR/influxdb
docker build -t influxdb:openleaf .

cd $DOCKER_DIR/kapacitor
docker build -t kapacitor:openleaf .

cd $DOCKER_DIR/telegraf
if [ ! -d "src/telegraf" ]:
then
    cd src
    git clone https://github.com/influxdata/telegraf.git
    cd telegraf
    git checkout release-1.3
    rm -rf .git
    cd ..
    mkdir -p telegraf/plugins/inputs/snaproute
    cd ..
fi
cp $ANALYTICS_DIR/telegraf-snaproute/all/all.go src/telegraf/plugins/inputs/all/
cp $ANALYTICS_DIR/telegraf-snaproute/snaproute/snaproute.go src/telegraf/plugins/inputs/snaproute/
docker build -t telegraf:openleaf .
