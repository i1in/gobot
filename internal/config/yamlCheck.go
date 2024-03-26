package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var c *Config

func C() *Config {
	if c == nil {
		cfg, err := LoadConfig()
		if err != nil {
			log.Fatalln(err)
		}
		c = cfg
	}
	return c
}

type RoleType string

type RolePermission string

const (
	RoleAdmin RoleType = "admin"
	RoleModer RoleType = "moder"

	PermBan  RolePermission = "ban"
	PermKick RolePermission = "kick"
	PermMute RolePermission = "mute"
)

type RolesConfig struct {
	Id     string           `yaml:"id"`
	Type   RoleType         `yaml:"type"`
	Parent []RoleType       `yaml:"parent"`
	Perm   []RolePermission `yaml:"perm"`
}

type Config struct {
	Roles []RolesConfig `yaml:"roles"`
}

func (r *RolesConfig) Run() {
	l, err := LoadConfig()
	if err != nil {
		log.Fatal("%s", err)
	}

	log.Println(l)
}

func LoadConfig() (*Config, error) {
	y := &Config{}

	f, err := os.ReadFile("testRoles.yaml")
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(f, &y); err != nil {
		return nil, err
	}

	return y, nil
}
