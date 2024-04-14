import yaml
from calc_p2p_net import calc_p2p_net, calc_p2p_port

def generate_wireguard_configs(config_path):
    with open(config_path, 'r') as file:
        config = yaml.safe_load(file)

    router_net = config['router_net']

    for peer in config['peers']:
        my_idx = config['idx']
        target_idx = peer['idx']
        private_key = config['privatekey']
        public_key = peer['public_key']
        destination = peer['destination']

        wg_conf = f"[Interface]\n"
        wg_conf += f"PrivateKey = {private_key}\n"
        wg_conf += f"Address = {calc_p2p_net(router_net, my_idx, target_idx)}\n"
        wg_conf += f"ListenPort = {calc_p2p_port(my_idx, target_idx)}\n"
        wg_conf += "Table = off\n"
#        wg_conf += "PostUp = sysctl -w net.ipv4.ip_forward=1\n"
#        wg_conf += "PostDown = sysctl -w net.ipv4.ip_forward=0\n"
        wg_conf += f"\n"
        wg_conf += f"[Peer]\n"
        wg_conf += f"PublicKey = {public_key}\n"
        wg_conf += f"AllowedIPs = {calc_p2p_net(router_net, target_idx, my_idx)}, 0.0.0.0/0\n"
        wg_conf += f"Endpoint = {destination}:{calc_p2p_port(my_idx, target_idx)}\n"

        interface_name = f"wg1{my_idx:02d}{target_idx:02d}"
        file_name = f"/etc/wireguard/{interface_name}.conf"
        with open(file_name, 'w') as f:
            f.write(wg_conf)

if __name__ == "__main__":
    generate_wireguard_configs("/magic-wan/config.yml")
