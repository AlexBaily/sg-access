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
		t.Errorf("MostSpecificRoute(ip1, &routeTable) = %v; want true",
			routeTable.Routes[0].MostSpecific)
	}
	if routeTable.Routes[1].MostSpecific {
		t.Errorf("MostSpecificRoute(ip1, &routeTable) = %v; want false",
			routeTable.Routes[0].MostSpecific)
	}
	if routeTable.Routes[2].MostSpecific {
		t.Errorf("MostSpecificRoute(ip1, &routeTable) = %v; want false",
			routeTable.Routes[0].MostSpecific)
	}

}
