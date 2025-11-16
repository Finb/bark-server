package database

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lithammer/shortuuid/v3"
	"github.com/mritd/logger"
)

type MySQL struct {
}

var mysqlDB *sql.DB

const (
	dbSchema = "" +
		"CREATE TABLE IF NOT EXISTS `devices` (" +
		"    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT," +
		"    `key` VARCHAR(255) NOT NULL," +
		"    `token` VARCHAR(255) NOT NULL," +
		"    PRIMARY KEY (`id`)," +
		"    UNIQUE KEY `key` (`key`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"
)

func NewMySQL(dsn string) Database {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatalf("failed to open database connection (%s): %v", dsn, err)
	}

	_, err = db.Exec(dbSchema)
	if err != nil {
		logger.Fatalf("failed to init database schema(%s): %v", dbSchema, err)
	}

	mysqlDB = db
	return &MySQL{}
}

func NewMySQLWithTLS(dsn, tlsName, caPath, certPath, keyPath string, isSkipVerify bool) Database {
	// 1. Load and register TLS configuration
	logger.Infof("MySQL TLS CA: %v", caPath)
	logger.Infof("MySQL TLS client cert: %v", certPath)
	logger.Infof("MySQL TLS client key: %v", keyPath)
	logger.Infof("Server certificate verification skipped: %v", isSkipVerify)
	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile(caPath)
	if err != nil {
		logger.Fatalf("failed to read CA cert: %v", err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		logger.Fatalf("failed to append CA cert")
	}

	var certs []tls.Certificate
	if certPath != "" && keyPath != "" {
		clientCert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			logger.Fatalf("failed to load client cert and key: %v", err)
		}
		certs = []tls.Certificate{clientCert}
	}

	tlsConfig := &tls.Config{
		RootCAs:            rootCertPool,
		Certificates:       certs,
		InsecureSkipVerify: isSkipVerify,
	}

	if err := mysql.RegisterTLSConfig(tlsName, tlsConfig); err != nil {
		logger.Fatalf("failed to register TLS config: %v", err)
	}

	// 2. Append TLS parameter to DSN if missing
	if !strings.Contains(dsn, "tls=") {
		if strings.Contains(dsn, "?") {
			dsn = dsn + "&tls=" + tlsName
		} else {
			dsn = dsn + "?tls=" + tlsName
		}
	}

	// 3. Create and return the Database instance
	return NewMySQL(dsn)
}

func (d *MySQL) CountAll() (int, error) {
	var count int
	err := mysqlDB.QueryRow("SELECT COUNT(1) FROM `devices`").Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (d *MySQL) DeviceTokenByKey(key string) (string, error) {
	var token string
	err := mysqlDB.QueryRow("SELECT `token` FROM `devices` WHERE `key`=? ", key).Scan(&token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (d *MySQL) SaveDeviceTokenByKey(key, token string) (string, error) {
	if key == "" {
		// Generate a new UUID as the deviceKey when a new device register
		key = shortuuid.New()
	}

	_, err := mysqlDB.Exec("INSERT INTO `devices` (`key`,`token`) VALUES (?,?) ON DUPLICATE KEY UPDATE `token`=?", key, token, token)
	if err != nil {
		return "", err
	}

	return key, nil
}

func (d *MySQL) DeleteDeviceByKey(key string) error {
	_, err := mysqlDB.Exec("DELETE FROM `devices` WHERE `key`=?", key)
	return err
}

func (d *MySQL) Close() error {
	return mysqlDB.Close()
}
