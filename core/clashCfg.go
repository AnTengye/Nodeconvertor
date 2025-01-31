package core

import (
	"fmt"
	"github.com/AnTengye/NodeConvertor/lib/yemoji"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"strings"
)

const (
	ClashProxies      = "proxies"
	ClashProxiesGroup = "proxy-groups"
	ClashRules        = "rules"
)

type baseClash struct {
	ProxyGroups []ProxyGroup `yaml:"proxy-groups"`
	//Rules       []string               `yaml:"rules"`
}
type ProxyGroup struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	// 其他字段
	OtherFields map[string]interface{} `yaml:",inline"`
	Proxies     []string               `yaml:"proxies"`
}

type Clash struct {
	data    map[string]any
	Proxies []Node
	base    *baseClash
}

func NewClash(filePath string) (*Clash, error) {
	var data map[string]any
	yamlData, err := yamlFromFile(filePath)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(yamlData, &data); err != nil {
		return nil, err
	}
	var base baseClash
	if err = yaml.Unmarshal(yamlData, &base); err != nil {
		return nil, err
	}
	return &Clash{
		data:    data,
		Proxies: make([]Node, 0),
		base:    &base,
	}, nil
}

func (c *Clash) AddProxy(n ...Node) {
	c.Proxies = append(c.Proxies, n...)
}

func (c *Clash) ToYaml() (string, error) {
	c.data[ClashProxies] = c.Proxies
	c.withProxyGroup()
	d, err := yaml.Marshal(c.data)
	if err != nil {
		zap.S().Errorw("ToYaml error", "err", err)
		return "", err
	}
	unicodePoints, err := yemoji.ParseUnicodePoints(d)
	if err != nil {
		return "", fmt.Errorf("generate clash with unicode error: %v", err)
	}
	return string(unicodePoints), nil
}

func (c *Clash) withProxyGroup() {
	proxies := make([]string, len(c.Proxies))
	for i, v := range c.Proxies {
		proxies[i] = v.Name()
	}
	for i, v := range c.base.ProxyGroups {
		for _, groupName := range viper.GetStringSlice("Advanced.ProxyGroup") {
			if strings.Contains(v.Name, groupName) {
				c.base.ProxyGroups[i].Proxies = proxies
			}
		}
	}
	c.data[ClashProxiesGroup] = c.base.ProxyGroups
}
