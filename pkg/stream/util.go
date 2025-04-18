package stream

import (
	"encoding/base64"
	"errors"
	"net/url"
	"strings"

	"github.com/nice-pink/goutil/pkg/log"
)

type ConnTarget struct {
	Scheme     string
	Domain     string
	MountPoint string
	Port       string
	BasicAuth  string
	UserAgent  string
}

func GetConnTarget(fullUrl string) (ConnTarget, error) {
	full, err := url.Parse(fullUrl)
	if err != nil {
		return ConnTarget{}, err
	}

	domain := full.Hostname()
	if domain == "" {
		return ConnTarget{}, errors.New("invalid url")
	}

	scheme := full.Scheme
	mountPoint := full.Path
	port := full.Port()
	if port == "" {
		if strings.HasPrefix(fullUrl, "https://") {
			port = "443"
		} else {
			port = "80"
		}
	}
	password, hasPassword := full.User.Password()
	basicAuth := ""
	if hasPassword {
		basicAuth = getBasicAuth(full.User.Username(), password)
	}
	return ConnTarget{Scheme: scheme, Domain: domain, MountPoint: mountPoint, Port: port, BasicAuth: basicAuth}, nil
}

func getBasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func convertToByteHeader(header string, print bool) []byte {
	if print {
		log.Info("Header:\n" + header)
		log.Info("Header size:", len(header))
	}
	return []byte(header)
}
