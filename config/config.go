package config

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path"
	"gopkg.in/yaml.v2"
	"github.com/BurntSushi/toml"
)

type Config struct {
	Monitor struct {
		IP string
		Port int
	}
	Database struct {
		Driver string
		Dsn string
		ShowSQL bool
		Migrate bool
	}
	HTTPServer struct {
		IP string
		Port int
	}
}

func Parse(file string) (*Config, error) {
	var c Config
	
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	
	switch path.Ext(file) {
	case ".yml", ".yaml":
		yaml.Unmarshal(buf, &c)
	case ".toml":
		_, err = toml.DecodeReader(bytes.NewReader(buf), &c)
	default:
		return nil, errors.New("file type not suport, use .yml .yaml or .toml instead")
	}

	if err != nil {
		return nil, err
	}

	return &c, nil
}
