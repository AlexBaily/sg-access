package api

import "testing"

/*
TestIntIPCalc will test the GetIntFromIP() function to make sure it produces
the correct integer from the specified IP address.
*/
func TestIntIPCalc(t *testing.T) {
	ip1 := GetIntFromIP("192.168.0.2")
	if ip1 != 3232235522 {
		t.Errorf("GetIntFromIP(\"192.168.0.2\") = %d; want 3232235522",
			ip1)
	}
	ip2 := GetIntFromIP("10.103.2.5")
	if ip2 != 174522885 {
		t.Errorf("GetIntFromIP(\"10.103.2.5\") = %d; want 174522885",
			ip2)
	}
}

/*
TestIPCompare will test comparing a NetRange object to an IP address
They should test true if the subnet is in the NetRange subnet and false
if it is not.
*/
func TestIPCompare(t *testing.T) {
	ip1 := GetIntFromIP("192.168.0.2")
	n := NewNetRange("192.168.0.0", "16", 3232235520)
	compTrue := CompareIntIP(ip1, n)
	if compTrue != true {
		t.Errorf("CompareIntIP(\"192.168.0.2\", n) = %v; want true",
			compTrue)
	}
	ip2 := GetIntFromIP("172.16.0.2")
	compFalse := CompareIntIP(ip2, n)
	if compFalse != false {
		t.Errorf("CompareIntIP(\"192.168.0.2\", n) = %v; want false",
			compFalse)
	}

}
