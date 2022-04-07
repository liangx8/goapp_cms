package mgr

import "io"

type (
	DBI interface {
		Load(d any) error
		Close()
	}
	View func(w io.Writer, data any) error
)
