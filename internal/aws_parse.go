package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//Create the session var that will be used throught the package
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

func ParseRange(ec2IpRangeArray []*ec2.IpRange) (ipRangeArray []NetRange) {
	//Loop through each *ec2.IpRange object to build a NetRange array.
	for _, ipRange := range ec2IpRangeArray {
		cidrIp := strings.Split(*ipRange.CidrIp, "/")
		//Check and load into our "cache"
		var intIp int64
		if val, ok := NetCache[cidrIp[0]]; ok {
			intIp = val
		} else {
			intIp = GetIntFromIP(cidrIp[0])
			NetCache[cidrIp[0]] = intIp
		}
		//Create a new range object which is a subnet and
		//it's IP addresses corresponding IP in integer form.
		nRange := NetRange{cidrIp[0], cidrIp[1], intIp}
		ipRangeArray = append(ipRangeArray, nRange)
	}
	return ipRangeArray
}

/*
ParseIPPermissions will take the *ec2.IpPermission object and parse it into SecurityGroupRules.
Need to add checking on egress traffic.
*/
func ParseIPPermissions(perm []*ec2.IpPermission) (ipPermission []SecurityGroupRule) {
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
		sgRule := SecurityGroupRule{portRange, ipRangeArray}
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
		perms := ParseIPPermissions(sg.IpPermissions)

		sGroup := SecurityGroup{*sg.GroupName, perms}

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
