package network

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nice-pink/audio-tool/pkg/util"
	"github.com/nice-pink/goutil/pkg/log"
	"golang.org/x/net/proxy"
)

type ConnectionType int

const (
	HttpConnection ConnectionType = iota
	SocketConnection
)

type Connection struct {
	url               string
	port              int
	proxyUrl          string
	proxyPort         int
	Meta              []byte
	SrcInfo           string
	DestInfo          string
	VerboseLogs       bool
	timeout           time.Duration
	connectionType    ConnectionType
	httpClient        *http.Client
	socketConn        net.Conn
	isSocketConnected bool
	metrics           *Metrics
}

func NewConnection(url, proxyUrl string, port, proxyPort int, timeout time.Duration, connectionType ConnectionType, metrics util.MetricsControl) *Connection {
	// get connection
	c := &Connection{
		url:       url,
		port:      port,
		proxyUrl:  proxyUrl,
		proxyPort: proxyPort,
		timeout:   timeout,
	}

	// init http connection
	if connectionType == HttpConnection {
		c.getHttpClient()
	}

	// init metrics
	if metrics.Enabled {
		c.metrics = NewMetrics(metrics.Prefix, metrics.Labels)
	}

	return c
}

func (c *Connection) Close() {
	if c.socketConn != nil {
		c.CloseSocket()
	}
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
			log.Err(err, "proxy url", pUrl)
			return nil, err
		}
		log.Info("Use proxy", proxyUrl)
		c.httpClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	return c.httpClient, nil
}

func (c *Connection) GetSocketConn() (net.Conn, error) {
	if c.isSocketConnected {
		// return socket
		return c.socketConn, nil
	}

	// init new socket connection
	addr := c.GetAddr()
	var err error

	if strings.HasPrefix(c.url, "https://") {
		// is tls connection (https)
		log.Info("Use tls - no proxy!")
		c.socketConn, err = tls.Dial(TCP_PROTO, addr, &tls.Config{
			InsecureSkipVerify: true,
		})
		if err != nil {
			c.isSocketConnected = true
		}
		return c.socketConn, err
	}

	if c.proxyUrl == "" || c.proxyPort == 0 {
		// no proxy
		c.socketConn, err = net.Dial(TCP_PROTO, addr)
		if err != nil {
			c.isSocketConnected = true
		}
		return c.socketConn, err
	}

	// setup proxy
	proxyAddr := c.GetProxyAddr()
	dialer, err := proxy.SOCKS5(TCP_PROTO, proxyAddr, nil, proxy.Direct)
	if err != nil {
		log.Err(err, "could not create socks5 proxy", proxyAddr)
		return nil, err
	}

	if c.VerboseLogs {
		log.Info("Use proxy", proxyAddr, "to connect socket to", addr)
	}
	c.socketConn, err = dialer.Dial(TCP_PROTO, addr)
	if err != nil {
		c.isSocketConnected = true
	}
	return c.socketConn, err
}

func (c *Connection) CloseSocket() {
	c.isSocketConnected = false
	c.socketConn.Close()
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
