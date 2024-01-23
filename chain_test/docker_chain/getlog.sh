#! /bin/sh
if [ ! -d logs ]; then
    mkdir logs
fi
docker-compose logs --no-color > logs/chain.log
grep "nddn-miner-1" logs/chain.log > logs/nddn-miner-1.log
grep "nddn-miner-2" logs/chain.log > logs/nddn-miner-2.log
grep "nddn-miner-3" logs/chain.log > logs/nddn-miner-3.log
grep "nddn-miner-4" logs/chain.log > logs/nddn-miner-4.log
grep "nddn-miner-5" logs/chain.log > logs/nddn-miner-5.log
grep "nddn-node" logs/chain.log > logs/nddn-node.log