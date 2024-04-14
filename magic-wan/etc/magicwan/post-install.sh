#!/etc/bash

# script to start service
chmod +x /etc/magicwan/start.sh

# cli cmd
chmod +x /etc/magicwan/cmd.sh
sudo ln -s /etc/magicwan/cmd.sh /usr/local/bin/magicwan

# create config storage
touch /etc/magicwan/config.yaml

# magic-wan service
sudo ln -s /etc/magicwan/service/magicwan.service /etc/systemd/system/magicwan.service
systemctl daemon-reload
