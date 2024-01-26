# Some integration test of NDDN

Currently this project contains the integration test of NPoS.

For the integration test cases of NPoS and staking, refer to: [npos test](npos_test.md) .

Make sure you have installed the basic dependencies:
- Docker
- Node.js

and already know how to run a local chain by using the docker compose project in `../docker_chain`.

## Preparation

Install js dependencies:
```shell
yarn
```

## Run

Start up the chain and run all test:

```shell
./run-all-test.sh
```

Or if you start the chain mannually, then you can run the test as follows:
```shell
yarn test-staking
yarn test-dao
```