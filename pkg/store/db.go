package store

import (
	"context"

	"github.com/aarongodin/pagebin/pkg/core"
	bolt "go.etcd.io/bbolt"
)

var (
	bucketApp               = "app"
	bucketThemes            = "themes"
	bucketPages             = "pages"
	bucketVersions          = "versions"
	bucketBlobs             = "blobs"
	bucketIndex             = "index"
	buckets                 = [...]string{bucketApp, bucketThemes, bucketPages, bucketVersions, bucketBlobs, bucketIndex}
	bucketIndexPageVersions = "page-versions"
	nestedBuckets           = map[string]string{
		bucketIndex: bucketIndexPageVersions,
	}

	contextKeyTransaction         = core.ContextKey("transaction")
	contextKeyTransactionWritable = core.ContextKey("transaction-writable")
)

func (s *store) StartTx(ctx context.Context, writable bool) (context.Context, error) {
	tx, err := s.db.Begin(writable)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, contextKeyTransaction, tx)
	ctx = context.WithValue(ctx, contextKeyTransactionWritable, writable)
	return ctx, nil
}

func (s *store) EndTx(ctx context.Context, txErr error) error {
	tx, txOK := ctx.Value(contextKeyTransaction).(*bolt.Tx)
	if !txOK || tx == nil {
		return core.ErrTransactionNotFound.NewWithNoMessage()
	}
	writable, writableOK := ctx.Value(contextKeyTransactionWritable).(bool)
	if !writableOK {
		return core.ErrTransactionNotFound.NewWithNoMessage()
	}

	if txErr != nil {
		if err := tx.Rollback(); err != nil {
			return core.ErrTransactionEnd.Wrap(err, "failed ending transaction; original err: %s", txErr.Error())
		}
		return txErr
	}

	var err error
	if writable {
		err = tx.Commit()
	} else {
		err = tx.Rollback()
	}
	if err != nil {
		return core.ErrTransactionEnd.WrapWithNoMessage(err)
	}
	return nil
}

func transactCtx(ctx context.Context, db *bolt.DB, writable bool, fn func(tx *bolt.Tx) error) error {
	ctxWritable, ctxWritableOk := ctx.Value(contextKeyTransactionWritable).(bool)
	if writable && ctxWritableOk && !ctxWritable {
		return core.ErrTransactionPrivilege.New("expected transaction to be writable")
	}
	tx, externalTX := ctx.Value(contextKeyTransaction).(*bolt.Tx)
	var err error
	if !externalTX {
		tx, err = db.Begin(writable)
		if err != nil {
			return err
		}
	}

	done := make(chan error, 1)
	go func() {
		done <- fn(tx)
	}()

	select {
	case <-ctx.Done():
		if err := tx.Rollback(); err != nil {
			return err
		}
		return ctx.Err()
	case userErr := <-done:
		if externalTX {
			if userErr == nil {
				return nil
			} else {
				if err := tx.Rollback(); err != nil {
					return err
				}
				return userErr
			}
		} else {
			if userErr == nil {
				if writable {
					return tx.Commit()
				} else {
					return tx.Rollback()
				}
			} else {
				if err := tx.Rollback(); err != nil {
					return err
				}
				return userErr
			}
		}
	}
}
