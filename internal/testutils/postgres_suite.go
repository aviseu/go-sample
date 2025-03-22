package testutils

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	mpg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

type PostgresSuite struct {
	suite.Suite

	container testcontainers.Container
	m         *migrate.Migrate
	DB        *sqlx.DB
	BadDB     *sqlx.DB
}

func (suite *PostgresSuite) SetupSuite() {
	ctx := context.Background()

	suite.container, suite.m, suite.DB, suite.BadDB = suite.createDependencies(ctx)
}

func (suite *PostgresSuite) createDependencies(ctx context.Context) (testcontainers.Container, *migrate.Migrate, *sqlx.DB, *sqlx.DB) {
	c, err := postgres.Run(
		ctx,
		"postgres:17",
		postgres.WithUsername("api"),
		postgres.WithPassword("pwd"),
		postgres.WithDatabase("todo"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	suite.NoError(err)

	dsn, err := c.ConnectionString(ctx, "sslmode=disable")
	suite.NoError(err)

	db, err := sqlx.Connect("postgres", dsn)
	suite.NoError(err)

	badDB, err := sqlx.Connect("postgres", dsn)
	suite.NoError(err)
	suite.NoError(badDB.Close())

	driver, err := mpg.WithInstance(db.DB, &mpg.Config{})
	suite.NoError(err)

	m, err := migrate.NewWithDatabaseInstance("file://../../../../configs/migrations", "postgres", driver)
	suite.NoError(err)

	return c, m, db, badDB
}

func (suite *PostgresSuite) TearDownSuite() {
	suite.NoError(suite.DB.Close())
	suite.NoError(suite.container.Terminate(context.Background()))
}

func (suite *PostgresSuite) SetupTest() {
	suite.NoError(suite.m.Up())
}

func (suite *PostgresSuite) TearDownTest() {
	suite.NoError(suite.m.Down())
}
