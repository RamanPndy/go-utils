package goutils

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/dlmiddlecote/sqlstats"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/sirupsen/logrus"
)

const (
	HOTLOADDBType = "hotload"
)

type DatabaseConfig struct {
	Type            string
	Address         string
	Port            int
	User            string
	Password        string
	Name            string
	SSL             string
	DSN             string
	LogMode         bool
	MaxOpenConns    int
	MaxIdleConns    int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

type DBConnInterface interface {
	Connect() (*gorm.DB, error)
	Close(*gorm.DB) error
	Ready() error
}

type DBConn struct {
	dbConfig *DatabaseConfig
	logger   *logrus.Logger
}

func NewDBConn(dbConfig *DatabaseConfig, logger *logrus.Logger) (*DBConn, error) {
	return &DBConn{
		dbConfig: dbConfig,
		logger:   logger,
	}, nil
}

func (c *DBConn) Connect() (*gorm.DB, error) {
	db, err := new(c.dbConfig, c.logger)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (c *DBConn) Close(db *gorm.DB) error {
	if db != nil {
		if err := db.Close(); err != nil {
			c.logger.WithError(err).Error("failed to close database connection")
			return err
		}
	}
	return nil
}

func new(dbConfig *DatabaseConfig, logger *logrus.Logger) (*gorm.DB, error) {
	dbType := strings.TrimSpace(dbConfig.Type)

	if strings.TrimSpace(dbConfig.DSN) == "" {
		logger.Error("database.dsn is required")
		return nil, fmt.Errorf("database.dsn is required")
	}

	var (
		db  *gorm.DB
		err error
	)
	if dbType == HOTLOADDBType {
		if err := setupGormWithHotload(dbConfig); err != nil {
			return nil, err
		}

	}
	db, err = gorm.Open(dbType, dbConfig.DSN)
	if err != nil {
		logger.WithError(err).Error("failed to open database connection")
		return nil, err
	}
	configureDBConns(db, dbConfig)
	registerSQLStats(logger, db)
	return db, nil
}

func (c *DBConn) Ready() error {
	db, err := sql.Open(c.dbConfig.Type, c.dbConfig.DSN)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Ping()
}

// registerSQLStats exposes database/sql pool stats as go_sql_stats_connections_* metrics.
// It is idempotent: a duplicate registration is silently ignored.
func registerSQLStats(logger *logrus.Logger, db *gorm.DB) {
	collector := sqlstats.NewStatsCollector("axur_integrator", db.DB())
	if err := prometheus.Register(collector); err != nil {
		if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
			logger.WithError(err).Warn("failed to register sqlstats collector")
		}
	}
}

// configureDBConns sets the db connection string
func configureDBConns(db *gorm.DB, dbConfig *DatabaseConfig) {
	db.LogMode(dbConfig.LogMode)
	db.DB().SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.DB().SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.DB().SetConnMaxLifetime(dbConfig.MaxConnLifetime)
	db.DB().SetConnMaxIdleTime(dbConfig.MaxConnIdleTime)
}

func setupGormWithHotload(dbConfig *DatabaseConfig) error {
	if dbType := dbConfig.Type; dbType == HOTLOADDBType {
		u, err := url.Parse(dbConfig.DSN)
		if err != nil {
			return fmt.Errorf("could not parse hotload dsn, it must be an URL in case if hotload is enabled: %w", err)
		}
		// hostname in hotload URLs contains the name of the underlying driver
		dialect, ok := gorm.GetDialect(u.Hostname())
		if !ok {
			return fmt.Errorf("could not find underlying dialect %v", u.Hostname())
		}
		gorm.RegisterDialect(dbType, dialect)
	}
	return nil
}

func BuildDSN(dbConfig *DatabaseConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s dbname=%s",
		dbConfig.Address, dbConfig.Port,
		dbConfig.User, dbConfig.Password,
		dbConfig.SSL, dbConfig.Name)
}

func OpenDB(dbConfig *DatabaseConfig, logger *logrus.Logger) (*gorm.DB, error) {
	// For hotload, register the dialect before gorm.Open (same as grpcserver.setupGormWithHotload).
	var (
		dbType string
		dsn    string
	)

	if dbConfig.Type == HOTLOADDBType {
		if err := registerHotloadDialect(dbConfig.DSN); err != nil {
			return nil, err
		}
	} else {
		dsn = BuildDSN(dbConfig)
	}

	database, err := gorm.Open(dbType, dsn)
	if err != nil {
		return nil, fmt.Errorf("openDB: %w", err)
	}
	database.SetLogger(logger)
	configureDBConns(database, dbConfig)
	registerSQLStats(logger, database)

	return database, nil
}

// registerHotloadDialect registers the hotload dialect with gorm. This must be
// called before gorm.Open when database.type is "hotload". The DSN is expected
// to be a URL like "hotload://postgres/..." where the hostname is the underlying
// driver name.
func registerHotloadDialect(dsn string) error {
	u, err := url.Parse(dsn)
	if err != nil {
		return fmt.Errorf("could not parse hotload dsn: %w", err)
	}
	dialect, ok := gorm.GetDialect(u.Hostname())
	if !ok {
		return fmt.Errorf("could not find underlying dialect %v", u.Hostname())
	}
	gorm.RegisterDialect("hotload", dialect)
	return nil
}
