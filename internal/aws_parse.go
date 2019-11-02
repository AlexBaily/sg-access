package internal

import (
	"fmt"
	"os"

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

//QueryIPPermissions
func QueryIPPermissions(perm []*ec2.IpPermission) (ipPermission []NetRange) {
	return ipPermission
}

/*
QuerySecurityGroups will get the DescribSecurityGroupsOutput and parse it into the types that we want.
We can then take this output and pass it through to see if we get a match on the IP we want.
*/
/*func QuerySecurityGroups(securityGroups *ec2.DescribeSecurityGroupsOutput) (parsedGroup []SecurityGroup) {
	for _, sg := range securityGroups.SecurityGroups {
		perms := QueryIPPermissions(sg.IpPermissions)

		sGroup := SecurityGroupRule{*sg.GroupName, perms}
		parsedGroup = append(parsedGroup, sGroup)
	}
	return parsedGroup
} */

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
