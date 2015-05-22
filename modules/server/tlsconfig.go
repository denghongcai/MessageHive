package server

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/spf13/viper"
)

// TLS设置初始化
func tlsConfig() (tls.Config, error) {
	certPaths := viper.GetStringSlice("tls.certificates")
	certificates := make([][]byte, 0)
	for _, value := range certPaths {
		csr, err := ioutil.ReadFile(value)
		if err != nil {
			return tls.Config{}, err
		}
		certificates = append(certificates, csr)
	}
	/*
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
	*/
	// 私钥
	privkey, err := ioutil.ReadFile(viper.GetString("tls.key"))
	if err != nil {
		return tls.Config{}, err
	}
	priv, err := x509.ParsePKCS1PrivateKey(privkey)
	if err != nil {
		return tls.Config{}, err
	}

	cert := tls.Certificate{
		Certificate: certificates,
		PrivateKey:  priv,
	}

	config := tls.Config{
		ClientAuth:   tls.NoClientCert,
		Certificates: []tls.Certificate{cert},
	}
	config.Rand = rand.Reader

	return config, nil
}
