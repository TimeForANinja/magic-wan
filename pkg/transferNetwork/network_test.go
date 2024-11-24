package transferNetwork

import (
	"net"
	"testing"
)

/* Manually calculated expectation
Below you can find (parts) the calculated `/31` networks:

|        | Target | 0 |   1 |   2 |    3 |    4 |    5 |    6 |    7 |    8 |    9 |    10 |    11 |    12 |    13 |    14 |    15 |    16 |
| ------ | ------ | - | --- | --- | ---- | ---- | ---- | ---- | ---- | ---- | ---- | ----- | ----- | ----- | ----- | ----- | ----- | ----- |
| Source |        |   |     |     |      |      |      |      |      |      |      |       |       |       |       |       |       |       |
|      0 |        | * | 0.0 | 0.2 | 0.06 | 0.12 | 0.20 | 0.30 | 0.42 | 0.56 | 0.72 |  0.90 | 0.110 | 0.132 | 0.156 | 0.182 | 0.210 | 0.240 |
|      1 |        | * |   * | 0.4 | 0.08 | 0.14 | 0.22 | 0.32 | 0.44 | 0.58 | 0.74 |  0.92 | 0.112 | 0.134 | 0.158 | 0.184 | 0.212 | 0.242 |
|      2 |        | * |   * |   * | 0.10 | 0.16 | 0.24 | 0.34 | 0.46 | 0.60 | 0.76 |  0.94 | 0.114 | 0.136 | 0.160 | 0.186 | 0.214 | 0.244 |
|      3 |        | * |   * |   * |    * | 0.18 | 0.26 | 0.36 | 0.48 | 0.62 | 0.78 |  0.96 | 0.116 | 0.138 | 0.162 | 0.188 | 0.216 | 0.246 |
|      4 |        | * |   * |   * |    * |    * | 0.28 | 0.38 | 0.50 | 0.64 | 0.80 |  0.98 | 0.118 | 0.140 | 0.164 | 0.190 | 0.218 | 0.248 |
|      5 |        | * |   * |   * |    * |    * |    * | 0.40 | 0.52 | 0.66 | 0.82 | 0.100 | 0.120 | 0.142 | 0.166 | 0.192 | 0.220 | 0.250 |
|      6 |        | * |   * |   * |    * |    * |    * |    * | 0.54 | 0.68 | 0.84 | 0.102 | 0.122 | 0.144 | 0.168 | 0.194 | 0.222 | 0.252 |
|      7 |        | * |   * |   * |    * |    * |    * |    * |    * | 0.70 | 0.86 | 0.104 | 0.124 | 0.146 | 0.170 | 0.196 | 0.224 | 0.254 |
|      8 |        | * |   * |   * |    * |    * |    * |    * |    * |    * | 0.88 | 0.106 | 0.126 | 0.148 | 0.172 | 0.198 | 0.226 |   1.0 |
|      9 |        | * |   * |   * |    * |    * |    * |    * |    * |    * |    * | 0.108 | 0.128 | 0.150 | 0.174 | 0.200 | 0.228 |   1.2 |
|     10 |        | * |   * |   * |    * |    * |    * |    * |    * |    * |    * |     * | 0.130 | 0.152 | 0.176 | 0.202 | 0.230 |   1.4 |
|     11 |        | * |   * |   * |    * |    * |    * |    * |    * |    * |    * |     * |     * | 0.154 | 0.178 | 0.204 | 0.232 |   1.6 |
|     12 |        | * |   * |   * |    * |    * |    * |    * |    * |    * |    * |     * |     * |     * | 0.180 | 0.206 | 0.234 |   1.8 |
|     13 |        | * |   * |   * |    * |    * |    * |    * |    * |    * |    * |     * |     * |     * |     * | 0.208 | 0.236 |  1.10 |
|     14 |        | * |   * |   * |    * |    * |    * |    * |    * |    * |    * |     * |     * |     * |     * |     * | 0.238 |  1.12 |
|     15 |        | * |   * |   * |    * |    * |    * |    * |    * |    * |    * |     * |     * |     * |     * |     * |     * |  1.14 |

*/

func Test_getPeerToPeerNet(t *testing.T) {
	tests := []struct {
		name           string
		myIDX, peerIDX uint8
		baseNet        string
		wantMyIP       string
		wantPeerIP     string
		wantP2PNet     string
	}{
		{
			name:       "First Address",
			myIDX:      0,
			peerIDX:    1,
			baseNet:    "192.168.1.0/24",
			wantMyIP:   "192.168.1.0",
			wantPeerIP: "192.168.1.1",
			wantP2PNet: "192.168.1.0/31",
		},
		{
			name:       "First Address (inverted)",
			myIDX:      1,
			peerIDX:    0,
			baseNet:    "192.168.1.0/24",
			wantMyIP:   "192.168.1.1",
			wantPeerIP: "192.168.1.0",
			wantP2PNet: "192.168.1.0/31",
		},
		{
			name:       "Medium Address",
			myIDX:      2,
			peerIDX:    4,
			baseNet:    "192.168.1.0/24",
			wantMyIP:   "192.168.1.16",
			wantPeerIP: "192.168.1.17",
			wantP2PNet: "192.168.1.16/31",
		},
		{
			name:       "Medium Address (inverted)",
			myIDX:      4,
			peerIDX:    2,
			baseNet:    "192.168.1.0/24",
			wantMyIP:   "192.168.1.17",
			wantPeerIP: "192.168.1.16",
			wantP2PNet: "192.168.1.16/31",
		},
		{
			name:       "High Address",
			myIDX:      16,
			peerIDX:    13,
			baseNet:    "10.0.0.0/16",
			wantMyIP:   "10.0.1.11",
			wantPeerIP: "10.0.1.10",
			wantP2PNet: "10.0.1.10/31",
		},
		{
			name:       "High Address (inverted)",
			myIDX:      13,
			peerIDX:    16,
			baseNet:    "10.0.0.0/16",
			wantMyIP:   "10.0.1.10",
			wantPeerIP: "10.0.1.11",
			wantP2PNet: "10.0.1.10/31",
		},
		{
			name:       "IPv6 Address",
			myIDX:      0,
			peerIDX:    1,
			baseNet:    "2001:db8::/64",
			wantMyIP:   "2001:db8::",
			wantPeerIP: "2001:db8::1",
			wantP2PNet: "2001:db8::/127",
		},
		{
			name:       "IPv6 Address (inverted)",
			myIDX:      1,
			peerIDX:    0,
			baseNet:    "2001:db8::/64",
			wantMyIP:   "2001:db8::1",
			wantPeerIP: "2001:db8::",
			wantP2PNet: "2001:db8::/127",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, baseNet, err := net.ParseCIDR(tt.baseNet)
			if err != nil {
				t.Fatalf("Invalid Testdata - Unexpected error: %v", err)
			}
			myIP, peerIP, p2pNet, err := GetPeerToPeerNet(tt.myIDX, tt.peerIDX, baseNet)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			} else {
				if myIP.String() != tt.wantMyIP {
					t.Errorf("Expected myIP: %v, got: %v", tt.wantMyIP, myIP)
				}
				if peerIP.String() != tt.wantPeerIP {
					t.Errorf("Expected peerIP: %v, got: %v", tt.wantPeerIP, peerIP)
				}
				if p2pNet.String() != tt.wantP2PNet {
					t.Errorf("Expected p2pNet: %v, got: %v", tt.wantP2PNet, p2pNet)
				}
			}
		})
	}
}
