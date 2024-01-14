package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func ParseDatabaseURI(uri string) (string, error) {
	parsedURL, err := url.Parse(uri)
	ErrorPanicPrinter(err, true)

	//enchantech-cluster-do-user-12948347-0.c.db.ondigitalocean.com
	//enchantech-codex-cluster-tf-do-user-12948347-0.c.db.ondigitalocean.com

	username := parsedURL.User.Username()
	password, _ := parsedURL.User.Password()

	host := parsedURL.Hostname()
	port := parsedURL.Port()

	database := strings.TrimPrefix(parsedURL.Path, "/")

	newFormat := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=skip-verify", username, password, host, port, database)
	println("newFormat: " + newFormat)

	return newFormat, err

	//connectionString := "mysql://doadmin:AVNS_z0OORLs1Abp8Vbun_0u@enchantech-cluster-do-user-12948347-0.c.db.ondigitalocean.com:25060/enchantech-db?ssl-mode=REQUIRED"
	//doadmin:AVNS_z0OORLs1Abp8Vbun_0u@tcp(enchantech-cluster-do-user-12948347-0.c.db.ondigitalocean.com:25060)/enchantech-db?tls=skip-verify
	//doadmin:AVNS_mquevJ0IZghOFg3PHXP@tcp(enchantech-codex-cluster-tf-do-user-12948347-0.c.db.ondigitalocean.com:25060)/defaultdb?tls=skip-verify
}
