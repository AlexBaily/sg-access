package internal

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
func CheckCache(cidrIp string) (intIp int64) {
	if val, ok := NetCache[cidrIp]; ok {
		intIp = val
	} else {
		intIp = GetIntFromIP(cidrIp)
		NetCache[cidrIp] = intIp
	}
	return intIp
}

/*
ParseRange takes an []*ec2.IpRange parses it and convert it into a []NetRange array.
*/
func ParseRange(ec2IpRangeArray []*ec2.IpRange) (ipRangeArray []NetRange) {
	//Loop through each *ec2.IpRange object to build a NetRange array.
	for _, ipRange := range ec2IpRangeArray {
		cidrIp := strings.Split(*ipRange.CidrIp, "/")
		//Check and load into our "cache"
		intIp := CheckCache(cidrIp[0])

		//Create a new range object which is a subnet and
		//it's IP addresses corresponding IP in integer form.
		nRange := NewNetRange(cidrIp[0], cidrIp[1], intIp)
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
GetSecurityGroups will build a list of all RouteTables for parsing later.
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
ParsedRouteTables will take *ec2.DescribeRouteTablesOutput and output a parse RoutTable Array.
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

 */
func ParseRoutes(routes []*ec2.Route) (parsedRoutes []NetRange) {
	for _, route := range routes {
		//Find the next hop destination
		dest := ParseRouteDestination(*route)
		//Parse the CIDR range
		cidrIp := strings.Split(*route.DestinationCidrBlock, "/")
		//Check and load into our "cache"
		intIP := CheckCache(cidrIp[0])

		nRange := NewNetRange(cidrIp[0], cidrIp[1], intIP)
		nRange.RouteTableDestination = dest
		parsedRoutes = append(parsedRoutes, nRange)
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

func MostSpecificRoute()