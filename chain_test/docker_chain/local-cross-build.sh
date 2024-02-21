#!/bin/bash
set -e

# build
build_path=../../
cd ${build_path}
docker build -f Dockerfile.debian -t nddn-debian-client:build  .
cd -

# copy binary out of the local image
docker run -d --name nddn-debian-for-copy nddn-debian-client:build /bin/bash
docker cp nddn-debian-for-copy:/go-ethereum/build/bin/geth ${PWD}/set_up_file/geth
docker rm nddn-debian-for-copy
ls -l ./set_up_file