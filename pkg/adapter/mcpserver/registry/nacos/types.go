package nacos

// Nacos configuration data structures
type NacosTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

type RequestTemplate struct {
	URL            string              `json:"url"`
	Method         string              `json:"method"`
	Headers        []map[string]string `json:"headers"`
	ArgsToJsonBody bool                `json:"argsToJsonBody"`
	ArgsToUrlParam bool                `json:"argsToUrlParam"`
}

type ResponseTemplate struct {
	PrependBody string `json:"prependBody"`
}

type JsonGoTemplate struct {
	RequestTemplate  RequestTemplate  `json:"requestTemplate"`
	ResponseTemplate ResponseTemplate `json:"responseTemplate"`
}

type ToolMeta struct {
	Enabled   bool                   `json:"enabled"`
	Templates map[string]interface{} `json:"templates"`
}

type ToolsSpec struct {
	Tools     []NacosTool         `json:"tools"`
	ToolsMeta map[string]ToolMeta `json:"toolsMeta"`
}
