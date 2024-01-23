#!/bin/bash

set -e

sysName="$(uname -s)"
arch="$(uname -m)"
case "${sysName}" in
    Linux*)     target=geth;;
    Darwin*)
      case "${arch}" in
        x86_64*)  target=geth-linux-amd64;;
        *)  target=geth-linux-arm64;;
      esac
      ;;
    *)
      echo "system unsupported"
      exit 1
esac
echo "build target: ${target}"

build_path=../../../nddn

cd ${build_path}

make "${target}"

cd -
# check output
if [ -f "${build_path}/build/bin/${target}" ]; then
  echo "build success"

  cp ${build_path}/build/bin/${target} ./set_up_file/geth
  ls -l ./set_up_file
else
  echo "build failed, ${build_path}/build/bin/${target} does not exist"
  exit 1
fi
