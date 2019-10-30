package internal

import (
	"strconv"
	"strings"
)

//NetCache will instantiate a singleton map for storing ranges in a cache.
var NetCache = make(map[string][]string)

type SecurityGroup struct {
	name  string
	rules []SecurityGroupRule
}

type SecurityGroupRule struct {
	ports    string
	networks []NetRange
}

type NetRange struct {
	cidr         string
	networkRange []string
}

func GetIntFromIP(cidrRange string) (i int64) {
	cidr := strings.Split(cidrRange, ".")

	quad1, _ := strconv.ParseInt(cidr[0], 10, 64)
	quad2, _ := strconv.ParseInt(cidr[1], 10, 64)
	quad3, _ := strconv.ParseInt(cidr[2], 10, 64)
	quad4, _ := strconv.ParseInt(cidr[3], 10, 64)

	return int64(quad1<<24 + quad2<<16 + quad3<<8 + quad4)
}

/*
func CheckRange(cidr string) (n NetRange, b bool) {

}*/
