#!/bin/bash

CONFIG_FILE="/magic-wan/config.yml"

# Paths to Python scripts
WIREGUARD_SCRIPT="/magic-wan/build_wireguard_config.py"
FRR_SCRIPT="/magic-wan/build_frr_config.py"

# Check if the config file exists
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Error: Config file '$CONFIG_FILE' not found. Exiting."
    exit 1
fi

# Call Python scripts to build configurations
echo "Building WireGuard configuration..."
python3 "$WIREGUARD_SCRIPT" "$CONFIG_FILE"
echo "WireGuard configuration built."

echo "Building FRR configuration..."
python3 "$FRR_SCRIPT" "$CONFIG_FILE"
echo "FRR configuration built."

# Start WireGuard service
echo "Starting WireGuard service..."
# Example command to start WireGuard service
for wg_cfg in /etc/wireguard/wg*.conf; do
    wg-quick up "$wg_cfg"
done
echo "WireGuard service started."

# Start FRR service
echo "Starting FRR service..."
# Example command to start FRR service
/usr/lib/frr/frrinit.sh start
echo "FRR service started."

# Keep the container running & print debug log every few seconds
while true; do
    /magic-wan/debug.sh
    sleep 10
done
