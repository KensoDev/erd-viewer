package db

// This package is deprecated. Use github.com/kensodev/erd-viewer/pkg/erd instead.
// These type aliases are provided for backward compatibility.

import (
	"github.com/kensodev/erd-viewer/pkg/erd"
	"github.com/kensodev/erd-viewer/pkg/erd/postgres"
)

// Deprecated: Use erd.Column instead
type Column = erd.Column

// Deprecated: Use erd.Table instead
type Table = erd.Table

// Deprecated: Use erd.ForeignKey instead
type ForeignKey = erd.ForeignKey

// Deprecated: Use erd.SchemaData instead
type SchemaData = erd.SchemaData

// Deprecated: Use postgres.Introspector instead
type Introspector = postgres.Introspector

// Deprecated: Use postgres.NewIntrospector instead
var NewIntrospector = postgres.NewIntrospector
