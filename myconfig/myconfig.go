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

func (self *IniConfig) GetInt64(key string) int64 {
	return self.GetSectionInt64("", key)
}

func (self *IniConfig) GetBool(key string) bool {
	re, _ := self.GetSectionBool("", key)
	return re
}

func (self *IniConfig) GetFloat64(key string) float64 {
	re, _ := self.GetSectionFloat64("", key)
	return re
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
		return false, fmt.Errorf("self.conf is null")
	}
	s := self.conf.Section(section)
	return s.Key(key).Bool()
}

func (self *IniConfig) GetSectionInt64(section string, key string) int64 {
	if self.conf == nil {
		return 0
	}
	s := self.conf.Section(section)
	v, _ := s.Key(key).Int64()
	return v
}

func (self *IniConfig) GetSectionFloat64(section string, key string) (float64, error) {
	if self.conf == nil {
		return 0, fmt.Errorf("self.conf is null")
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
