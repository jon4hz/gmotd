package config

type Config struct {
	Hostname *Hostname `mapstructure:"hostname"`
	Uptime   *Uptime   `mapstructure:"uptime"`
	SysInfo  *SysInfo  `mapstructure:"sysinfo"`
	Zpool    *Zpool    `mapstructure:"zpool"`
	Docker   *Docker   `mapstructure:"docker"`
	Smart    *Smart    `mapstructure:"smart"`
	Systemd  *Systemd  `mapstructure:"systemd"`
}

type Hostname struct {
	Disabled      bool   `mapstructure:"disabled"`
	Color         string `mapstructure:"color"`
	Figlet        bool   `mapstructure:"figlet"`
	FigletFont    string `mapstructure:"figlet_font"`
	FigletFontDir string `mapstructure:"figlet_font_dir"`
}

type Uptime struct {
	Disabled  bool `mapstructure:"disabled"`
	Precision int  `mapstructure:"precision"`
}

type SysInfo struct {
	Disabled        bool `mapstructure:"disabled"`
	UptimePrecision int  `mapstructure:"uptime_precision"`
}

type Zpool struct {
	Disabled bool `mapstructure:"disabled"`
}

type Docker struct {
	Disabled         bool     `mapstructure:"disabled"`
	IgnoreContainers []string `mapstructure:"ignore_containers"`
}

type Smart struct {
	Disabled bool     `mapstructure:"disabled"`
	Disks    []string `mapstructure:"disks"`
}

type Systemd struct {
	Disabled bool     `mapstructure:"disabled"`
	Units    []string `mapstructure:"units"`
}
