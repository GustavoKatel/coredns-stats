package stats_test

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-faker/faker/v4"
	_ "github.com/lib/pq"

	stats "github.com/GustavoKatel/coredns-stats"
)

func (s *StatsBackendSuite) TestBackendPostgres() {
	s.Run("Start runs the init query", func() {
		prefix := faker.Word()
		_, err := stats.PrepareStatsBackend(s.connString, 1, 1*time.Second, prefix, time.Hour, "@daily", s.logger)
		s.Require().NoError(err)

		db, err := sql.Open("postgres", s.connString)
		s.Require().NoError(err)

		res, err := db.Exec(fmt.Sprintf("INSERT INTO %s_stats (domain, question_type, metadata) VALUES ('example.com', 'A', '{}'::jsonb)", prefix))
		s.Require().NoError(err)

		rows, err := res.RowsAffected()
		s.Require().NoError(err)
		s.Equal(int64(1), rows)
	})

	s.Run("store stats", func() {
		prefix := faker.Word()
		backend, err := stats.PrepareStatsBackend(s.connString, 1, 1*time.Second, prefix, time.Hour, "@daily", s.logger)
		s.Require().NoError(err)
		s.Require().NotNil(backend)

		backend.Store("example.com", "A", stats.Metadata{"key": "value"})

		<-time.After(1 * time.Second)

		db, err := sql.Open("postgres", s.connString)
		s.Require().NoError(err)

		var domain, questionType string
		var metadata stats.Metadata
		err = db.QueryRow(fmt.Sprintf("SELECT domain, question_type, metadata FROM %s_stats", prefix)).Scan(&domain, &questionType, &metadata)
		s.Require().NoError(err)

		s.Equal("example.com", domain)
		s.Equal("A", questionType)
		s.Equal(stats.Metadata{"key": "value"}, metadata)
	})

	s.Run("cleanup old entries on cron", func() {
		prefix := faker.Word()
		backend, err := stats.PrepareStatsBackend(s.connString, 1, 1*time.Second, prefix, time.Second, "@every 2s", s.logger)
		s.Require().NoError(err)
		s.Require().NotNil(backend)

		backend.Store("example.com", "A", stats.Metadata{"key": "value"})

		<-time.After(3 * time.Second)

		db, err := sql.Open("postgres", s.connString)
		s.Require().NoError(err)

		var count int64
		err = db.QueryRow(fmt.Sprintf("SELECT count(*) FROM %s_stats", prefix)).Scan(&count)
		s.Require().NoError(err)

		s.Equal(int64(0), count)
	})

}
