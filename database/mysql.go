package database

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"os"
	"strings"
	"time"

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
		"    `device_type` VARCHAR(20) NOT NULL DEFAULT 'ios'," +
		"    `updated_at` BIGINT NOT NULL DEFAULT 0," +
		"    PRIMARY KEY (`id`)," +
		"    UNIQUE KEY `key` (`key`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"

	// migrate old devices table to add the device_type and updated_at columns
	dbSchemaMigration = "" +
		"ALTER TABLE `devices` " +
		"ADD COLUMN `device_type` VARCHAR(20) NOT NULL DEFAULT 'ios', " +
		"ADD COLUMN `updated_at` BIGINT NOT NULL DEFAULT 0"
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

	migrateMySQLDevice(db)

	mysqlDB = db
	return &MySQL{}
}

// migrateMySQLDevice adds the device_type and updated_at columns to a pre-device_type devices table.
func migrateMySQLDevice(db *sql.DB) {
	var count int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'devices' AND COLUMN_NAME = 'device_type'",
	).Scan(&count)
	if err != nil {
		logger.Fatalf("failed to check database schema: %v", err)
	}

	if count == 0 {
		if _, err := db.Exec(dbSchemaMigration); err != nil {
			logger.Fatalf("failed to migrate database schema(%s): %v", dbSchemaMigration, err)
		}
	}
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

func (d *MySQL) DeviceTokenByKey(key string) (*DeviceRecord, error) {
	var token, deviceType string
	var updatedAt int64
	err := mysqlDB.QueryRow("SELECT `token`, `device_type`, `updated_at` FROM `devices` WHERE `key`=?", key).Scan(&token, &deviceType, &updatedAt)
	if err != nil {
		return nil, err
	}

	return &DeviceRecord{Token: token, Type: normalizeDeviceType(deviceType), UpdatedAt: updatedAt}, nil
}

func (d *MySQL) SaveDeviceTokenByKey(key, token, deviceType string) (string, error) {
	if key == "" {
		// Generate a new UUID as the deviceKey when a new device register
		key = shortuuid.New()
	}

	deviceType = normalizeDeviceType(deviceType)
	now := time.Now().Unix()
	_, err := mysqlDB.Exec("INSERT INTO `devices` (`key`,`token`,`device_type`,`updated_at`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE `token`=?, `device_type`=?, `updated_at`=?", key, token, deviceType, now, token, deviceType, now)
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
