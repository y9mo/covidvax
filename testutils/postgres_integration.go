package testutils

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/stretchr/testify/suite"
)

type IntegrationSuite struct {
	suite.Suite
	postgresContainer *postgresContainer
	db                *gorm.DB
	fixtures          *testfixtures.Loader
}

type postgresContainer struct {
	testcontainers.Container
	PGConnect string
}

func (s *IntegrationSuite) SetupPostgres(ctx context.Context) {
	srcMount, err := filepath.Abs("../docker.d/db/init.sql")
	s.Require().NoError(err)

	req := testcontainers.ContainerRequest{
		Image:        "postgres:13.4-alpine",
		ExposedPorts: []string{"5432/tcp"},
		BindMounts:   map[string]string{srcMount: "/docker-entrypoint-initdb.d/init-db.sql"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(time.Second * 5),
		Env:          map[string]string{"POSTGRES_PASSWORD": "password"},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)

	mappedPort, err := container.MappedPort(ctx, "5432")
	s.Require().NoError(err)

	hostIP, err := container.Host(ctx)
	s.Require().NoError(err)

	pgConnect := fmt.Sprintf("host=%s port=%s sslmode=disable user=admin dbname=covidvax-test password=admin-pwd",
		hostIP, mappedPort.Port())
	s.postgresContainer = &postgresContainer{Container: container, PGConnect: pgConnect}
}

func (s *IntegrationSuite) InitSQLClient() {
	var err error
	s.db, err = gorm.Open("postgres", s.postgresContainer.PGConnect)
	s.Require().NoError(err)
}

func (s *IntegrationSuite) PurgeContainer() {
	if s.postgresContainer != nil {
		s.postgresContainer.Terminate(context.Background())
	}
}

func (s *IntegrationSuite) DB() *gorm.DB {
	return s.db
}

func (s *IntegrationSuite) ApplyMigrations() {
	driver, err := postgres.WithInstance(s.db.DB(), &postgres.Config{})
	if err != nil {
		log.Fatalf("error while creating postgres-migrate: %s", err)
	}
	migrate, err := migrate.NewWithDatabaseInstance(
		"file://../db/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("error while creating migrate instance: %s", err)
	}
	err = migrate.Up()
	if err != nil {
		log.Fatalf("error while applying migration: %s", err)
	}
}

func (s *IntegrationSuite) SetupFixtures() {
	var err error
	s.fixtures, err = testfixtures.New(
		testfixtures.Database(s.db.DB()),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("../testdata/fixtures"),
		testfixtures.UseAlterConstraint(),
	)
	s.Require().NoError(err)
}

func (s *IntegrationSuite) LoadDatabaseWithFixtures() {
	s.Require().NoError(s.fixtures.Load())
}

func (s *IntegrationSuite) Cleanup() {
	truncateQuery := `TRUNCATE TABLE appointment_bookings, appointments, treatment_centers, patients;`

	err := s.db.Exec(truncateQuery).Error
	if err != nil {
		log.Fatal("impossible to cleanup DB", err)
	}
}

func (s *IntegrationSuite) SetupTest() {
	s.LoadDatabaseWithFixtures()
}

func (s *IntegrationSuite) TearDownTest() {
	s.Cleanup()
}

func (s *IntegrationSuite) SetupSuite() {
	s.SetupPostgres(context.Background())
	s.InitSQLClient()
	s.ApplyMigrations()
	s.SetupFixtures()
}

func (s *IntegrationSuite) TearDownSuite() {
	s.PurgeContainer()
}
