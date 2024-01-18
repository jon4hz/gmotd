package config

type Config struct {
	Hostname Hostname
	Smart    Smart
}

type Hostname struct {
	Disabled      bool
	Color         string
	Figlet        bool
	FigletFont    string
	FigletFontDir string
}

type Smart struct {
	Disabled bool
	Disks    []string
}
