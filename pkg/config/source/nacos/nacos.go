package nacos

import "github.com/ncraft-io/ncraft-go/pkg/config/source"

type nacosSource struct {
}

func (n *nacosSource) Read() (*source.ChangeSet, error) {

    return nil, nil
}

func (n *nacosSource) Write(*source.ChangeSet) error {
    return nil
}

func (n *nacosSource) Watch() (source.Watcher, error) {
    return nil, nil
}

func (n *nacosSource) String() string {
    return ""
}

func NewSource(opts ...source.Option) source.Source {
    return &nacosSource{}
}
