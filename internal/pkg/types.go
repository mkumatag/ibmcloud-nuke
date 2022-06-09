package pkg

type Nuke interface {
	NukeIt() error
	Validate() error
}

type Filter struct {
	Name  string `yaml:"name"`
	Regex string `yaml:"regex"`
	ID    string `yaml:"ID"`
}

type Resource struct {
	Nuke

	Type       string `yaml:"type"`
	Attributes `yaml:",inline"`
}

type Attributes struct {
	Region string `yaml:"region"`
	Zone   string `yaml:"zone"`
	Filter Filter `yaml:"filter"`
}

type Config struct {
	Global    Attributes `yaml:",inline"`
	Resources []Resource `yaml:"resources"`
}

func (c *Config) GetGlobalRegion() string {
	return c.Global.Region
}

func (c *Config) GetGlobalZone() string {
	return c.Global.Zone
}

func (c *Config) GetGlobalFilter() Filter {
	return c.Global.Filter
}
