package parse

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Variables struct {
		Node         string `yaml:"node"`
		Job          string `yaml:"job"`
		Namespace    string `yaml:"namespace"`
		Interval     string `yaml:"interval"`
		RateInterval string `yaml:"rate_interval"`
	}
}

func ParseConfig() *Config {
	path, err := filepath.Abs("../../conf/conf.yaml")
	if err != nil {
		log.Fatal(err)
	}
	conf, err := ReadConfig(path)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}

func ReadConfig(filename string) (*Config, error) {
	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	conf := &Config{}
	err = yaml.Unmarshal(buffer, conf)
	if err != nil {
		log.Fatal(err)
	}
	return conf, err
}

func ParseEpoch(timestamp interface{}) time.Time {
	flt, _, err := big.ParseFloat(fmt.Sprint("", timestamp), 10, 0, big.ToNearestEven)
	if err != nil {
		log.Fatal(err)

	}
	i, _ := flt.Int64()
	return time.Unix(i, 0)
}

func ParseQuery(query string, conf *Config) (string, error) {
	if len(query) > 0 {
		if strings.Contains(query, "$node") || strings.Contains(query, "$job") || strings.Contains(query, "$namespace") || strings.Contains(query, "$interval") || strings.Contains(query, "$__rate_interval") {
			replacer := strings.NewReplacer("$node", conf.Variables.Node, "$job", conf.Variables.Job, "$namespace", conf.Variables.Namespace, "$interval", conf.Variables.Interval, "$__rate_interval", conf.Variables.RateInterval)
			resquery := replacer.Replace(query)
			return resquery, nil
		}
	} else {
		return "", errors.New("no query to parse")
	}

	return query, nil
}
