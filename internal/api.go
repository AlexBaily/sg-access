package internal

type SecurityGroup struct {
    name string
    rules[] SecurityGroupRule
}

type SecurityGroupRule struct {
    ports string
    networks[] NetRange
}


type NetRange struct {
    cidr string
    range[]
}


//This "cache" will be 
type NetCache map[string]string[]