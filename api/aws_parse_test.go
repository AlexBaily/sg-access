package api

import "testing"

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
	nR1 := NewNetRange("10.0.0.0", "0", 167772160)
	nR2 := NewNetRange("0.0.0.0", "0", 0)

	nR1.RouteTableDestination = "pxc-000000"
	nR2.RouteTableDestination = "local"

	trueCheck := isMoreSpecific(nR1, nR2)

	if !trueCheck {
		t.Errorf("isMoreSpecific(nR1, nR2)= %v; want true",
			trueCheck)
	}

}
