package store

import (
	"context"
	"os"
	"path/filepath"

	"github.com/aarongodin/pagebin/pkg/config"
	"github.com/aarongodin/pagebin/pkg/core"
	"github.com/oklog/ulid/v2"
)

const (
	BlobBackendLocalFS = "localfs"
	BlobBackendS3      = "s3"
)

type BlobStore interface {
	PutBytes(ctx context.Context, raw []byte) (ulid.ULID, error)
	GetBytes(ctx context.Context, uid ulid.ULID) ([]byte, error)
}

type localFSBlobStore struct {
	rootDir string
}

func (s localFSBlobStore) PutBytes(ctx context.Context, raw []byte) (ulid.ULID, error) {
	uid := ulid.Make()
	if err := os.WriteFile(filepath.Join(s.rootDir, uid.String()), raw, 0755); err != nil {
		return ulid.ULID{}, err
	}
	return uid, nil
}

func (s localFSBlobStore) GetBytes(ctx context.Context, uid ulid.ULID) ([]byte, error) {
	return os.ReadFile(filepath.Join(s.rootDir, uid.String()))
}

func NewBlobStore(rc *config.RuntimeConfig) (BlobStore, error) {
	switch rc.BlobBackend {
	case BlobBackendLocalFS:
		if _, err := os.Stat(rc.BlobLocalFSRootDir); os.IsNotExist(err) {
			err := os.Mkdir(rc.BlobLocalFSRootDir, 0755)
			if err != nil {
				return nil, err
			}
		}
		return &localFSBlobStore{rootDir: rc.BlobLocalFSRootDir}, nil
	default:
		return nil, core.ErrUnknownBlobBackend.NewWithNoMessage()
	}
}
