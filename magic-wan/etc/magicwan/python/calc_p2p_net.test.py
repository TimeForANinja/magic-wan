import unittest
from calc_p2p_net import calc_p2p_net, _offset_ipnet, calc_p2p_port

class Test_offset_ipnet(unittest.TestCase):
    def test_valid_cases(self):
        self.assertEqual(_offset_ipnet("192.168.0.0/16", 0), "192.168.0.0")
        self.assertEqual(_offset_ipnet("192.168.0.0/16", 1), "192.168.0.1")

        self.assertEqual(_offset_ipnet("192.168.0.0/16", 254), "192.168.0.254")
        self.assertEqual(_offset_ipnet("192.168.0.0/16", 255), "192.168.0.255")
        self.assertEqual(_offset_ipnet("192.168.0.0/16", 256), "192.168.1.0")
        self.assertEqual(_offset_ipnet("192.168.0.0/16", 257), "192.168.1.1")

        self.assertEqual(_offset_ipnet("192.168.0.0/16", 510), "192.168.1.254")
        self.assertEqual(_offset_ipnet("192.168.0.0/16", 511), "192.168.1.255")
        self.assertEqual(_offset_ipnet("192.168.0.0/16", 512), "192.168.2.0")
        self.assertEqual(_offset_ipnet("192.168.0.0/16", 513), "192.168.2.1")

        self.assertEqual(_offset_ipnet("192.168.0.0/16", 65534), "192.168.255.254")
        self.assertEqual(_offset_ipnet("192.168.0.0/16", 65535), "192.168.255.255")

    def test_invalid_case(self):
        with self.assertRaises(ValueError):
            _offset_ipnet("172.31.0.0/24", -1)
            _offset_ipnet("172.31.0.0/24", 256)
            _offset_ipnet("172.31.0.0/16", 65536)

class Testcalc_p2p_port(unittest.TestCase):
    def test_valid_cases(self):
        self.assertEqual(calc_p2p_port(0, 0), "51820")
        self.assertEqual(calc_p2p_port(0, 1), "51821")
        self.assertEqual(calc_p2p_port(0, 2), "51822")
        self.assertEqual(calc_p2p_port(0, 3), "51823")
        self.assertEqual(calc_p2p_port(0, 4), "51824")

        self.assertEqual(calc_p2p_port(0, 0), "51820")
        self.assertEqual(calc_p2p_port(1, 0), "51821")
        self.assertEqual(calc_p2p_port(2, 0), "51822")
        self.assertEqual(calc_p2p_port(3, 0), "51823")
        self.assertEqual(calc_p2p_port(4, 0), "51824")

        self.assertEqual(calc_p2p_port(1, 1), "51820")
        self.assertEqual(calc_p2p_port(1, 2), "51821")
        self.assertEqual(calc_p2p_port(1, 3), "51822")
        self.assertEqual(calc_p2p_port(1, 4), "51823")

        self.assertEqual(calc_p2p_port(1, 1), "51820")
        self.assertEqual(calc_p2p_port(2, 1), "51821")
        self.assertEqual(calc_p2p_port(3, 1), "51822")
        self.assertEqual(calc_p2p_port(4, 1), "51823")

class Testcalc_p2p_net(unittest.TestCase):
    def test_valid_cases(self):
        # network between 1 and 2
        self.assertEqual(calc_p2p_net("192.168.1.0/24", 1, 2), "192.168.1.0/31")
        self.assertEqual(calc_p2p_net("192.168.1.0/24", 2, 1), "192.168.1.1/31")
        # network between 1 and 3
        self.assertEqual(calc_p2p_net("192.168.1.0/24", 1, 3), "192.168.1.2/31")
        self.assertEqual(calc_p2p_net("192.168.1.0/24", 3, 1), "192.168.1.3/31")
        # network between 2 and 3
        self.assertEqual(calc_p2p_net("192.168.1.0/24", 2, 3), "192.168.1.4/31")
        self.assertEqual(calc_p2p_net("192.168.1.0/24", 3, 2), "192.168.1.5/31")
        # network between 1 and 4
        self.assertEqual(calc_p2p_net("192.168.1.0/24", 1, 4), "192.168.1.6/31")
        self.assertEqual(calc_p2p_net("192.168.1.0/24", 4, 1), "192.168.1.7/31")
        # network between 2 and 4
        self.assertEqual(calc_p2p_net("192.168.1.0/24", 2, 4), "192.168.1.8/31")
        self.assertEqual(calc_p2p_net("192.168.1.0/24", 4, 2), "192.168.1.9/31")
        # network between 3 and 4
        self.assertEqual(calc_p2p_net("192.168.1.0/24", 3, 4), "192.168.1.10/31")
        self.assertEqual(calc_p2p_net("192.168.1.0/24", 4, 3), "192.168.1.11/31")

if __name__ == '__main__':
    unittest.main()
