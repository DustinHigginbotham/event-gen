package eventgen

type DomainSchema struct {
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Entity      Entity       `yaml:"entity"`
	Commands    []Command    `yaml:"commands"`
	Events      []Event      `yaml:"events"`
	Reactors    []Reactor    `yaml:"reactors"`
	Projections []Projection `yaml:"projections"`

	Package string `yaml:"-"`
}
