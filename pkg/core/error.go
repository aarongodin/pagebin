package core

import (
	"github.com/joomcode/errorx"
)

var (
	traitUnexpected = errorx.RegisterTrait("unexpected")

	errApp                   = errorx.NewNamespace("app")
	ErrUnknown               = errorx.NewType(errApp, "unknown", traitUnexpected)
	ErrUnknownBlobBackend    = errorx.NewType(errApp, "unknown_blob_backend", traitUnexpected)
	ErrPageNotFound          = errorx.NewType(errApp, "page_not_found", errorx.NotFound())
	ErrThemeNotCompiled      = errorx.NewType(errApp, "theme_not_compiled", traitUnexpected)
	ErrThemeTemplateNotFound = errorx.NewType(errApp, "theme_template_not_found")
	ErrThemeTemplateExec     = errorx.NewType(errApp, "theme_template_exec")
	ErrVersionNotCompiled    = errorx.NewType(errApp, "version_not_compiled", traitUnexpected)
	ErrReservedPath          = errorx.NewType(errApp, "reserved_path")
	ErrInvalidVersion        = errorx.NewType(errApp, "invalid_version")
	ErrUIDRequired           = errorx.NewType(errApp, "uid_required")

	errStore                = errorx.NewNamespace("store")
	ErrItemNotFound         = errorx.NewType(errStore, "item_not_found", errorx.NotFound())
	ErrBucketNotFound       = errorx.NewType(errStore, "bucket_not_found", errorx.NotFound())
	ErrTransactionNotFound  = errorx.NewType(errStore, "tx_not_found", traitUnexpected)
	ErrTransactionEnd       = errorx.NewType(errStore, "tx_end", traitUnexpected)
	ErrTransactionPrivilege = errorx.NewType(errStore, "tx_privilege", traitUnexpected)
)
