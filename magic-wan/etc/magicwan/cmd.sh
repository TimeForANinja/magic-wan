#!/bin/bash

# Function to start the service
start_service() {
    sudo systemctl start magicwan
}

# Function to auto-start the service
auto_start_service() {
    sudo systemctl enable magicwan
}

# Function to stop the service
stop_service() {
    sudo systemctl stop magicwan
}

# Function to view the configuration
view_config() {
    cat /etc/magic_wan/config.yaml
}

# Function to load a new configuration
load_config() {
    sudo cp "$1" /etc/magic_wan/config.yaml
    sudo systemctl restart magicwan
}

# Main function to handle subcommands
main() {
    case "$1" in
        start)
            start_service
            ;;
        auto_start)
            auto_start_service
            ;;
        stop)
            stop_service
            ;;
        view)
            view_config
            ;;
        load)
            if [ -z "$2" ]; then
                echo "Error: Please provide a path to the configuration file."
                exit 1
            fi
            load_config "$2"
            ;;
        *)
            echo "Error: Invalid subcommand. Usage: magicwan [start|auto_start|stop|view|load /path/to/config.yaml]"
            exit 1
            ;;
    esac
}

# Call the main function with the provided arguments
main "$@"
