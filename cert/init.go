package cert

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/go-sql-driver/mysql"
)

func Init() {
	caCert, err := os.ReadFile("./cert/ca-certificate.crt")
	if err != nil {
		panic(err)
	}

	caPool := x509.NewCertPool()
	if ok := caPool.AppendCertsFromPEM(caCert); !ok {
		panic("failed to append CA cert")
	}

	// 注册一个自定义 TLS 配置
	err = mysql.RegisterTLSConfig("do", &tls.Config{
		RootCAs:    caPool,
		ServerName: "db-mysql-sgp1-11646-do-user-6185766-0.m.db.ondigitalocean.com",
	})
	if err != nil {
		panic(err)
	}
}
