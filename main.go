package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
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
func printMatchesTab(ipAddressInt int64, awsGroups []*ec2.DescribeSecurityGroupsOutput, pretty bool := false) {
	w := NewWriter(os.Stdout)
	for _, group := range awsGroups {
		parsedGroups := internal.ParseSecurityGroups(group)
		for _, parsedGroup := range parsedGroups {
			for _, rule := range parsedGroup.Rules {
				for _, ipRange := range rule.Networks {
					if internal.CompareIntIP(ipAddressInt, ipRange) {
						if pretty {
							
						} else {
							w.Write([]string{parsedGroup.Name, rule.TrafficDirection,
								rule.Ports, ipRange.Cidr})
						}

					}
				}

			}
		}

	}
	w.Flush()
}


func isValidIP(ip string) bool {
	pattern := "^(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\\." +
		"(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9][0-9]?)"
	matched, _ := regexp.MatchString(pattern, ip)
	return matched

}

func main() {

	ipAddress := flag.String("ip", "", "Required - IP Address to search SGs for.")
	printTab := flag.Bool("print-tab", false, "If flag is set then this will print the results tabulated.")
	flag.Parse()

	if *ipAddress == "" {
		fmt.Println("Error, missing required arguement:")
		flag.PrintDefaults()
		os.Exit(1)
	} else if !isValidIP(*ipAddress) {
		fmt.Println("Error, please enter a valid IP address")
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
