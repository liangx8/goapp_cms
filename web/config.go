package web
import (
	"time"
	yaml "gopkg.in/yaml.v2"
)
type (
	Config struct{
		Passphase      string    `yaml:"pass-phase"`
		SourcesPrefix  string    `yaml:"source-prefix"`
		SourcesPkg     string    `yaml:"sources-pkg"`
		ResetId        string    `yaml:"reset-id"` // passPhase is expected a reset?
		ResetTimeout   time.Time `yaml:"reset-timeout"`
	}
)
var defaultConfig =Config{
	Passphase:     "",
	SourcesPrefix: "web",
	SourcesPkg:    "src.zip",
	ResetId:       "",
}
func (c *Config)Yaml()([]byte,error){
	return yaml.Marshal(c)
}
