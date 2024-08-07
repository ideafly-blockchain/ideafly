#! /bin/bash
cd $(dirname $0)
cd ../docker_chain
docker compose down --remove-orphans
# docker compose up -d
docker compose up -d --build

cd -
echo "==Start staking test=="
yarn test-staking
a=$?

if [ "$a" -ne 0 ]; then
  echo "==**== node logs ==**=="
  cd ../docker_chain
  docker compose logs --no-color -n 500
  cd -
  echo "==**== end of node logs ==**=="
fi

echo "==Start dao govenence test=="
yarn test-dao
b=$?

if [ "$b" -ne 0 ]; then
  echo "==**== node logs ==**=="
  cd ../docker_chain
  docker compose logs --no-color -n 500
  cd -
  echo "==**== end of node logs ==**=="
fi

# echo "All is done."
# read -r -p "Do you want to stop the chain [y/n]: " input
# case $input in
    # [yY][eE][sS]|[yY])
        cd ../docker_chain
        docker compose down --remove-orphans
        cd -
        # ;;
    # *)
        # echo "you may need to stop the chain manually."
        # ;;
# esac
exit $((a|b))