#!/bin/bash

# prep configs
python build_cfg_wg.py /etc/magicwan/config.yaml
python build_cfg_frr.py /etc/magicwan/config.yaml

# start components
sudo /etc/magicwan/start_wg_services.sh
sudo systemctl start frr
