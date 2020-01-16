package service

type MetadataCenter interface {
	GetProviderMetaData(m *MetadataIdentifier) *MetadataInfo
}
type MetadataIdentifier struct {
	ServiceInterface string
	Version          string
	Group            string
	Side             string
	Application      string
}
type MetadataInfo struct {
	Parameters    map[string]string
	CanonicalName string
	CodeSource    string
	Methods       []MethodDefinition
	Types         []TypeDefinition
}
type MethodDefinition struct {
	Name           string
	ParameterTypes []string
	ReturnType     string
}
type TypeDefinition struct {
	Id              string
	Type            string
	Items           []TypeDefinition
	Enums           []string
	Ref             string `json:"$ref"`
	TypeBuilderName string
}
