package goutils

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/dlmiddlecote/sqlstats"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/sirupsen/logrus"
)

const (
	HOTLOADDBType = "hotload"
	POSTGRESQL    = "postgres"
	MYSQL         = "mysql"
	SQLITE3       = "sqlite3"
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
	Migrate(db *gorm.DB, models ...interface{}) error
	Rollback(db *gorm.DB) error
	MigrateFromPath(migrationsPath string) error
	RollbackFromPath(migrationsPath string, steps int) error
	Version() (string, error)
	GetDBConfig() *DatabaseConfig
	IsHotload() bool
	GetDSN() string
	GetDBType() string
	IsReady() bool
	IsDBSupported() bool
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

func (c *DBConn) Migrate(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...).Error
}

func (c *DBConn) Rollback(db *gorm.DB) error {
	return db.Close()
}

func (c *DBConn) MigrateFromPath(migrationsPath string) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		c.dbConfig.DSN,
	)
	if err != nil {
		c.logger.Errorf("Run : initializing migrator: %v", err)
		return fmt.Errorf("initializing migrator: %w", err)
	}
	defer func() {
		_, err = m.Close()
		if err != nil {
			c.logger.Errorf("closing migrator: %v", err)
		}
	}()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		c.logger.Errorf("running migrations: %v", err)
		return fmt.Errorf("running migrations: %w", err)
	}
	return nil
}

func (c *DBConn) RollbackFromPath(migrationsPath string, steps int) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		c.dbConfig.DSN,
	)
	if err != nil {
		c.logger.Errorf("Run : initializing migrator: %v", err)
		return fmt.Errorf("initializing migrator: %w", err)
	}
	defer func() {
		_, err = m.Close()
		if err != nil {
			c.logger.Errorf("closing migrator: %v", err)
		}
	}()

	if steps != 0 {
		if err := m.Steps(-steps); err != nil && err != migrate.ErrNoChange {
			c.logger.Errorf("running rollback migrations: %v", err)
			return fmt.Errorf("running rollback migrations: %w", err)
		}
	} else {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			c.logger.Errorf("running rollback migrations: %v", err)
			return fmt.Errorf("running rollback migrations: %w", err)
		}
	}
	return nil
}

func (c *DBConn) Version() (string, error) {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", "migrations"),
		c.dbConfig.DSN,
	)
	if err != nil {
		c.logger.Errorf("Run : initializing migrator: %v", err)
		return "", fmt.Errorf("initializing migrator: %w", err)
	}
	defer func() {
		_, err = m.Close()
		if err != nil {
			c.logger.Errorf("closing migrator: %v", err)
		}
	}()

	version, dirty, err := m.Version()
	if err != nil {
		c.logger.Errorf("getting migration version: %v", err)
		return "", fmt.Errorf("getting migration version: %w", err)
	}

	if dirty {
		return fmt.Sprintf("%d (dirty)", version), nil
	}
	return fmt.Sprintf("%d", version), nil
}

func (c *DBConn) GetDBConfig() *DatabaseConfig {
	return c.dbConfig
}

func (c *DBConn) IsHotload() bool {
	return c.dbConfig.Type == HOTLOADDBType
}

func (c *DBConn) GetDSN() string {
	return c.dbConfig.DSN
}

func (c *DBConn) GetDBType() string {
	return c.dbConfig.Type
}

func (c *DBConn) IsReady() bool {
	db, err := sql.Open(c.dbConfig.Type, c.dbConfig.DSN)
	if err != nil {
		c.logger.WithError(err).Error("failed to open database connection")
		return false
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		c.logger.WithError(err).Error("database ping failed")
		return false
	}
	return true
}

func (c *DBConn) IsDBSupported() bool {
	supportedDBs := []string{POSTGRESQL, MYSQL, SQLITE3, HOTLOADDBType}
	for _, db := range supportedDBs {
		if c.dbConfig.Type == db {
			return true
		}
	}
	return false
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
