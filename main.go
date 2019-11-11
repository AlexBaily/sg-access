package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"text/tabwriter"

	"github.com/AlexBaily/sg-access/api"

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
		parsedGroups := api.ParseSecurityGroups(group)
		for _, parsedGroup := range parsedGroups {
			for _, rule := range parsedGroup.Rules {
				for _, ipRange := range rule.Networks {
					if api.CompareIntIP(ipAddressInt, ipRange) {
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
		parsedTables := api.ParseRouteTables(table)
		for _, parsedTable := range parsedTables {
			for _, route := range parsedTable.Routes {
				if api.CompareIntIP(ipAddressInt, route) {
					fmt.Fprintln(tw, parsedTable.VpcID, "\t", parsedTable.RouteTableID, "\t",
						route.Cidr+"/"+route.Mask, "\t", route.RouteTableDestination, "\t")
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
		parsedTables := api.ParseRouteTables(table)
		for _, parsedTable := range parsedTables {
			for _, route := range parsedTable.Routes {
				if api.CompareIntIP(ipAddressInt, route) {
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

func main() {

	ipAddress := flag.String("ip", "", "Required - IP Address to search SGs for.")
	printTab := flag.Bool("print-tab", false, "If flag is set then this will print the results tabulated.")
	routes := flag.Bool("routes", false, "This will show the route tables associated to the IP address.")
	pretty := flag.Bool("pretty", false, "Set this flag to print out in a pretty format.")
	flag.Parse()

	if *ipAddress == "" {
		fmt.Println("Error, missing required argu	ment:")
		flag.PrintDefaults()
		os.Exit(1)
	} else if !isValidIP(*ipAddress) {
		fmt.Println("Error, please enter a valid IP address")
		flag.PrintDefaults()
		os.Exit(1)
	}

	ipAddressInt := api.GetIntFromIP(*ipAddress)

	if *routes {
		awsRoutes := api.GetRouteTables()
		if *printTab {
			printRouteMatches(ipAddressInt, awsRoutes)
		} else if *pretty {
			printRoutePretty(ipAddressInt, awsRoutes)
		}
	} else {
		awsGroups := api.GetSecurityGroups()
		if *printTab {
			printMatches(ipAddressInt, awsGroups, *pretty)
		}
	}

}
