# running a private chain using docker compose

## Prerequisites

Need `docker` (with `docker compose`), and `golang:1.21-bookworm` image for cross build(from non-ubuntu OS), and `ubuntu` image for running the nodes.

```bash
docker pull ubuntu:22.04
docker pull golang:1.21-bookworm
```

> Notice:  
> By default it will pull the image base on your host platform,   
> for example, it will pull a `linux/amd64` image on macOS with `intel` chip,  
> and pull a `linux/arm64` image on macOS with `Apple M1` chip.

Try the `local-cross-build.sh` to build the binary for ubuntu .

## config

Prepare your config file, and the docker compose.yml, which is relate to how much nodes do you want to run and each miner's key and address.

And then config the miner's(validator's) addresses to the `genesis.json`.

> The node keys for used in p2p connection is recorded at `nodekeys.txt`

## run

When all config is fine, then you can run up your private chain:

```bash
docker compose up --build -d
```

## the accounts in the current genesis.json 

|address     |privatekey      |publickey|
|--|--|--|
|0x3a696FeAe901DAe50967F28D7A2225577052F394      |ebaa2febee077847f41b9bd23b28ba7318f37d92658ccbe194a2df432a93810f        |04aa560ea7a3a11bb3831a7f461132e5d8f6928de996784367575e25be66b775e1fcd4dac12e127c8596ea3fbe8bcd6b8ef87800233683c3074b292f68f8cdf763|
|0xa27D573683766F78A818F169C20E287149D26b09      |5e9561af4f2963911d4c04c0fe830666f57b0d87f9bd24ffc4f65aad2a2c2de1|        040fef6aa67ec70a8741d7524db255145fe6dd052e572aa92af5d447a813458ae278073beeb576da21304bfe4a8c84f990c1e84a2b5f39bb12d8f635b96e8642e6|

## some usefull commands

```bash
docker compose up --build -d

docker compose down
```

get logs:

```bash
docker compose logs --no-color > chain.log
```
there's also a shell script to do this:
```bash
./getlog.sh
```

enter the running container:

```bash
docker compose exec -- nddn-node sh
```

directory attach the geth.ipc and run a query:

```bash
docker compose exec -- nddn-miner-1 geth attach --exec eth.blockNumber  /root/data/geth.ipc

docker compose exec -- nddn-miner-1 geth attach --exec 'eth.getBlockByNumber("0x1",false)' /root/data/geth.ipc
```

using json rpc:

```bash
curl --location --request POST 'localhost:8545' --header 'Content-Type: application/json' --data-raw '{"jsonrpc": "2.0","method": "eth_blockNumber","params": [],"id": 10}'
```
