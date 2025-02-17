package eventgen

type Projection struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Type        string   `yaml:"type"`
	ReactsTo    []string `yaml:"reactsTo"`
}
