package config

import (
	"io"
	"io/ioutil"
	"reflect"
	"regexp"
	"sync"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type App struct {
	Name      string `yaml:"name"`
	Secret    string `yaml:"secret"`
	Version   string `yaml:"version"`
	Env       string `yaml:"env"`
	Port      int    `yaml:"port"`
	PortRange struct {
		Min int `yaml:"min"`
		Max int `yaml:"max"`
	} `yaml:"port_range"`
}

type Logger struct {
	Filename string `yaml:"filename"`
	LogLevel int    `yaml:"log_level"`
}

type Sentry struct {
	Dsn string `yaml:"dsn"`
}

type Consul struct {
	Addr string `yaml:"addr"`
}

type Redis struct {
	Addr string `yaml:"addr"`
	Auth string `yaml:"auth"`
	Size int    `yaml:"size"`
}

type Memcache struct {
	Addr []string `yaml:"addr"`
}

type Mysql struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Dbname       string `yaml:"dbname"`
	Charset      string `yaml:"charset"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type NSQ struct {
	Nsqlookupd []string `yaml:"nsqlookupd"`
	Nsqd       string   `yaml:"nsqd"`
	Topic      string   `yaml:"topic"`
	Channel    string   `yaml:"channel"`
	Size       int      `yaml:"size"`
}

type Influxdb struct {
	Addr      string `yaml:"addr"`
	Database  string `yaml:"database"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Precision string `yaml:"precision"`
	Size      int    `yaml:"size"`
	BatchSize int    `yaml:"batch_size"`
}

type Alilog struct {
	Endpoint      string        `yaml:"endpoint"`
	AccessKey     string        `yaml:"access_key"`
	AccessSecret  string        `yaml:"access_secret"`
	Project       string        `yaml:"project"`
	LogStore      string        `yaml:"log_store"`
	Topic         string        `yaml:"topic"`
	Size          int           `yaml:"size"`
	UploadTimeout time.Duration `yaml:"upload_timeout"`
}

type Tablestore struct {
	Endpoint        string `yaml:"endpoint"`
	InstanceName    string `yaml:"instance_name"`
	AccessKeyId     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
}

var once sync.Once
var _conf interface{}

func CreateConfig(reader io.Reader, conf interface{}) error {
	var err error
	once.Do(func() {
		content, err := ioutil.ReadAll(reader)
		if err != nil {
			return
		}
		err = yaml.Unmarshal(content, conf)
		_conf = conf
	})
	return err
}

func GetConfig() interface{} {
	return _conf
}

func ReplaceWithMark(c interface{}, secretMap map[string]string) {
	if reflect.TypeOf(c).Kind() != reflect.Ptr {
		panic("need pointer")
	}
	v := reflect.ValueOf(c).Elem()
	updateStructValueWithSecretMap(v, secretMap)
}

func RenderWithData(src []byte, data map[string]string) []byte {
	pattern := regexp.MustCompile("{{(.*?)}}")
	dst := pattern.ReplaceAllFunc(src, func(matched []byte) []byte {
		val, ok := data[string(matched[2:len(matched)-2])]
		if ok {
			return []byte(val)
		}
		return matched
	})
	return dst
}

func updateStructValueWithSecretMap(v reflect.Value, secretMap map[string]string) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		switch t.Field(i).Type.Kind() {
		case reflect.Struct:
			updateStructValueWithSecretMap(v.Field(i), secretMap)
		case reflect.Slice, reflect.Array:
			switch v.Field(i).Interface().(type) {
			case []string:
				l := v.Field(i).Len()
				for j := 0; j < l; j++ {
					val := v.Field(i).Index(j).String()
					if len(val) > 4 && val[0:2] == "{{" && val[len(val)-2:] == "}}" {
						secret, ok := secretMap[val[2:len(val)-2]]
						if ok {
							v.Field(i).Index(j).SetString(secret)
						}
					}
				}
			}
		case reflect.String:
			val := v.Field(i).String()
			if len(val) > 4 && val[0:2] == "{{" && val[len(val)-2:] == "}}" {
				secret, ok := secretMap[val[2:len(val)-2]]
				if ok {
					v.Field(i).SetString(secret)
				}
			}
		}
	}
}
