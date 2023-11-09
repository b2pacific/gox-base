package goxServer

import (
	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-base/config"
	"net/http"
	"sync"
)

type Server interface {
	Start(handler http.Handler, appConfig *config.App) error
	Stop() chan bool
}

type ServerShutdownHook interface {
	StopFunction() func()
}

func NewServer(cf gox.CrossFunction) (Server, error) {
	s := &serverImpl{CrossFunction: cf, stopOnce: &sync.Once{}}
	return s, nil
}

func NewServerWithShutdownHookFunc(cf gox.CrossFunction, shutdownHookFunc func()) (Server, error) {
	s := &serverImpl{CrossFunction: cf, shutdownHookFunc: shutdownHookFunc, stopOnce: &sync.Once{}}
	return s, nil
}

func NewServerWithShutdownHook(cf gox.CrossFunction, serverShutdownHook ServerShutdownHook) (Server, error) {
	s := &serverImpl{CrossFunction: cf, shutdownHookFunc: serverShutdownHook.StopFunction(), stopOnce: &sync.Once{}}
	return s, nil
}
