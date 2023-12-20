package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type DBOptions struct {
	Timeout time.Duration
}

// DBCredentials structure.
type DBCredentials struct {
	DBUser           string `json:"DBUser"`
	DBPassword       string `json:"DBPassword"`
	DBHost           string `json:"DBHost"`
	DBPort           int    `json:"DBPort"`
	DB               string `json:"DB"`
	CloudSQLInstance string `json:"CloudSqlInstance"`
}

type database struct {
	db  *bun.DB
	ctx context.Context
	// DbTx is a transaction object that can be used to run queries in a transaction.
	DbTx DatabaseTransactionInterface
}

type databaseTransaction struct {
	db  *bun.DB
	ctx context.Context
}
type DatabaseTransactionInterface interface {
	// RunInTransaction runs the given function in a database transaction.
	// auto rollback the transaction if the function returns an error.
	//
	// Parameters:
	// - function: func() error
	//
	// Returns:
	// - error: error
	RunInTransaction(function func() error) error
}

// Returns the database operational results
type Result struct {
	// RowsAffected: number of rows affected by the query
	RowsAffected int64
}

type DatabaseInterface interface {
	// Connect to database using the dsn ("host= database= user= port= password= ").
	// It returns the context and database connection.
	// It uses the default timeout of 30 seconds.
	// If you need to change the timeout, use ConnectWithTimeout.
	//
	// Parameters:
	// - dbCredentials: database.DBCredentials
	//
	// Returns:
	// - error: error
	Connect(dbCredentials DBCredentials) error

	// ConnectWithTimeout to database using the dsn ("host= database= user= port= password= ") and DBOptions
	// It returns the context and database connection.
	//
	// Parameters:
	// - dbCredentials: database.DBCredentials
	// - options: *database.DBOptions
	//
	// Returns:
	// - error: error
	ConnectWithTimeout(dbCredentials DBCredentials, options *DBOptions) error

	// Exec executes the given sql statement.
	//
	// Parameters:
	// - sql: string
	//
	// Returns:
	// - error: error
	Exec(sql string) error
}

// New instance of CommunityDatabase
//
// Returns:
// - CommunityDatabaseInterface: CommunityDatabaseInterface
func New() DatabaseInterface {
	return &database{}
}

// Connect to database using the dsn ("host= database= user= port= password= ").
// It returns the context and database connection.
// It uses the default timeout of 30 seconds.
// If you need to change the timeout, use ConnectWithTimeout.
//
// Parameters:
// - dbCredentials: database.DBCredentials
//
// Returns:
// - error: error
func (udb *database) Connect(dbCredentials DBCredentials) error {
	return udb.ConnectWithTimeout(dbCredentials, &DBOptions{Timeout: time.Second * 30})
}

// ConnectWithTimeout to database using the dsn ("host= database= user= port= password= ") and DBOptions
// It returns the context and database connection.
//
// Parameters:
// - dbCredentials: database.DBCredentials
// - options: *database.DBOptions
//
// Returns:
// - error: error
func (udb *database) ConnectWithTimeout(
	dbCredentials DBCredentials,
	options *DBOptions,
) error {
	db, ctx, err := connect(dbCredentials, options)
	if err != nil {
		err = fmt.Errorf("unable to connect to database: %v", err)
		log.Println(err)
		return err
	}
	udb.db = db
	udb.ctx = ctx
	udb.DbTx = &databaseTransaction{db: db, ctx: ctx}

	// Todo: add tables
	return nil
}

// RunInTransaction runs the given function in a database transaction.
// auto rollback the transaction if the function returns an error.
//
// Parameters:
// - function: func() error
//
// Returns:
// - error: error
func (udbTx *databaseTransaction) RunInTransaction(function func() error) error {
	return udbTx.db.RunInTx(udbTx.ctx, nil, func(context context.Context, tx bun.Tx) error {
		return function()
	})
}

// Exec executes the given sql statement.
//
// Parameters:
// - sql: string
//
// Returns:
// - error: error
func (udb *database) Exec(sql string) error {
	_, err := udb.db.ExecContext(udb.ctx, sql)
	if err != nil {
		err = fmt.Errorf("unable to execute the sql on database: %v", err)
		log.Println(err)
		return err
	}
	return nil
}

// connect to database using the dsn ("host= database= user= port= password= ").
// It returns the context and database connection.
func connect(dbCredentials DBCredentials, opts *DBOptions) (*bun.DB, context.Context, error) {
	ctx := context.Background()
	network := "tcp"
	if strings.HasPrefix(dbCredentials.DBHost, "/") {
		network = "unix"
	} else {
		dbCredentials.DBHost += ":" + strconv.Itoa(dbCredentials.DBPort)
	}
	conn := pgdriver.NewConnector(
		pgdriver.WithNetwork(network),
		pgdriver.WithAddr(dbCredentials.DBHost),
		pgdriver.WithUser(dbCredentials.DBUser),
		pgdriver.WithPassword(dbCredentials.DBPassword),
		pgdriver.WithDatabase(dbCredentials.DB),
		pgdriver.WithInsecure(true),
		pgdriver.WithTimeout(opts.Timeout),
	)

	sqldb := sql.OpenDB(conn)
	db := bun.NewDB(sqldb, pgdialect.New())

	if getOsEnv("DB_QUERY_DEBUG", "") != "" {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
	}
	_, err := db.ExecContext(ctx, "SELECT 1")
	if err != nil {
		err = fmt.Errorf("unable to connect to database: %v", err)
		log.Println(err)
		return db, ctx, err
	}
	return db, ctx, nil
}

func getOsEnv(key string, defaultValue string) string {
	value, lookup := os.LookupEnv(key)
	if !lookup {
		return defaultValue
	}
	return value
}
