package code

import (
    "net"
    "os"
)

func CompareSelfIP(ip string) bool {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        os.Stderr.WriteString("Oops: " + err.Error() + "\n")
        os.Exit(1)
    }

    for _, a := range addrs {
        ipnet, ok := a.(*net.IPNet);
        if !ok || ipnet.IP.IsLoopback() {
            continue
        }
        if ipnet.IP.To4() == nil {
            continue
        }
        if ipnet.IP.String() == ip {
            return true
        }
    }

    return false
}

func GetLocalIp() string {
    var localIp string
    addrs, err := net.InterfaceAddrs()
    lanCheck := net.IPNet{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(24, 32)}
    if err == nil {
        for _, a := range addrs {
            ipnet, ok := a.(*net.IPNet)
            if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil && lanCheck.Contains(ipnet.IP) {
                localIp = ipnet.IP.String();
                break;
            }
        }
    }
    return localIp
}
