package config

type Config struct {
	Hostname HostnameConfig
}

type HostnameConfig struct {
	Disable       bool
	Figlet        bool
	FigletFont    string
	FigletFontDir string
	FigletColor   string
}
