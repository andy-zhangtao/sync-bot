package types

import "github.com/BurntSushi/toml"

// Use .toml as default config file.
// The content of config file as below:
// [run]
// 	  token = "bot api token"
//    debug = true/false
// 	  interval = 1 (query message in one second interval)

type Config struct {
	Run     Run          `toml:"run"`
	GHelper GitHubHelper `toml:"github"`
}

type Run struct {
	Token    string `toml:"token"`
	Debug    bool   `toml:"debug"`
	Interval int    `toml:"interval"`
}

type GitHubHelper struct {
	Token string `toml:"token"`
	Repo  string `toml:"repo"`
	Auth  string `toml:"auth"`
	Email string `toml:"email"`
}

func ParserConfig(path string) (c Config, err error) {
	_, err = toml.DecodeFile(path, &c)
	return
}
