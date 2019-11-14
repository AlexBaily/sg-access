package api

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//Create the session var that will be used throughout the package
var sess *session.Session

//Initliase the session in init()
func init() {
	var region string
	//If there is not default region in environ then just set it to eu-west-1 for the moment.
	if os.Getenv("AWS_DEFAULT_REGION") == "" {
		region = "eu-west-1"
	} else {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}
	sess = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region)}))
}

//CheckCache will take a string and check if we have already worked out the int64 version of the IP.
func CheckCache(cidrIP string) (intIP int64) {
	if val, ok := NetCache[cidrIP]; ok {
		intIP = val
	} else {
		intIP = GetIntFromIP(cidrIP)
		NetCache[cidrIP] = intIP
	}
	return intIP
}

/*
ParseRange takes an []*ec2.IpRange parses it and convert it into a []NetRange array.
*/
func ParseRange(ec2IpRangeArray []*ec2.IpRange) (ipRangeArray []NetRange) {
	//Loop through each *ec2.IpRange object to build a NetRange array.
	for _, ipRange := range ec2IpRangeArray {
		cidrIP := strings.Split(*ipRange.CidrIp, "/")
		//Check and load into our "cache"
		intIP := CheckCache(cidrIP[0])

		//Create a new range object which is a subnet and
		//it's IP addresses corresponding IP in integer form.
		nRange := NewNetRange(cidrIP[0], cidrIP[1], intIP)
		ipRangeArray = append(ipRangeArray, nRange)
	}
	return ipRangeArray
}

/*
ParseIPPermissions will take the *ec2.IpPermission object and parse it into SecurityGroupRules.
Need to add checking on egress traffic.
*/
func ParseIPPermissions(perm []*ec2.IpPermission, trafficD string) (ipPermission []SecurityGroupRule) {
	//Loop through every permission and build a SecurityGroupRule for it.
	for _, permission := range perm {
		//Loop through all of the IpRanges and build the IPRange list for the
		//SecurityGroupRule type.
		ipRangeArray := ParseRange(permission.IpRanges)
		var portRange string
		//Get the port range, create the sgRule object and add to the ipPermission Array
		if *permission.IpProtocol != "-1" {
			portRange = (strconv.FormatInt(*permission.FromPort, 10) + "-" +
				strconv.FormatInt(*permission.ToPort, 10))
		} else {
			portRange = "all"
		}
		//Create a rule object and add to the Array.
		sgRule := SecurityGroupRule{portRange, ipRangeArray, trafficD}
		ipPermission = append(ipPermission, sgRule)
	}

	return ipPermission
}

/*
ParseSecurityGroups will get the DescribSecurityGroupsOutput and parse it into the types that we want.
We can then take this output and pass it through to see if we get a match on the IP we want.
*/
func ParseSecurityGroups(securityGroups *ec2.DescribeSecurityGroupsOutput) (parsedGroup []SecurityGroup) {
	for _, sg := range securityGroups.SecurityGroups {
		permsIngress := ParseIPPermissions(sg.IpPermissions, "ingress")
		permsEgress := ParseIPPermissions(sg.IpPermissions, "egress")
		//Adding three dots to the end of permsEgress will add the whole slice together.
		perms := append(permsIngress, permsEgress...)
		sGroup := SecurityGroup{*sg.GroupName, *sg.VpcId, perms}

		parsedGroup = append(parsedGroup, sGroup)
	}
	return parsedGroup
}

/*
GetSecurityGroups will build a list of all SecurityGroups for parsing later.
We can change this function in the future to specify which region we want to use
Or we can set it so that it uses scans all regions.
*/
func GetSecurityGroups() []*ec2.DescribeSecurityGroupsOutput {
	svc := ec2.New(sess)

	var sgArray []*ec2.DescribeSecurityGroupsOutput
	sgInput := &ec2.DescribeSecurityGroupsInput{}

	for {

		securityGroups, err := svc.DescribeSecurityGroups(sgInput)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
		}
		sgArray = append(sgArray, securityGroups)

		if securityGroups.NextToken == nil {
			break
		}
	}
	return sgArray
}

/*
GetRouteTables will build a list of all RouteTables for parsing later.
We can change this function in the future to specify which region we want to use
Or we can set it so that it uses scans all regions.
*/
func GetRouteTables() []*ec2.DescribeRouteTablesOutput {
	svc := ec2.New(sess)

	var rtArray []*ec2.DescribeRouteTablesOutput
	rtInput := &ec2.DescribeRouteTablesInput{}

	for {

		routeTables, err := svc.DescribeRouteTables(rtInput)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
		}
		rtArray = append(rtArray, routeTables)

		if routeTables.NextToken == nil {
			break
		}
	}
	return rtArray

}

/*
ParseRouteTables will take *ec2.DescribeRouteTablesOutput and output a parse RoutTable Array.
*/
func ParseRouteTables(routeTables *ec2.DescribeRouteTablesOutput) (parsedTable []RouteTable) {
	for _, rt := range routeTables.RouteTables {

		routes := ParseRoutes(rt.Routes)
		routeTable := RouteTable{*rt.RouteTableId, *rt.VpcId, routes}

		parsedTable = append(parsedTable, routeTable)
	}
	return parsedTable
}

