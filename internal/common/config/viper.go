package config

import (
	"github.com/spf13/viper"
	"strings"
)

func NewViperConfig() error {
	viper.SetConfigName("global")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../common/config")
	viper.EnvKeyReplacer(strings.NewReplacer("-", "_"))
	_ = viper.BindEnv("stripe-key", "STRIPE_KEY", "endpoint-stripe-secret", "ENDPOINT_STRIPE_SECRET") //别名
	viper.AutomaticEnv()
	return viper.ReadInConfig()
}
