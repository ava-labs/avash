#!/bin/bash
set -e

# Node name
NAME=""

# Required flag defaults
H_PORT="9650"
S_PORT="9651"
PUBLIC_IP="127.0.0.1"

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
		--xput-server-port=*|\
		--signature-verification-enabled=*|\
		--api-ipcs-enabled=*|\
		--http-tls-enabled=*|\
		--http-tls-cert-file=*|\
		--http-tls-key-file=*|\
		--bootstrap-ips=*|\
		--bootstrap-ids=*|\
		--db-enabled=*|\
		--db-dir=*|\
		--log-level=*|\
		--log-dir=*|\
		--snow-avalanche-batch-size=*|\
		--snow-avalanche-num-parents=*|\
		--snow-sample-size=*|\
		--snow-quorum-size=*|\
		--snow-virtuous-commit-threshold=*|\
		--snow-rogue-commit-threshold=*|\
		--staking-tls-enabled=*|\
		--staking-tls-key-file=*|\
		--staking-tls-cert-file=*)
            FLAGS+="${arg} "
            ;;
        --http-port=*)
            H_PORT="${arg#*=}"
            ;;
        --staking-port=*)
            S_PORT="${arg#*=}"
            ;;
        --public-ip=*)
            PUBLIC_IP="${arg#*=}"
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
FLAGS+="--public-ip=${PUBLIC_IP} "

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
docker run -d --name $NAME -p $H_PORT:$H_PORT -p $S_PORT:$S_PORT gecko-$GECKO_COMMIT /gecko/build/ava $FLAGS
