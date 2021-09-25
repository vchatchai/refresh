package model

type Config struct {
	URLRefresh string `yaml:"URLRefresh"`
	Database   struct {
		Debug        bool   `yaml:"Debug"`
		Password     string `yaml:"Password"`
		DatabasePort int    `yaml:"Port"`
		Server       string `yaml:"Server"`
		User         string `yaml:"User"`
		Database     string `yaml:"Name"`
	} `yaml:"Database"`

	Booking struct {
		HeaderFile string `yaml:"HeaderFile"`
		TypeFile   string `yaml:"TypeFile"`
		DetailFile string `yaml:"DetailFile"`

		Days int `yaml:"RefreshDays"`
	} `yaml:"Booking"`

	Laden struct {
		ContainerFile string `yaml:"ContainerFile"`
		ContainerUser string `yaml:"ContainerUser"`
	} `yaml:"Laden"`
}
