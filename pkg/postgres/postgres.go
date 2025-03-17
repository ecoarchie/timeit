package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Postgres -.
type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Pool *pgxpool.Pool
}

// New -.
func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize) //nolint:gosec // skip integer overflow conversion int -> int32
	// poolConfig.AfterConnect(context.Background(), )
	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		// Manually register the composite type rr_tod
		rrTOD, err := conn.LoadType(ctx, "rr_tod")
		if err != nil {
			fmt.Println("error loading type rr_tod")
			return err
		}
		conn.TypeMap().RegisterType(rrTOD)

		rrTODarray, err := conn.LoadType(ctx, "_rr_tod")
		if err != nil {
			fmt.Println("error loading type rr_tod")
			return err
		}
		conn.TypeMap().RegisterType(rrTODarray)
		log.Println("Registered custom type rr_tod")
		return nil
	}
	// Fields: []pgtype.CompositeFields{
	// 	{Name: "reader_id", Type: &pgtype.UUID{}},
	// 	{Name: "tod", Type: &pgtype.Timestamp{}},
	// },

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			if pingErr := pg.Pool.Ping(context.Background()); pingErr == nil {
				break
			} else {
				// TODO Should I close the old connection here?
				log.Println("failed ping to database")
			}
		}
		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}
	if err := pg.Pool.Ping(context.Background()); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
