package model

// RemoteConfig remote server info which offer server discovery or k/v or distribution function
type RemoteConfig struct {
	Protocol string `yaml:"protocol" json:"protocol" default:"zookeeper"`
	Timeout  string `yaml:"timeout" json:"timeout"`
	Address  string `yaml:"address" json:"address"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}
