package configfile

import (
	"fmt"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

func TestEncodeAuth(t *testing.T) {
	newAuthConfig := &types.AuthConfig{Username: "ken", Password: "test"}
	authStr := encodeAuth(newAuthConfig)
	decAuthConfig := &types.AuthConfig{}
	var err error
	decAuthConfig.Username, decAuthConfig.Password, err = decodeAuth(authStr)
	if err != nil {
		t.Fatal(err)
	}
	if newAuthConfig.Username != decAuthConfig.Username {
		t.Fatal("Encode Username doesn't match decoded Username")
	}
	if newAuthConfig.Password != decAuthConfig.Password {
		t.Fatal("Encode Password doesn't match decoded Password")
	}
	if authStr != "a2VuOnRlc3Q=" {
		t.Fatal("AuthString encoding isn't correct.")
	}
}

func TestProxyConfig(t *testing.T) {
	httpProxy := "http://proxy.mycorp.com:3128"
	httpsProxy := "https://user:password@proxy.mycorp.com:3129"
	ftpProxy := "http://ftpproxy.mycorp.com:21"
	noProxy := "*.intra.mycorp.com"
	defaultProxyConfig := ProxyConfig{
		HTTPProxy:  httpProxy,
		HTTPSProxy: httpsProxy,
		FTPProxy:   ftpProxy,
		NoProxy:    noProxy,
	}

	cfg := ConfigFile{
		Proxies: map[string]ProxyConfig{
			"default": defaultProxyConfig,
		},
	}

	proxyConfig := cfg.ParseProxyConfig("/var/run/docker.sock", []string{})
	expected := map[string]*string{
		"HTTP_PROXY":  &httpProxy,
		"http_proxy":  &httpProxy,
		"HTTPS_PROXY": &httpsProxy,
		"https_proxy": &httpsProxy,
		"FTP_PROXY":   &ftpProxy,
		"ftp_proxy":   &ftpProxy,
		"NO_PROXY":    &noProxy,
		"no_proxy":    &noProxy,
	}
	assert.Equal(t, expected, proxyConfig)
}

func TestProxyConfigOverride(t *testing.T) {
	httpProxy := "http://proxy.mycorp.com:3128"
	overrideHTTPProxy := "http://proxy.example.com:3128"
	overrideNoProxy := ""
	httpsProxy := "https://user:password@proxy.mycorp.com:3129"
	ftpProxy := "http://ftpproxy.mycorp.com:21"
	noProxy := "*.intra.mycorp.com"
	defaultProxyConfig := ProxyConfig{
		HTTPProxy:  httpProxy,
		HTTPSProxy: httpsProxy,
		FTPProxy:   ftpProxy,
		NoProxy:    noProxy,
	}

	cfg := ConfigFile{
		Proxies: map[string]ProxyConfig{
			"default": defaultProxyConfig,
		},
	}

	ropts := []string{
		fmt.Sprintf("HTTP_PROXY=%s", overrideHTTPProxy),
		"NO_PROXY=",
	}
	proxyConfig := cfg.ParseProxyConfig("/var/run/docker.sock", ropts)
	expected := map[string]*string{
		"HTTP_PROXY":  &overrideHTTPProxy,
		"http_proxy":  &httpProxy,
		"HTTPS_PROXY": &httpsProxy,
		"https_proxy": &httpsProxy,
		"FTP_PROXY":   &ftpProxy,
		"ftp_proxy":   &ftpProxy,
		"NO_PROXY":    &overrideNoProxy,
		"no_proxy":    &noProxy,
	}
	assert.Equal(t, expected, proxyConfig)
}

func TestProxyConfigPerHost(t *testing.T) {
	httpProxy := "http://proxy.mycorp.com:3128"
	httpsProxy := "https://user:password@proxy.mycorp.com:3129"
	ftpProxy := "http://ftpproxy.mycorp.com:21"
	noProxy := "*.intra.mycorp.com"

	extHTTPProxy := "http://proxy.example.com:3128"
	extHTTPSProxy := "https://user:password@proxy.example.com:3129"
	extFTPProxy := "http://ftpproxy.example.com:21"
	extNoProxy := "*.intra.example.com"

	defaultProxyConfig := ProxyConfig{
		HTTPProxy:  httpProxy,
		HTTPSProxy: httpsProxy,
		FTPProxy:   ftpProxy,
		NoProxy:    noProxy,
	}
	externalProxyConfig := ProxyConfig{
		HTTPProxy:  extHTTPProxy,
		HTTPSProxy: extHTTPSProxy,
		FTPProxy:   extFTPProxy,
		NoProxy:    extNoProxy,
	}

	cfg := ConfigFile{
		Proxies: map[string]ProxyConfig{
			"default":                       defaultProxyConfig,
			"tcp://example.docker.com:2376": externalProxyConfig,
		},
	}

	proxyConfig := cfg.ParseProxyConfig("tcp://example.docker.com:2376", []string{})
	expected := map[string]*string{
		"HTTP_PROXY":  &extHTTPProxy,
		"http_proxy":  &extHTTPProxy,
		"HTTPS_PROXY": &extHTTPSProxy,
		"https_proxy": &extHTTPSProxy,
		"FTP_PROXY":   &extFTPProxy,
		"ftp_proxy":   &extFTPProxy,
		"NO_PROXY":    &extNoProxy,
		"no_proxy":    &extNoProxy,
	}
	assert.Equal(t, expected, proxyConfig)
}

func TestConfigFile(t *testing.T) {
	configFilename := "configFilename"
	configFile := NewConfigFile(configFilename)

	assert.Equal(t, configFilename, configFile.Filename)
}
