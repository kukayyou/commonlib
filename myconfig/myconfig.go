package myconfig

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
)

type IniConfig struct {
	conf *ini.File
}

func (self *IniConfig) Load(filename string) {
	conf, err := ini.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
	self.conf = conf
}

func (self *IniConfig) GetString(key string) string {
	return self.GetSectionString("", key)
}

func (self *IniConfig) GetInt64(key string) (int64, error) {
	return self.GetSectionInt64("", key)
}

func (self *IniConfig) GetInt(key string) (int, error) {
	return self.GetSectionInt("", key)
}

func (self *IniConfig) GetBool(key string) (bool, error) {
	return self.GetSectionBool("", key)
}

func (self *IniConfig) GetFloat64(key string) (float64, error) {
	return self.GetSectionFloat64("", key)
}

func (self *IniConfig) GetSectionString(section string, key string) string {
	if self.conf == nil {
		return ""
	}
	s := self.conf.Section(section)
	return s.Key(key).String()
}

func (self *IniConfig) GetSectionBool(section string, key string) (bool, error) {
	if self.conf == nil {
		return false, fmt.Errorf("config is not init")
	}
	s := self.conf.Section(section)
	return s.Key(key).Bool()
}

func (self *IniConfig) GetSectionInt64(section string, key string) (int64, error) {
	if self.conf == nil {
		return 0, fmt.Errorf("config is not init")
	}
	s := self.conf.Section(section)
	return s.Key(key).Int64()
}

func (self *IniConfig) GetSectionInt(section string, key string) (int, error) {
	if self.conf == nil {
		return 0, fmt.Errorf("config is not init")
	}
	s := self.conf.Section(section)
	return s.Key(key).Int()
}

func (self *IniConfig) GetSectionFloat64(section string, key string) (float64, error) {
	if self.conf == nil {
		return 0, fmt.Errorf("config is not init")
	}
	s := self.conf.Section(section)
	return s.Key(key).Float64()
}

var Config *IniConfig

func init() {
	Config = &IniConfig{}
}

func LoadConfig(filename string) {
	Config.Load(filename)
}
