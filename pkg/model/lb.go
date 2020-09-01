package model

// LbPolicy the load balance policy enum
type LbPolicy int32

const (
	RoundRobin LbPolicy = 0
	IPHash     LbPolicy = 1
	WightRobin LbPolicy = 2
	Rand       LbPolicy = 3
)

// LbPolicyName key int32 for LbPolicy, value string
var LbPolicyName = map[int32]string{
	0: "RoundRobin",
	1: "IPHash",
	2: "WightRobin",
	3: "Rand",
}

// LbPolicyValue key string, value int32 for LbPolicy
var LbPolicyValue = map[string]int32{
	"RoundRobin": 0,
	"IPHash":     1,
	"WightRobin": 2,
	"Rand":       3,
}
