package model

type AuthorityConfiguration struct {
	Rules []AuthorityRule `yaml:"authority_rule" json:"authority_rule"` //Rules the authority rule list
}

type AuthorityRule struct {
	Strategy StrategyType `yaml:"strategy" json:"strategy"` //Strategy the authority rule strategy
	Limit    LimitType    `yaml:"limit" json:"limit"`       //Limit the authority rule limit
	Items    []string     `yaml:"items" json:"items"`       //Items the authority rule items
}

// StrategyType the authority rule strategy enum
type StrategyType int32

const (
	Whitelist StrategyType = 0
	Blacklist StrategyType = 1
)

// StrategyTypeName key int32 for StrategyType, value string
var StrategyTypeName = map[int32]string{
	0: "Whitelist",
	1: "Blacklist",
}

// StrategyTypeValue key string, value int32 for StrategyType
var StrategyTypeValue = map[string]int32{
	"Whitelist": 0,
	"Blacklist": 1,
}

// LimitType the authority rule limit enum
type LimitType int32

const (
	IP  LimitType = 0
	App LimitType = 1
)

// LimitTypeName key int32 for LimitType, value string
var LimitTypeName = map[int32]string{
	0: "IP",
	1: "App",
}

// LimitTypeValue key string, value int32 for LimitType
var LimitTypeValue = map[string]int32{
	"IP":  0,
	"App": 1,
}
