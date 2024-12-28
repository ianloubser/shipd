#!/bin/sh

# how to run this script ? 
# curl -o- https://raw.githubusercontent.com/ianloubser/shipd/vm-bootstrap.sh {wildcard-domain-here} | bash

set -e;

# install nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash

# install latest node
nvm install node 

# install pm2
npm i -g pm2

# pull deploy script
curl -o- https://raw.githubusercontent.com/ianloubser/shipd/latest/install.sh | bash

# run deploy script with pm2
pm2 start /usr/local/bin/shipd --name ci-server

# install caddyserver
if [ $(uname) == 'Darwin ']
    wget https://github.com/caddyserver/caddy/releases/download/v2.9.0-beta.2/caddy_2.9.0-beta.2_mac_amd64.tar.gz
    tar -zxf caddy.tar.gz
    mv caddy /usr/local/bin/
fi

# setup caddyserver
WILDCARD_DOMAIN=$1

# dump new caddyfile
DEPLOY_SERVER_KEY=$(openssl rand -base64 48)
CADDY_CONTENTS=$(cat <<EOF
ci.${WILDCARD_DOMAIN} {
    @apirequests {
        header X-Auth-token ${DEPLOY_SERVER_KEY}
    }
    route {
        reverse_proxy @apirequests localhost:5000
        respond "You don't look like a teapot ?" 418
    }
}
EOF
)

echo $CADDY_CONTENTS > ~/Caddyfile

# run the caddyserver
pm2 start caddy ~/Caddyfile --name proxy

echo "Done!!"