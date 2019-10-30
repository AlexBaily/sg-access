package internal

import (
	"strconv"
	"strings"
)

//NetCache will instantiate a singleton map for storing ranges in a cache.
var NetCache = make(map[string]int64)

type SecurityGroup struct {
	name  string
	rules []SecurityGroupRule
}

//SecurityGroupRule is the struct for the port and networks associated with that port
//We use a reference to NetRange as we don't know how many times this range will be used
//Should save some address space?
type SecurityGroupRule struct {
	ports    string
	networks []*NetRange
}

//NetRange is a struct that contains information about a network
type NetRange struct {
	cidr         string
	mask         string
	networkRange int64
}

func GetIntFromIP(cidrRange string) (i int64) {

	cidr := strings.Split(cidrRange, ".")
	//Parse each octet and convert to a base 10 int.
	quad1, _ := strconv.ParseInt(cidr[0], 10, 64)
	quad2, _ := strconv.ParseInt(cidr[1], 10, 64)
	quad3, _ := strconv.ParseInt(cidr[2], 10, 64)
	quad4, _ := strconv.ParseInt(cidr[3], 10, 64)

	//Shift each octet based on where it is positioned
	//The first octet needs to be shifted the most.
	return int64(quad1<<24 + quad2<<16 + quad3<<8 + quad4)
}

/*
To make this a function or not to make this a function, that is the question.
*/
func CompareIntIP(ipAddr int64, subnet NetRange) bool {
	mask, _ := strconv.ParseInt(subnet.mask, 10, 64)
	//Compare the normalised values to see if they match
	return (uint64(ipAddr) >> uint32(32-mask)) == (uint64(subnet.networkRange) >> uint32((32 - mask)))
}

/*
func main() {

	subnet := strings.Split("192.0.0.0/16", "/")
	subnetInt := GetIntFromIP(subnet[0])
	n := NetRange{subnet[0], subnet[1], subnetInt}

	ipToTest := "192.14.0.1"
	ipInt := GetIntFromIP(ipToTest)

	sameIP := CompareIntIP(ipInt, n)

	fmt.Println("%v", sameIP)
}
}*/
