package config

import (
	"fmt"
	"github.com/cloudentity/cac/internal/cac/client"
	"github.com/cloudentity/cac/internal/cac/logging"
	"github.com/cloudentity/cac/internal/cac/storage"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"net/url"
	"reflect"
	"strings"
	"time"
)

var (
	DefaultConfig = Configuration{
		Client:  client.DefaultConfig,
		Logging: logging.DefaultLoggingConfig,
		Storage: storage.DefaultMultiStorageConfig,
	}
)

type Configuration struct {
	Client  client.Configuration              `json:"client"`
	Logging logging.Configuration             `json:"logging"`
	Storage storage.MultiStorageConfiguration `json:"storage"`
}

func InitConfig(path string) (_ *Configuration, err error) {
	var (
		decoder    *mapstructure.Decoder
		decodedMap map[string]any
		config     = DefaultConfig
		dconf      = mapstructure.DecoderConfig{
			Result: &decodedMap,
		}
	)

	configureDecoder(&dconf)
	v := viper.GetViper()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if decoder, err = mapstructure.NewDecoder(&dconf); err != nil {
		return nil, err
	}

	if err = decoder.Decode(config); err != nil {
		return nil, err
	}

	for k, val := range decodedMap {
		v.SetDefault(k, val)
	}

	if path != "" {
		v.SetConfigFile(path)

		if err := v.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	if err := v.Unmarshal(&config, configureDecoder); err != nil {
		return nil, err
	}

	slog.With("config", config).Debug("Initiated configuration")

	return &config, nil
}

func configureDecoder(config *mapstructure.DecoderConfig) {
	config.TagName = "json"
	config.WeaklyTypedInput = true
	config.DecodeHook = mapstructure.ComposeDecodeHookFunc(urlDecoder(), timeDecoder())
}

func urlDecoder() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(url.URL{}) {
			return data, nil
		}

		var (
			str string
			err error
			ok  bool
			u   *url.URL
		)

		if str, ok = data.(string); !ok {
			return nil, fmt.Errorf("cannot map %v", reflect.TypeOf(data))
		}

		if u, err = url.Parse(str); err != nil {
			return nil, err
		}

		return u, nil
	}
}

func timeDecoder() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
	}
}
