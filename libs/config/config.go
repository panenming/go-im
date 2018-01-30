package config

import (
	"errors"
	"io/ioutil"
	"path"
	"strings"

	"github.com/panenming/go-im/libs/config/ini"
	"github.com/panenming/go-im/libs/config/json"
	"github.com/panenming/go-im/libs/config/toml"
	"github.com/panenming/go-im/libs/config/xml"
	"github.com/panenming/go-im/libs/config/yaml"
)

func New(f string, c interface{}) error {
	buf, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	fileSuffix := getFileSuffix(f)
	switch fileSuffix {
	case "json":
		err = json.Unmarshal(buf, c)
	case "yaml":
		err = yaml.Unmarshal(buf, c)
	case "xml":
		err = xml.Unmarshal(buf, c)
	case "ini":
		err = ini.Unmarshal(buf, c)
	case "toml":
		err = toml.Unmarshal(buf, c)
	default:
		err = errors.New("格式不支持")
	}
	return err
}

func getFileSuffix(f string) string {
	filePath := path.Base(f)
	fileSuffix := path.Ext(filePath)
	return strings.Trim(fileSuffix, ".")
}
