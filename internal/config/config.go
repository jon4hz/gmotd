package config

type Config struct {
	Hostname Hostname
}

type Hostname struct {
	Disable       bool
	Color         string
	Figlet        bool
	FigletFont    string
	FigletFontDir string
}
