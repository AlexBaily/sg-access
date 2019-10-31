package internal

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
