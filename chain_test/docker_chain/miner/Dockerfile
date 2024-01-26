FROM ubuntu:22.04
WORKDIR "/root"
ADD ./set_up_file/geth  /usr/local/bin/geth
ADD ./set_up_file/genesis.json  ./genesis.json
ADD ./set_up_file/config.toml  ./config.toml

RUN mkdir -p ./keystore \
&& mkdir -p ./data \
&& chmod 755 /usr/local/bin/geth


ARG password
ARG privatekey
RUN echo $password > ./.accountpassword
RUN echo $privatekey > ./.privatekey


RUN geth init --datadir /root/data ./genesis.json
RUN geth account import --datadir /root/data --password ./.accountpassword  ./.privatekey

ENV address=""
ENV nodekeyhex=""

ENTRYPOINT ["sh","-c","geth --config config.toml --nodekeyhex $nodekeyhex  --mine --miner.etherbase $address --unlock $address --password ./.accountpassword --allow-insecure-unlock --internaltx.disabled"]

EXPOSE 8545 8546 30303 30303/udp
