package store

import (
	"context"
	"crypto/sha256"
	"os"
	"path/filepath"

	"github.com/aarongodin/pagebin/pkg/config"
	"github.com/aarongodin/pagebin/pkg/core"
	"github.com/oklog/ulid/v2"
	bolt "go.etcd.io/bbolt"
)

const (
	BlobBackendLocalFS = "localfs"
	BlobBackendS3      = "s3"
)

type BlobStore interface {
	GetBlob(ctx context.Context, uid ulid.ULID) (core.Blob, error)
	GetBytes(ctx context.Context, uid ulid.ULID) ([]byte, error)
	CreateBlob(ctx context.Context, raw []byte) (core.Blob, error)
	UpdateBlob(ctx context.Context, uid ulid.ULID, raw []byte) (core.Blob, error)
}

type localFSBlobStore struct {
	rootDir string
	db      documentDB[core.Blob]
}

func (s localFSBlobStore) GetBlob(ctx context.Context, uid ulid.ULID) (core.Blob, error) {
	return s.db.One(ctx, bucketBlobs, uid.String())
}

func (s localFSBlobStore) GetBytes(ctx context.Context, uid ulid.ULID) ([]byte, error) {
	return os.ReadFile(filepath.Join(s.rootDir, uid.String()))
}

func (s localFSBlobStore) CreateBlob(ctx context.Context, raw []byte) (core.Blob, error) {
	uid := ulid.Make()
	if err := os.WriteFile(filepath.Join(s.rootDir, uid.String()), raw, 0755); err != nil {
		return core.Blob{}, err
	}
	hasher := sha256.New()
	if _, err := hasher.Write(raw); err != nil {
		return core.Blob{}, err
	}
	blob := core.Blob{
		UID:  uid,
		Hash: hasher.Sum(nil),
	}
	if err := s.db.Save(ctx, bucketBlobs, uid.String(), blob); err != nil {
		return core.Blob{}, err
	}
	return blob, nil
}

func (s localFSBlobStore) UpdateBlob(ctx context.Context, uid ulid.ULID, raw []byte) (core.Blob, error) {
	blob, err := s.db.One(ctx, bucketBlobs, uid.String())
	if err != nil {
		return core.Blob{}, err
	}
	hasher := sha256.New()
	if _, err := hasher.Write(raw); err != nil {
		return core.Blob{}, err
	}
	blob.Hash = hasher.Sum(nil)
	if err := os.WriteFile(filepath.Join(s.rootDir, uid.String()), raw, 0755); err != nil {
		return core.Blob{}, err
	}
	if err := s.db.Save(ctx, bucketBlobs, uid.String(), blob); err != nil {
		return core.Blob{}, err
	}
	return blob, nil
}

func NewBlobStore(rc *config.RuntimeConfig, db *bolt.DB) (BlobStore, error) {
	switch rc.BlobBackend {
	case BlobBackendLocalFS:
		if _, err := os.Stat(rc.BlobLocalFSRootDir); os.IsNotExist(err) {
			err := os.Mkdir(rc.BlobLocalFSRootDir, 0755)
			if err != nil {
				return nil, err
			}
		}
		return &localFSBlobStore{
			rootDir: rc.BlobLocalFSRootDir,
			db:      docDB[core.Blob]{db},
		}, nil
	default:
		return nil, core.ErrUnknownBlobBackend.NewWithNoMessage()
	}
}
