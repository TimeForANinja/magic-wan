import yaml
from calc_p2p_net import calc_p2p_net

def generate_frr_config(config_path, frr_config_path):
    with open(config_path, 'r') as file:
        config = yaml.safe_load(file)

    my_idx = config['idx']
    router_net = config['router_net']

    # Create FRR configuration
    frr_conf = "log syslog informational\n"
    frr_conf += "!\nrouter ospf\n"
    frr_conf += f" ospf router-id 0.0.0.{my_idx}\n"

    for peer in config['peers']:
        target_idx = peer['idx']
        p2p_net = calc_p2p_net(router_net, my_idx, target_idx)
        frr_conf += f" network {p2p_net} area 0.0.0.0\n"

    frr_conf += "exit\n!\n"

    for peer in config['peers']:
        target_idx = peer['idx']
        interface_name = f"wg1{my_idx:02d}{target_idx:02d}"
        frr_conf += f"interface {interface_name}\n"
        frr_conf += " ip router ospf area 0.0.0.0\n"
        frr_conf += " ip ospf network point-to-point\n"
        frr_conf += "exit\n"
        frr_conf += "!\n"

    with open(frr_config_path, 'w') as f:
        f.write(frr_conf)

def update_ospf_daemons_config(file_path):
    with open(file_path, 'r') as file:
        file_content = file.read()

    # Replace "ospfd=no" with "ospfd=yes"
    updated_content = file_content.replace("ospfd=no", "ospfd=yes")

    with open(file_path, 'w') as file:
        file.write(updated_content)

if __name__ == "__main__":
    generate_frr_config("/magic-wan/config.yml", "/etc/frr/frr.conf")
    update_ospf_daemons_config("/etc/frr/daemons")
