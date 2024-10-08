package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/milosgajdos/go-hypher"
)

// Syncer syncs graph to sqlite.
type Syncer struct {
	db *DB
}

// NewSyncer creates a new sqlite syncer and returns it.
func NewSyncer(db *DB) (*Syncer, error) {
	return &Syncer{
		db: db,
	}, nil
}

// Sync sync graph g to sqlite DB.
func (s *Syncer) Sync(ctx context.Context, g hypher.Graph) error {
	tx, err := s.db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// nolint:errcheck
	defer tx.Rollback()

	if err := s.syncGraph(ctx, tx, g); err != nil {
		return err
	}

	nodes := g.Nodes()
	for nodes.Next() {
		n, ok := nodes.Node().(hypher.Node)
		if !ok {
			continue
		}
		if err := s.syncNode(ctx, tx, g.UID(), n); err != nil {
			return err
		}
	}

	edges := g.Edges()
	for edges.Next() {
		e, ok := edges.Edge().(hypher.Edge)
		if !ok {
			continue
		}
		if err := s.syncEdge(ctx, tx, g.UID(), e); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// syncGraph initializes the graph entry in the database.
func (s *Syncer) syncGraph(ctx context.Context, tx *sql.Tx, g hypher.Graph) error {
	attrs, err := json.Marshal(g.Attrs())
	if err != nil {
		return err
	}

	createdAt := time.Now()
	updatedAt := createdAt

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO graphs (
			uid,
			label,
			attrs,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?)
	`,
		g.UID(),
		g.Label(),
		attrs,
		(*NullTime)(&createdAt),
		(*NullTime)(&updatedAt),
	); err != nil {
		return err
	}

	return nil
}

// syncNode stores node in the sqlite DB.
func (s *Syncer) syncNode(ctx context.Context, tx *sql.Tx, graphUID string, n hypher.Node) error {
	createdAt := time.Now()
	updatedAt := createdAt

	attrs, err := json.Marshal(n.Attrs())
	if err != nil {
		return err
	}

	// Execute insertion query.
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO nodes (
			uid,
			graph,
			label,
			attrs,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		n.UID(),
		graphUID,
		n.Label(),
		string(attrs),
		(*NullTime)(&createdAt),
		(*NullTime)(&updatedAt),
	); err != nil {
		return err
	}

	return nil
}

// syncEdge stores edge in the sqlite DB.
func (s *Syncer) syncEdge(ctx context.Context, tx *sql.Tx, graphUID string, e hypher.Edge) error {
	createdAt := time.Now()
	updatedAt := createdAt

	attrs, err := json.Marshal(e.Attrs())
	if err != nil {
		return err
	}

	// Retrieve source and target node UIDs directly
	sourceUID := e.From().(hypher.Node).UID()
	targetUID := e.To().(hypher.Node).UID()

	// Execute insertion query.
	_, err = tx.ExecContext(ctx, `
		INSERT INTO edges (
			uid,
			graph,
			source,
			target,
			label,
			weight,
			attrs,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
	`,
		e.UID(),
		graphUID,
		sourceUID,
		targetUID,
		e.Label(),
		e.Weight(),
		string(attrs),
		(*NullTime)(&createdAt),
		(*NullTime)(&updatedAt),
	)
	if err != nil {
		return err
	}

	return nil
}
