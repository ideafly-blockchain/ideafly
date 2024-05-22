#! /bin/bash

cd ../docker_chain
docker compose down --remove-orphans
# docker compose up -d
docker compose up -d --build

cd -
yarn test-jail

echo "All is done."
read -r -p "Do you want to stop the chain [y/n]: " input
case $input in
    [yY][eE][sS]|[yY])
        cd ../docker_chain
        docker compose down --remove-orphans
        cd -
        ;;
    *)
        echo "you may need to stop the chain manually."
        ;;
esac
