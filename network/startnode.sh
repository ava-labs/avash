#!/bin/bash
set -e

DATA_DIR="$(pwd)"
CTNR_DIR="/data"

# Node name
NAME=""

# Required flag defaults
H_PORT="9650"
S_PORT="9651"

# Process flags
FLAGS=""
for arg in "$@"
do
    case "$arg" in
        --name=*)
            NAME="${arg#*=}"
            ;;
        --assertions-enabled=*|\
		--ava-tx-fee=*|\
		--network-id=*|\
        --public-ip=*|\
		--xput-server-port=*|\
        --xput-server-enabled=*|\
		--signature-verification-enabled=*|\
        --api-admin-enabled=*|\
		--api-ipcs-enabled=*|\
        --api-keystore-enabled=*|\
        --api-metrics-enabled=*|\
		--http-tls-enabled=*|\
		--http-tls-cert-file=*|\
		--http-tls-key-file=*|\
		--bootstrap-ips=*|\
		--bootstrap-ids=*|\
		--db-enabled=*|\
		--log-level=*|\
		--snow-avalanche-batch-size=*|\
		--snow-avalanche-num-parents=*|\
		--snow-sample-size=*|\
		--snow-quorum-size=*|\
		--snow-virtuous-commit-threshold=*|\
		--snow-rogue-commit-threshold=*|\
        --p2p-tls-enabled=*|\
		--staking-tls-enabled=*|\
		--staking-tls-key-file=*|\
		--staking-tls-cert-file=*)
            FLAGS+="${arg} "
            ;;
        --db-dir=*|\
        --log-dir=*|\
        --plugin-dir=*)
            FLAGS+="${arg%=*}=$CTNR_DIR/${arg#*=} "
            ;;
        --data-dir=*)
            DATA_DIR+="/${arg#*=}"
            ;;
        --http-port=*)
            H_PORT="${arg#*=}"
            ;;
        --staking-port=*)
            S_PORT="${arg#*=}"
            ;;
        *)
            echo
            echo "ERROR: Unsupported flag '${arg%%=*}'"
            echo
            exit 1
            ;;
    esac
done
FLAGS+="--http-port=${H_PORT} "
FLAGS+="--staking-port=${S_PORT} "

if [ -z $NAME ]; then
    echo
    echo "ERROR: Missing node name argument (--name=[node name] required)"
    echo
    exit 1
fi

# Build and run Docker image as daemon
GECKO_COMMIT="$(git --git-dir="./gecko/.git" rev-parse --short HEAD)"
docker_image="$(docker images -q gecko-$GECKO_COMMIT:latest 2> /dev/null)"
if [ -z $docker_image ]; then
    ./gecko/scripts/build_image.sh
fi
mkdir -p $DATA_DIR
docker run -d --name $NAME \
    -v $DATA_DIR:$CTNR_DIR \
    -p $H_PORT:$H_PORT \
    -p $S_PORT:$S_PORT \
    gecko-$GECKO_COMMIT \
    /gecko/build/ava $FLAGS
