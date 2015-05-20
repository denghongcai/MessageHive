package server

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

// TLS设置初始化
func tlsConfig() (tls.Config, error) {
	// 证书链叶子证书
	leafcsr, lerr := ioutil.ReadFile("./certificates/dhc.house.der.csr")
	if lerr != nil {
		return tls.Config{}, lerr
	}
	// 证书链根证书
	rootcsr, rerr := ioutil.ReadFile("./certificates/alphassl.der.csr")
	if rerr != nil {
		return tls.Config{}, rerr
	}
	// 私钥
	privkey, perr := ioutil.ReadFile("./certificates/dhc.house.der.key")
	if perr != nil {
		return tls.Config{}, perr
	}
	priv, pperr := x509.ParsePKCS1PrivateKey(privkey)
	if pperr != nil {
		return tls.Config{}, pperr
	}

	cert := tls.Certificate{
		Certificate: [][]byte{leafcsr, rootcsr},
		PrivateKey:  priv,
	}

	config := tls.Config{
		ClientAuth:   tls.NoClientCert,
		Certificates: []tls.Certificate{cert},
	}
	config.Rand = rand.Reader

	return config, nil
}
