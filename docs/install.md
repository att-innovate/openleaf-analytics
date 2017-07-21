## Installation
### Docker Engine
#### For the 3.18 kernel version of ONL:

	root@wedge:~# apt-get update
	root@wedge:~# apt-get install docker-engine

#### For the 3.16 Jessie kernel version of ONL:

Follow the steps as outlined on  Docker’s [installation guide for Debian](https://docs.docker.com/engine/installation/linux/debian/).

Additionally the storage-driver needs to be changed:

	root@barefoot:~# echo 'DOCKER_OPTS="--storage-driver=overlay"' >> /etc/default/docker
	root@barefoot:~# service docker start

### Docker Compose

	root@barefoot:~# curl -L "https://github.com/docker/compose/releases/download/1.11.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
	root@barefoot:~# chmod +x /usr/local/bin/docker-compose

### Analytics Stack

The scripts used for the building of the components require the stack to be installed under `/root/openleaf-analytics`

#### Clone Project

	root@barefoot:~# cd /root
	root@barefoot:~# git clone https://github.com/att-innovate/openleaf-analytics.git

#### Build Docker Containers

To build the components the switch has to have access to the internet to be able to pull the base images for example. The containers for the components itself get build locally.

	root@barefoot:~# cd openleaf-analytics/
	root@barefoot:~/openleaf-analytics# ./scripts/build_it.sh

#### Run Analytics Stack

	root@barefoot:~/openleaf-analytics# ./scripts/run_it.sh

Verify that containers are up and running:

	root@barefoot:~/openleaf-analytics# docker ps

List should show all the 4 services: grafana, kapacitor, telegraf, influxdb

#### Configure Grafana and Kapacitor

    root@barefoot:~/openleaf-analytics# ./scripts/configure_grafana.sh
    root@barefoot:~/openleaf-analytics# ./scripts/configure_kapacitor.sh 

#### Use Analytics Stack

That's it for the installation. Please check the [User Guide](userguide.md) for additional information.

#### Stop Analytics Stack

	root@barefoot:~/openleaf-analytics# ./scripts/stop_it.sh
