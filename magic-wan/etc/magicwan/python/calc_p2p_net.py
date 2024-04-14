import ipaddress

def calc_p2p_net(router_net, idx, peer_idx):
    if idx > peer_idx:
        host_offset = _calc_host_idx(idx, peer_idx) + 1
    else:
        host_offset = _calc_host_idx(peer_idx, idx)

    return _offset_ipnet(router_net, host_offset) + "/31"

def _offset_ipnet(subnet, host):
    # Parse the subnet
    network = ipaddress.ip_network(subnet)

    # already excludes the network and broadcast address
    if host > network.num_addresses:
        raise ValueError("Error: host exceeds subnet range")
    if host < 0:
        raise ValueError("Error: host exceeds subnet range")

    # Get the nth host within the subnet
    host = network.network_address + host
    return str(host)

def _calc_host_idx(row, col):
    # all params -1 to offset from 0 to 1 indexing
    return (row-2) * (row-1) + 2*(col-1)

def calc_p2p_port(idx, peer_idx):
    return str(51820 + abs(idx - peer_idx))
