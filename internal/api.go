package internal

import (
	"strconv"
	"strings"
)

//NetCache will instantiate a singleton map for storing ranges in a cache.
var NetCache = make(map[string]int64)

type SecurityGroup struct {
	Name  string
	Rules []SecurityGroupRule
}

//SecurityGroupRule is the struct for the port and networks associated with that port
//We use a reference to NetRange as we don't know how many times this range will be used
//Should save some address space?
type SecurityGroupRule struct {
	Ports            string
	Networks         []NetRange
	TrafficDirection string
}

//NetRange is a struct that contains information about a network
type NetRange struct {
	Cidr         string
	Mask         string
	NetworkRange int64
}

func GetIntFromIP(ipAdrr string) (i int64) {

	octets := strings.Split(ipAdrr, ".")
	//Parse each octet and convert to a base 10 int.
	oct1, _ := strconv.ParseInt(octets[0], 10, 64)
	oct2, _ := strconv.ParseInt(octets[1], 10, 64)
	oct3, _ := strconv.ParseInt(octets[2], 10, 64)
	oct4, _ := strconv.ParseInt(octets[3], 10, 64)

	//Shift each octet based on where it is positioned
	//The first octet needs to be shifted the most.
	return int64(oct1<<24 + oct2<<16 + oct3<<8 + oct4)
}

/*
To make this a function or not to make this a function, that is the question.
*/
func CompareIntIP(ipAddr int64, subnet NetRange) bool {
	mask, _ := strconv.ParseInt(subnet.Mask, 10, 64)
	//Compare the normalised values to see if they match
	return (uint64(ipAddr) >> uint32(32-mask)) == (uint64(subnet.NetworkRange) >> uint32((32 - mask)))
}
