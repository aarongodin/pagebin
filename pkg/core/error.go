package core

import (
	"github.com/joomcode/errorx"
)

var (
	// traitUnexpected = errorx.RegisterTrait("unexpected")

	errApp                   = errorx.NewNamespace("app")
	ErrUnknownBlobBackend    = errorx.NewType(errApp, "unknown_blob_backend")
	ErrPageNotFound          = errorx.NewType(errApp, "page_not_found", errorx.NotFound())
	ErrThemeNotCompiled      = errorx.NewType(errApp, "theme_not_compiled")
	ErrThemeTemplateNotFound = errorx.NewType(errApp, "theme_template_not_found")
	ErrThemeTemplateExec     = errorx.NewType(errApp, "theme_template_exec")
	ErrVersionNotCompiled    = errorx.NewType(errApp, "version_not_compiled")

	errStore          = errorx.NewNamespace("store")
	ErrItemNotFound   = errorx.NewType(errStore, "item_not_found", errorx.NotFound())
	ErrBucketNotFound = errorx.NewType(errStore, "bucket_not_found", errorx.NotFound())
)
