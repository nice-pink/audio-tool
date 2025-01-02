package network

import (
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/nice-pink/goutil/pkg/log"
	"golang.org/x/net/proxy"
)

type Connection struct {
	Url         string
	Port        int
	ProxyUrl    string
	ProxyPort   int
	Meta        []byte
	SrcInfo     string
	DestInfo    string
	VerboseLogs bool
}

func NewConnection(url string, port int) Connection {
	return Connection{Url: url, Port: port}
}

func (c Connection) DeepCopy() Connection {
	return Connection{
		Url:         c.Url,
		Port:        c.Port,
		ProxyUrl:    c.ProxyUrl,
		ProxyPort:   c.ProxyPort,
		Meta:        c.Meta,
		SrcInfo:     c.SrcInfo,
		DestInfo:    c.DestInfo,
		VerboseLogs: c.VerboseLogs,
	}
}

func (c Connection) GetHttpClient(timeout time.Duration) (*http.Client, error) {
	client := &http.Client{Timeout: timeout * time.Second}

	if c.ProxyUrl != "" && c.ProxyPort != 0 {
		// via proxy
		proxyUrl, err := url.Parse("http://" + c.GetProxyAddr())
		if err != nil {
			log.Err(err, "proxy url")
			return nil, err
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	return client, nil
}

func (c Connection) GetSocketConn() (net.Conn, error) {
	if c.ProxyUrl == "" || c.ProxyPort == 0 {
		addr := c.GetAddr()
		return net.Dial(TCP_PROTO, addr)
	}

	proxyAddr := c.GetProxyAddr()
	dialer, err := proxy.SOCKS5(TCP_PROTO, proxyAddr, nil, proxy.Direct)
	if err != nil {
		log.Err(err, "could not create socks5 proxy")
		return nil, err
	}

	addr := c.GetAddr()
	if c.VerboseLogs {
		log.Info("Use proxy", proxyAddr, "to connect socket to", addr)
	}

	return dialer.Dial(TCP_PROTO, addr)
}

func (c Connection) GetAddr() string {
	return c.Url + ":" + strconv.Itoa(c.Port)
}

func (c Connection) GetProxyAddr() string {
	return c.ProxyUrl + ":" + strconv.Itoa(c.ProxyPort)
}
