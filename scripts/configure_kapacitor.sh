#! /bin/sh

export TASKS_DIR=/root/openleaf-analytics/docker/kapacitor/tasks

docker run --rm --net=host -it -v $TASKS_DIR:/tasks kapacitor:openleaf kapacitor define batch_health -type batch -tick /tasks/health.tick -dbrp telegraf.autogen
echo --------- uploaded task ---------

docker run --rm --net=host -it kapacitor:openleaf kapacitor enable batch_health
echo --------- enabled task ---------

echo --------- show task ---------
docker run --rm --net=host -it kapacitor:openleaf kapacitor show batch_health

echo --------- list tasks ---------
docker run --rm --net=host -it kapacitor:openleaf kapacitor list tasks