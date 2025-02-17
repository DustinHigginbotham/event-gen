package eventgen

type Command struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Handler     string  `yaml:"handler"`
	Emits       string  `yaml:"emits"`
	Fields      []Field `yaml:"fields"`
}
