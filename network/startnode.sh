#!/bin/bash
set -e

command_exists() {
    command -v "$@" > /dev/null 2>&1
}

get_distribution() {
	lsb_dist=""
	# Every system that we officially support has /etc/os-release
	if [ -r /etc/os-release ]; then
		lsb_dist="$(. /etc/os-release && echo "$ID")"
	fi
	# Returning an empty string here should be alright since the
	# case statements don't act unless you provide an actual value
	echo "$lsb_dist"
}

# Set package manager based on OS distribution
lsb_dist=$( get_distribution )
lsb_dist="$(echo "$lsb_dist" | tr '[:upper:]' '[:lower:]')"
case "$lsb_dist" in
    ubuntu|debian|raspbian)
        pkg_manager="apt"
        ;;
    centos)
        pkg_manager="yum"
        ;;
    fedora)
        pkg_manager="dnf"
        ;;
    *)
        echo
        echo "ERROR: Unsupported distribution '$lsb_dist'"
        echo
        exit 1
        ;;
esac

# Install Git
if ! command_exists git; then
    if command_exists $pkg_manager
    then
        y | sudo $pkg_manager install git-all
    else
        echo
        echo "ERROR: Missing expected package manager '$pkg_manager'"
        echo
        exit 1
    fi
fi

# Install Docker
if ! command_exists docker; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
fi

# Download Gecko
gecko_remote="https://github.com/ava-labs/gecko.git"
if [ ! -d gecko ]
then
    git clone $gecko_remote
else
    cd gecko
    git_remote="$(git config --get remote.origin.url)"
    if [ git_remote != gecko_remote ]
    then
        git pull origin master
        cd ..
    else
        echo
        echo "ERROR: Existing directory 'gecko' not using '$gecko_remote' as git remote.origin.url"
        echo
        exit 1
    fi
fi

# Build and run Docker image as daemon
sudo service docker start
sudo usermod -aG docker $USER
sudo su -l $USER <<'END'
    GECKO_COMMIT="$(git --git-dir="./gecko/.git" rev-parse --short HEAD)"
    docker_image="$(docker images -q gecko-$GECKO_COMMIT:latest 2> /dev/null)"
    if [ -z $docker_image ]; then
        ./gecko/scripts/build_image.sh
    fi
    docker run -d -p 9650:9650 -p 9651:9651 gecko-$GECKO_COMMIT /gecko/build/ava \
        --public-ip=127.0.0.1 \
        --snow-sample-size=1 \
        --snow-quorum-size=1 \
        --staking-tls-enabled=false
END
