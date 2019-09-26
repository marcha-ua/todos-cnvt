package owl

import (
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
	"todos-cnvt/todos"
)

func ReadConfigEdge(cfg string)todos.EdgeTagName{

	viper.SetDefault("APP_PORT", "8080")
ext:="toml"
file:="config"
dir:="."
	if len(cfg)!=0{
		dir = path.Dir(cfg)
		ext=path.Ext(cfg)
		file=path.Base(cfg)
		file = strings.Trim(file, ext)
		ext = strings.TrimLeft(ext,".")
	}

	viper.SetDefault("subclass", "SubClass")
	viper.SetDefault("domain", "Domain")
	viper.SetDefault("individual", "Individual")
	viper.SetDefault("property", "Property")
	viper.SetDefault("default", "Default")
	viper.SetDefault("collection", "Collection")

	if os.Getenv("ENVIRONMENT") != "EDGE_NAME" {
		viper.SetConfigName(file)
		viper.SetConfigType(ext)
		viper.AddConfigPath(dir)
		viper.ReadInConfig()
	} else {  viper.AutomaticEnv() }

edge:=todos.EdgeTagName{
	SubClassGroup:viper.GetString("subclass"),
	DomainGroup:viper.GetString("domain"),
	IndividualGroup :viper.GetString("individual"),
	PropertyGroup:viper.GetString("property"),
	DefaultGroup    :viper.GetString("default"),
	CollectionGroup :viper.GetString("collection"),
}
return edge
}