/*
ParseRoutes will take []*ec2.Route and convert it to a []NetRange
*/
func ParseRoutes(routes []*ec2.Route) (parsedRoutes []NetRange) {
	for _, route := range routes {
		if route.DestinationCidrBlock != nil {
			//Find the next hop destination
			dest := ParseRouteDestination(*route)
			//Parse the CIDR range

			cidrIP := strings.Split(*route.DestinationCidrBlock, "/")
			//Check and load into our "cache"
			intIP := CheckCache(cidrIP[0])

			nRange := NewNetRange(cidrIP[0], cidrIP[1], intIP)
			nRange.RouteTableDestination = dest
			if *route.Origin == "EnableVgwRoutePropagation" {
				nRange.Propagated = true
			}
			parsedRoutes = append(parsedRoutes, nRange)
		}
	}
	return parsedRoutes
}

/*
ParseRouteDestination will look at the ec2.Route type and determine what the destination is .e.g.
VPG, GatewayId, InstanceId, NateGatewayID
We find this information by using reflection to get all of the fields,
then we can exclude fields we don't need and just search for the field that contains data
*/
func ParseRouteDestination(route ec2.Route) (dest string) {
	r := reflect.ValueOf(route)
	for i := 0; i < r.NumField(); i++ {
		//Check to see if the field is a pointer, this is because all ec2.Route fields are *strings.
		if r.Field(i).Kind() == reflect.Ptr {
			//Check if the field IsValid() this will check if we have a nil pointer,
			//if not make sure it's not one of the fields we don't care about.
			if r.Field(i).Elem().IsValid() && (r.Type().Field(i).Name != "DestinationCidrBlock" &&
				r.Type().Field(i).Name != "DestinationIpv6CidrBlock" && r.Type().Field(i).Name != "Origin" &&
				r.Type().Field(i).Name != "State") {
				dest = r.Field(i).Elem().String()
				break
			}
		}
	}
	return dest
}

/*
MostSpecificRoute will take an IP address and a dereferenced RouteTable
It will then see which one of the routes in the table is the most specific match.
*/
func MostSpecificRoute(ipAddressInt int64, table *RouteTable) {
	//Start off with declaring some variables to be used in the loops.
	var mostSpecific *NetRange
	var msInt int //msInt is [mostSpecific] Mask, declared here as we need to use in loop.
	for i := range table.Routes {
		if CompareIntIP(ipAddressInt, table.Routes[i]) {
			//This is the current Routes (i) Mask as int
			rmInt, _ := strconv.Atoi(table.Routes[i].Mask)
			//If mostSpecific is nil then this is the first match in the loop and we need
			//to set an initial value.
			if mostSpecific == nil {
				//This syntax is required to get the reference of the Route in our
				//[]NetRange in *RouteTable else we will not pass the reference.
				mostSpecific = (&table.Routes[i])
				mostSpecific.MostSpecific = true
				msInt, _ = strconv.Atoi(mostSpecific.Mask)
				//If the new Route Mask is more larger (more specific) than the current
				//mostSpecific then set it to the new one, first we need to clear the current
				//mostSpecific.MostSpecific to false as it is no longer the most specific.
			} else if rmInt > msInt {
				mostSpecific.MostSpecific = false
				mostSpecific = (&table.Routes[i])
				mostSpecific.MostSpecific = true
				msInt, _ = strconv.Atoi(mostSpecific.Mask)
			} else if rmInt == msInt {
				//If the routes have an equal prefix then we need to figure out the tiebreaker.
				if isMoreSpecific(*mostSpecific, table.Routes[i]) {
					mostSpecific.MostSpecific = false
					mostSpecific = (&table.Routes[i])
					mostSpecific.MostSpecific = true
					msInt, _ = strconv.Atoi(mostSpecific.Mask)
				}
			}
		}
	}

}

/*
isMoreSpecific will check the current mostSpecific NetRange and compare
it against the new route, it will first check route prepagation and then
the destination to see where it is going.
We will eventually want to see if we can tell if it's from BGP or a static VPN route.
BGP VPN routes go first.
*/
func isMoreSpecific(currMS NetRange, msToCompare NetRange) bool {
	//If they are both not propagated then we will check to see where they are going.
	//Local routes get precendence, then static routes and then propagated routes.
	//Of the propagated routes it goes DX BGP -> VPN BGP -> VPN Static.
	//There is no way to where the route came from when it's a VGW.
	//https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Route_Tables.html#route-tables-priority
	if currMS.RouteTableDestination[:3] == "loc" {
		return false
	} else if msToCompare.RouteTableDestination[:3] == "loc" {
		return true
	} else if currMS.Propagated && !msToCompare.Propagated {
		return false
	} else {
		//!currMS.Propagated && msToCompare.Propagated should be the final comparison
		//There will not be anymore comparisons after this due to the aformentioned issues with VGW.
		return true
	}
}
