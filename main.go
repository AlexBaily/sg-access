package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"

	sg "github.com/AlexBaily/sg-access/api"

	"github.com/aws/aws-sdk-go/service/ec2"
)

//Get a CSV writer to write out tab delimited CSV
func newWriter(w io.Writer) (writer *csv.Writer) {
	writer = csv.NewWriter(w)
	writer.Comma = '\t'

	return
}

//Way too many loops in this at the moment, need to split out parsing the groups etc.
func printMatches(ipAddressInt int64, awsGroups []*ec2.DescribeSecurityGroupsOutput, pretty bool) {
	w := newWriter(os.Stdout)
	for _, group := range awsGroups {
		parsedGroups := sg.ParseSecurityGroups(group)
		for _, parsedGroup := range parsedGroups {
			for _, rule := range parsedGroup.Rules {
				for _, ipRange := range rule.Networks {
					if sg.CompareIntIP(ipAddressInt, ipRange) {
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

func printRoutePretty(ipAddressInt int64, awsRoutes []*ec2.DescribeRouteTablesOutput) {
	tw := new(tabwriter.Writer)
	tw.Init(os.Stdout, 0, 0, 8, ' ', tabwriter.Debug|tabwriter.AlignRight)
	fmt.Fprintln(tw, "VPC\tRoute Table ID\tCIDR\tDestination\t")
	for _, table := range awsRoutes {
		parsedTables := sg.ParseRouteTables(table)
		for _, parsedTable := range parsedTables {
			//Pass in the reference to the parsedTable so we can set the most specific route.
			sg.MostSpecificRoute(ipAddressInt, &parsedTable)
			for _, route := range parsedTable.Routes {
				if sg.CompareIntIP(ipAddressInt, route) {
					fmt.Fprintln(tw, parsedTable.VpcID, "\t", parsedTable.RouteTableID, "\t",
						route.Cidr+"/"+route.Mask, "\t", route.RouteTableDestination, "\t",
						route.MostSpecific, "\t")
					fmt.Fprintln(tw)
				}

			}
		}

	}
	tw.Flush()
}

//This is basically the same as printMatches but for Routes.
func printRouteMatches(ipAddressInt int64, awsRoutes []*ec2.DescribeRouteTablesOutput) {
	w := newWriter(os.Stdout)
	for _, table := range awsRoutes {
		parsedTables := sg.ParseRouteTables(table)
		for _, parsedTable := range parsedTables {
			for _, route := range parsedTable.Routes {
				if sg.CompareIntIP(ipAddressInt, route) {
					w.Write([]string{parsedTable.RouteTableID, route.Cidr + "/" + route.Mask,
						route.RouteTableDestination})
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

//Parse IPs will parse multiple IPs based off of the IP addresses given in the command line args.
func parseIPs(ipAddresses *string, routes *bool, printTab *bool, pretty *bool) {
	//Parse multiple IPs based on
	for _, ipAddress := range strings.Split(*ipAddresses, ",") {
		ipAddress = strings.TrimSpace(ipAddress)
		if !isValidIP(ipAddress) {
			fmt.Println("Error, please enter a valid IP address")
			flag.PrintDefaults()
			os.Exit(1)
		}
		//Get the integer of the IP address.
		ipAddressInt := sg.GetIntFromIP(ipAddress)
		//Check if are going to parse routes rather than SG.
		if *routes {
			awsRoutes := sg.GetRouteTables()
			if *printTab {
				printRouteMatches(ipAddressInt, awsRoutes)
			} else if *pretty {
				printRoutePretty(ipAddressInt, awsRoutes)
			}
		} else {
			awsGroups := sg.GetSecurityGroups()
			if *printTab {
				printMatches(ipAddressInt, awsGroups, *pretty)
			}
		}
	}
}

func main() {

	//Parse all of the flags.
	ipAddresses := flag.String("ip", "", "Required - IP Address to search SGs for.")
	printTab := flag.Bool("print-tab", false, "If flag is set then this will print the results tabulated.")
	routes := flag.Bool("routes", false, "This will show the route tables associated to the IP address.")
	pretty := flag.Bool("pretty", false, "Set this flag to print out in a pretty format.")
	flag.Parse()

	//Check if we don't have an IP or if it's invalid.
	if *ipAddresses == "" {
		fmt.Println("Error, missing required argu	ment:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	parseIPs(ipAddresses, routes, printTab, pretty)

}
