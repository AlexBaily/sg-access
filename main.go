package main

import (
	"fmt"
	"sg-access/internal"
)

func main() {
	awsGroups := internal.GetSecurityGroups()
	//Test to see what we're getting out of the Query.
	for _, group := range awsGroups {
		test1 := internal.ParseSecurityGroups(group)
		fmt.Printf("%v", test1)

	}

}
