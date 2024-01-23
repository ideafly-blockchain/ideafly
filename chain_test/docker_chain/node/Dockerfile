FROM ubuntu:22.04
WORKDIR "/root"
ADD ./set_up_file/geth  /usr/local/bin/geth
ADD ./set_up_file/genesis.json  ./genesis.json
ADD ./set_up_file/config.toml  ./config.toml

RUN mkdir -p ./keystore \
&& mkdir -p ./data \
&& chmod 755 /usr/local/bin/geth

RUN geth init --datadir /root/data ./genesis.json


ENTRYPOINT ["geth","--config","config.toml","--allow-insecure-unlock"]
CMD ["--metrics","--metrics.addr", "0.0.0.0","--metrics.port","19090","--internaltx.all"]

EXPOSE 8545 8546 30303 30303/udp 19090
