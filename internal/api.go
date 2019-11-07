package internal

import (
	"strconv"
	"strings"
)

//NetCache will instantiate a singleton map for storing ranges in a cache.
var NetCache = make(map[string]int64)

//SecurityGroup struct that will house all of the SecurityGroupRule objects.
type SecurityGroup struct {
	Name  string
	VpcID string
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
	Cidr                  string
	Mask                  string
	NetworkRange          int64
	RouteTableDestination string
}

//RouteTable is a struct that contains information on an individual RouteTable
type RouteTable struct {
	RouteTableID string
	VpcID        string
	Routes       []NetRange
}

/*
NewNetRange will be the interface we use to create NetRange objects.
This is because we want to reuse the NetRange type for both SG and RouteTables.
RouteTableDestination is not require on SG so we give it a default here.
*/
func NewNetRange(Cidr string, Mask string, NetworkRange int64) NetRange {
	n := NetRange{}
	n.Cidr = Cidr
	n.Mask = Mask
	n.NetworkRange = NetworkRange
	n.RouteTableDestination = ""
	return n
}

/*
GetIntFromIP take an IP address and converts it into a 64bit integer.
*/
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
CompareIntIP will compare the IP address to the NetRange to see if they share the same Network address.
*/
func CompareIntIP(ipAddr int64, subnet NetRange) bool {
	mask, _ := strconv.ParseInt(subnet.Mask, 10, 64)
	//Compare the normalised values to see if they match
	return (uint64(ipAddr) >> uint32(32-mask)) == (uint64(subnet.NetworkRange) >> uint32((32 - mask)))
}
