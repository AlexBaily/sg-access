package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"sg-access/internal"

	"github.com/aws/aws-sdk-go/service/ec2"
)

//Get a CSV writer to write out tab delimited CSV
func NewWriter(w io.Writer) (writer *csv.Writer) {
	writer = csv.NewWriter(w)
	writer.Comma = '\t'

	return
}

//Way too many loops in this at the moment, need to split out parsing the groups etc.
func printMatchesTab(ipAddressInt int64, awsGroups []*ec2.DescribeSecurityGroupsOutput) {
	w := NewWriter(os.Stdout)
	for _, group := range awsGroups {
		parsedGroups := internal.ParseSecurityGroups(group)
		for _, parsedGroup := range parsedGroups {
			for _, rule := range parsedGroup.Rules {
				for _, ipRange := range rule.Networks {
					if internal.CompareIntIP(ipAddressInt, ipRange) {
						w.Write([]string{parsedGroup.Name, rule.TrafficDirection,
							rule.Ports, ipRange.Cidr, ipRange.Mask})
					}
				}

			}
		}

	}
	w.Flush()
}

func main() {

	ipAddress := flag.String("ip", "", "Required - IP Address to search SGs for.")
	printTab := flag.Bool("print-tab", false, "If flag is set then this will print the results tabulated.")
	flag.Parse()

	if *ipAddress == "" {
		fmt.Println("Error, missing required arguement:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	ipAddressInt := internal.GetIntFromIP(*ipAddress)
	awsGroups := internal.GetSecurityGroups()
	//Test to see what we're getting out of the Query.
	if *printTab {
		printMatchesTab(ipAddressInt, awsGroups)
	}

}
