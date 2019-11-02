package internal

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//Create the session var that will be used throught the package
var sess *session.Session

//Initliase the session in init()
func init() {
	sess = session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")}))
}

//QuerySecurityGroup will build a list of all SecurityGroups for parsing later.
func QuerySecurityGroups() []SecurityGroup {
	svc := ec2.New(sess)

	var sgArray []SecurityGroup
	sgInput := &ec2.DescribeSecurityGroupsInput{}

	for {

		securityGroups, err := svc.DescribeSecurityGroups(sgInput)

		for _, sg := range securityGroups.SecurityGroups {
			fmt.Printf("%v", sg)
		}
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

		if securityGroups.NextToken == nil {
			break
		}
	}
	return sgArray
}
