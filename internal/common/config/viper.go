package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func init() {
	if err := NewViperConfig(); err != nil {
		panic(err)
	}
}

var once sync.Once

func NewViperConfig() (err error) {
	once.Do(func() {
		err = newViperConfig()
	})
	return
}

func newViperConfig() error {
	relPath, err := getRelativePathFromCaller()
	if err != nil {
		return err
	}
	viper.SetConfigName("global")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(relPath)
	viper.EnvKeyReplacer(strings.NewReplacer("-", "_"))
	_ = viper.BindEnv("stripe-key", "STRIPE_KEY", "endpoint-stripe-secret", "ENDPOINT_STRIPE_SECRET") //别名
	viper.AutomaticEnv()
	return viper.ReadInConfig()
}

// 这样保证无论从项目的哪个目录启动程序，都能正确找到配置文件
func getRelativePathFromCaller() (relPath string, err error) {
	// 返回程序启动时的工作目录，哪里调了viper，就是哪里
	callerPwd, err := os.Getwd()
	if err != nil {
		return
	}
	// here 即“XX/gorder/internal/common/config/viper.go”
	_, here, _, _ := runtime.Caller(0)
	// 计算从callerPwd到here目录的相对路径
	relPath, err = filepath.Rel(callerPwd, filepath.Dir(here)) //Dir是获取当前文件所在目录
	fmt.Printf("caller from: %s, here: %s, relpath: %s", callerPwd, here, relPath)
	return
}
