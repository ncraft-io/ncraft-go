package memory

import (
    "github.com/ncraft-io/ncraft-go/pkg/config/loader"
    "github.com/ncraft-io/ncraft-go/pkg/config/reader"
    "github.com/ncraft-io/ncraft-go/pkg/config/source"
)

// WithSource appends a source to list of sources
func WithSource(s source.Source) loader.Option {
    return func(o *loader.Options) {
        o.Source = append(o.Source, s)
    }
}

// WithReader sets the config reader
func WithReader(r reader.Reader) loader.Option {
    return func(o *loader.Options) {
        o.Reader = r
    }
}
