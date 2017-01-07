package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	v6 := flag.Bool("v6", false, "output ipv6 address")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println(GetOutboundIP())
		os.Exit(0)
	}
	for _, arg := range flag.Args() {
		iface, err := net.InterfaceByName(arg)
		if err != nil {
			panic(err)
		}
		ips, err := extractIpsFromInterface(iface, *v6)
		if err != nil {
			panic(err)
		}
		for _, ip := range ips {
			fmt.Println(ip.String())
		}
	}
}

func isV4(ip net.IP) bool {
	return ip.To4() != nil
}

func extractIpsFromInterface(iface *net.Interface, onlyV6 bool) ([]net.IP, error) {

	if iface.Flags&net.FlagUp == 0 {
		return nil, errors.New("interface is down.")
	}
	if iface.Flags&net.FlagLoopback != 0 {
		return nil, errors.New("interface is a loopback interface")
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	result := make([]net.IP, 0)

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip == nil || ip.IsLoopback() {
			continue
		}

		v4 := ip.To4()

		if onlyV6 {
			if v4 == nil {
				result = append(result, ip)
			}
		} else {
			if v4 != nil {
				result = append(result, ip)
			}
		}
	}
	return result, nil
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")

	return localAddr[0:idx]
}
