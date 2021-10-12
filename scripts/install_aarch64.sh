#!/bin/bash

if ! command -v jq& /dev/null
then
    echo "jq could not be found"
    exit
fi

owner=eldius
repo=mqtt-listener-go

wget "$( curl https://api.github.com/repos/${owner}/${repo}/releases | jq -r '. | sort_by(.created_at) | last | .assets[] | select(.name | endswith(".aarch64")) | .browser_download_url' )"

mv mqtt-listener* mqtt-listener
sudo mv mqtt-listener /usr/bin
