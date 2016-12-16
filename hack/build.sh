#!/usr/bin/env bash
set -ex

# Do things differently when already in project_root/hack
if [ $(pwd | tr "/" "\n" | tail -1) == "hack" ]; then
echo "please execute in the project root"
exit 1;
fi

WORK_DIR=$(pwd)
GIT_BRANCH=${GIT_BRANCH:="any"}
REPO=074509403805.dkr.ecr.eu-west-1.amazonaws.com
IMAGE=aws-nuke

if [ $GIT_BRANCH = "origin/master" ]
then
    docker build -f hack/Dockerfile -t ${REPO}/${IMAGE}:latest $WORK_DIR
    docker push ${REPO}/${IMAGE}:latest
else
    docker build -f hack/Dockerfile -t ${REPO}/${IMAGE}:pr_build $WORK_DIR
fi
