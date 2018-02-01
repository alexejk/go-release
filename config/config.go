package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var p *viper.Viper

func init() {
	Reset()
}

func LoadConfig(name string) error {
	p.SetConfigFile(name)
	return p.ReadInConfig()
}

func BindFlag(name string, flagSet *pflag.FlagSet) {
	p.BindPFlag(name, flagSet.Lookup(name))
}

func Unmarshal(key string, rawVal interface{}) error {

	return p.UnmarshalKey(key, rawVal)
}

func IsSet(key string) bool {

	return p.IsSet(key)
}

func Reset() {

	p = viper.New()
	p.SetConfigType("yaml")
}

func Set(key string, value interface{}) {
	p.Set(key, value)
}

func Get(key string) interface{} {
	return p.Get(key)
}

func GetBool(key string) bool {
	return p.GetBool(key)
}

func GetString(key string) string {
	return p.GetString(key)
}

func GetInt64(key string) int64 {
	return p.GetInt64(key)
}
