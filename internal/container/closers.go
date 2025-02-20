package container

import "github.com/sknv/protomock/pkg/closer"

func (a *Application) AddCloser(closer closer.Closer) {
	a.closers.Add(closer)
}
