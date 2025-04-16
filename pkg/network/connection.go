package network

import (
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/nice-pink/goutil/pkg/log"
	"golang.org/x/net/proxy"
)

type ConnectionType int

const (
	HttpConnection ConnectionType = iota
	SocketConnection
)

type Connection struct {
	url            string
	port           int
	proxyUrl       string
	proxyPort      int
	Meta           []byte
	SrcInfo        string
	DestInfo       string
	VerboseLogs    bool
	timeout        time.Duration
	connectionType ConnectionType
	httpClient     *http.Client
}

func NewConnection(url, proxyUrl string, port, proxyPort int, timeout time.Duration, connectionType ConnectionType) *Connection {
	// get connection
	c := &Connection{
		url:       url,
		port:      port,
		proxyUrl:  proxyUrl,
		proxyPort: proxyPort,
		timeout:   timeout,
	}

	if connectionType == HttpConnection {
		c.getHttpClient()
	}

	return c
}

func (c Connection) DeepCopy() Connection {
	return Connection{
		url:         c.url,
		port:        c.port,
		proxyUrl:    c.proxyUrl,
		proxyPort:   c.proxyPort,
		Meta:        c.Meta,
		SrcInfo:     c.SrcInfo,
		DestInfo:    c.DestInfo,
		VerboseLogs: c.VerboseLogs,
	}
}

func (c *Connection) getHttpClient() (*http.Client, error) {
	if c.httpClient != nil {
		return c.httpClient, nil
	}

	c.httpClient = &http.Client{Timeout: c.timeout * time.Second}
	if c.proxyUrl != "" && c.proxyPort != 0 {
		// via proxy
		pUrl := "http://" + c.proxyUrl + ":" + strconv.Itoa(c.proxyPort)
		proxyUrl, err := url.Parse(pUrl)
		if err != nil {
			log.Err(err, "proxy url")
			return nil, err
		}
		c.httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	return c.httpClient, nil
}

func (c *Connection) getSocketConn() (net.Conn, error) {
	if c.proxyUrl == "" || c.proxyPort == 0 {
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

func (c *Connection) GetAddr() string {
	return c.url + ":" + strconv.Itoa(c.port)
}

func (c *Connection) GetProxyAddr() string {
	return c.proxyUrl + ":" + strconv.Itoa(c.proxyPort)
}

func (c *Connection) RemoveUrlProtocol() string {
	reg := regexp.MustCompile("(.*://)(.*)")
	return reg.ReplaceAllString(c.url, "${2}")
}
