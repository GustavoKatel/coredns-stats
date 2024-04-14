package stats

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type storeRequest struct {
	domain       string
	questionType string
	metadata     Metadata
}

type backendPostgres struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	uri          string
	workers      int64
	queryTimeout time.Duration
	statsPrefix  string

	logger Logger

	db        *sql.DB
	tableName string

	dataChan chan storeRequest
}

//go:embed sql/init_postgres.sql
var queryInit string

//go:embed sql/insert_stats_postgres.sql
var queryInsert string

func newBackendPostgres(uri string, workers int64, queryTimeout time.Duration, statsPrefix string, logger Logger) *backendPostgres {
	ctx, cancel := context.WithCancel(context.Background())

	return &backendPostgres{
		ctx:       ctx,
		ctxCancel: cancel,

		uri:          uri,
		workers:      workers,
		queryTimeout: queryTimeout,
		statsPrefix:  statsPrefix,

		tableName: statsPrefix + "_stats",

		logger: logger,
	}
}

func (b *backendPostgres) Start() error {
	db, err := sql.Open("postgres", b.uri)
	if err != nil {
		return err
	}

	b.db = db

	query := fmt.Sprintf(queryInit, b.tableName)
	if _, err := b.db.Exec(query); err != nil {
		return fmt.Errorf("error trying to init database: %w", err)
	}

	b.dataChan = make(chan storeRequest, b.workers)

	for i := int64(0); i < b.workers; i++ {
		go b.worker()
	}

	return nil
}

func (b *backendPostgres) Store(domain string, questionType string, metadata Metadata) {
	b.dataChan <- storeRequest{
		domain:       domain,
		questionType: questionType,
		metadata:     metadata,
	}
}

func (b *backendPostgres) Stop() error {
	b.ctxCancel()
	return b.db.Close()
}

func (b *backendPostgres) insert(data storeRequest) {
	ctx, cancel := context.WithTimeout(b.ctx, b.queryTimeout)
	defer cancel()

	query := fmt.Sprintf(queryInsert, b.tableName)
	_, err := b.db.ExecContext(ctx, query, data.domain, data.questionType, &data.metadata)
	if err != nil {
		b.logger.Error("failed to insert data into database", err)
	}
}

func (b *backendPostgres) worker() {
	for {
		select {
		case <-b.ctx.Done():
			return
		case data := <-b.dataChan:
			b.insert(data)
		}
	}
}
