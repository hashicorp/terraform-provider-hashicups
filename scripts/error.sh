#!/bin/bash

BASEDIR=$(dirname "$0")
echo ${BASEDIR}

DOCKERCOMPOSEFILE="${BASEDIR}/../docker_compose/docker-compose.yml"

echo ${DOCKERCOMPOSEFILE}

function finish {
    docker-compose -f ${DOCKERCOMPOSEFILE} down
}

trap finish EXIT

docker-compose -f ${DOCKERCOMPOSEFILE} up &

CNT=0
until curl -X POST localhost:19090/signup -d '{"username":"education", "password":"test123"}' || [ ${CNT} -gt 10 ]
do
  sleep 1
 ((CNT=CNT+1))
done

cd ${BASEDIR}/.. 

TF_ACC=1 go test -count=1 -v ./hashicups
