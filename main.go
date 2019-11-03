package main

import (
	"flag"
	"fmt"
	"sg-access/internal"

	"github.com/aws/aws-sdk-go/service/ec2"
)

//Way too many loops in this at the moment, need to split out parsing the groups etc.
func printMatchesTab(ipAddressInt int64, awsGroups []*ec2.DescribeSecurityGroupsOutput) {
	for _, group := range awsGroups {
		parsedGroups := internal.ParseSecurityGroups(group)
		for _, parsedGroup := range parsedGroups {
			for _, rule := range parsedGroup.Rules {
				for _, ipRange := range rule.Networks {
					if internal.CompareIntIP(ipAddressInt, ipRange) {
						fmt.Printf("%v\t%v\t%v/%v\n", parsedGroup.Name,
							rule.Ports, ipRange.Cidr, ipRange.Mask)
					}
				}

			}
		}

	}
}

func main() {

	ipAddress := flag.String("ip", "", "Required - IP Address to search SGs for.")
	printTab := flag.Bool("print-tab", false, "If flag is set then this will print the results tabulated.")
	flag.Parse()

	ipAddressInt := internal.GetIntFromIP(*ipAddress)
	awsGroups := internal.GetSecurityGroups()
	//Test to see what we're getting out of the Query.
	if *printTab {
		printMatchesTab(ipAddressInt, awsGroups)
	}

}
