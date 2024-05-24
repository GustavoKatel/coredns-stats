package stats_test

import (
	"context"
	"testing"
	"time"

	stats "github.com/GustavoKatel/coredns-stats"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type StatsBackendSuite struct {
	suite.Suite

	pgContainer *postgres.PostgresContainer
	connString  string
	logger      stats.Logger
}

func TestStatsBackendSuite(t *testing.T) {
	suite.Run(t, new(StatsBackendSuite))
}

func (s *StatsBackendSuite) SetupSuite() {
	ctx := context.Background()
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second),
			wait.ForExposedPort(),
		),
	)
	if err != nil {
		s.FailNow(err.Error())
	}

	s.pgContainer = pgContainer

	connStr, err := s.pgContainer.ConnectionString(ctx, "sslmode=disable")
	s.Require().NoError(err)

	s.connString = connStr

	s.logger = &logger{}
}

func (s *StatsBackendSuite) TearDownSuite() {
	if err := s.pgContainer.Terminate(context.Background()); err != nil {
		s.FailNowf("failed to terminate pgContainer: %s", err.Error())
	}
}

func (s *StatsBackendSuite) TestPrepareStatsBackend() {
	backend, err := stats.PrepareStatsBackend(s.connString, 1, 1*time.Second, "coredns", time.Hour, "@daily", nil)
	s.Require().NoError(err)
	s.Require().NotNil(backend)
}
