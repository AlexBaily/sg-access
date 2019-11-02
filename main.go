package main

import (
	"fmt"
	"sg-access/internal"
)

func main() {
	awsGroups := internal.GetSecurityGroups()
	fmt.Printf("%v", awsGroups)
}
