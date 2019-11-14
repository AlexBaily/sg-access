package api

import (
	"strconv"
	"testing"
)

/*
TestMostSpecificRoute will test the MostSpecificRoute() function to make sure
that it sets the highest specific route that matches the IP address specified.
*/
func TestMostSpecificRoute(t *testing.T) {

	//Set Variables erquired for the test.
	ip1 := GetIntFromIP("192.168.0.2")
	nR1 := NewNetRange("192.168.0.0", "16", 3232235520)
	nR2 := NewNetRange("10.0.0.0", "8", 167772160)
	nR3 := NewNetRange("0.0.0.0", "0", 0)

	nR1.RouteTableDestination = "pxc-000000"
	nR2.RouteTableDestination = "local"
	nR3.RouteTableDestination = "vgw-000000"

	routeTable := RouteTable{"rtb-test1", "vpc-test1", []NetRange{nR1, nR2, nR3}}

	MostSpecificRoute(ip1, &routeTable)

	if !routeTable.Routes[0].MostSpecific {
		t.Errorf("MostSpecificRoute(ip1, &routeTable) nR1 = %v; want true",
			routeTable.Routes[0].MostSpecific)
	}
	if routeTable.Routes[1].MostSpecific {
		t.Errorf("MostSpecificRoute(ip1, &routeTable) nR2 = %v; want false",
			routeTable.Routes[0].MostSpecific)
	}
	if routeTable.Routes[2].MostSpecific {
		t.Errorf("MostSpecificRoute(ip1, &routeTable) nR3 = %v; want false",
			routeTable.Routes[0].MostSpecific)
	}

}

/*
TestIsMostSpecific will test isMostSpecific to make sure that it returns true
when the second value is the higher priority path to take on an AWS Route Table.
*/
func TestIsMostSpecific(t *testing.T) {

	//Set Variables erquired for the test.
	nR1 := NewNetRange("10.0.0.0", "8", 167772160)
	nR2 := NewNetRange("10.0.0.0", "8", 167772160)
	nR3 := NewNetRange("10.0.0.0", "8", 167772160)

	nR1.RouteTableDestination = "igw-000000"
	nR2.RouteTableDestination = "vgw-000000"
	nR3.RouteTableDestination = "local"
	nR1.Propagated = true

	nR1Int, _ := strconv.Atoi(nR1.Mask)
	nR2Int, _ := strconv.Atoi(nR2.Mask)
	nR3Int, _ := strconv.Atoi(nR3.Mask)
	//Check to make sure that the non propagated route wins which is the second NetRange.
	trueCheck := isMoreSpecific(nR1Int, nR2Int, nR1, nR2)
	falseCheck := isMoreSpecific(nR3Int, nR2Int, nR3, nR2)

	if !trueCheck {
		t.Errorf("isMoreSpecific(nR1, nR2)= %v; want true",
			trueCheck)
	}
	if falseCheck {
		t.Errorf("isMoreSpecific(nR3, nR2)= %v; want false",
			falseCheck)
	}

}
