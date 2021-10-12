#!/bin/bash

if ! command -v jq& /dev/null
then
    echo "jq could not be found"
    exit
fi

owner=eldius
repo=mqtt-listener-go

wget "$( curl https://api.github.com/repos/${owner}/${repo}/releases | jq -r '. | sort_by(.created_at) | last | .assets[] | select(.name | endswith(".raspberry")) | .browser_download_url' )"

mv speedtest-wrapper* speedtest-wrapper
sudo mv speedtest-wrapper /usr/bin