package configcenter

type (
	ConfigClient interface {
		LoadConfig(properties map[string]interface{}) (string, error)

		ListenConfig(properties map[string]interface{}) (err error)
	}

	ListenConfig func(data string)
)
