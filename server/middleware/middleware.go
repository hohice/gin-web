package middleware

import (
	"github.com/hohice/gin-web/pkg/setting"
)

func Init() error {
	conf := &setting.Config
	for _, fn := range Initlist {
		if err, closeFn := fn(conf); err != nil {
			Close()
			return err
		} else {
			Closelist = append(Closelist, closeFn)
		}
	}
	return nil
}

func Close() {
	for _, fn := range Closelist {
		fn()
	}
	Closelist = []Closeble{}
}

type Closeble func()
type Register func(conf *setting.Configs) (error, Closeble)

var Initlist []Register
var Closelist []Closeble

func registerSelf(regfunc Register) {
	Initlist = append(Initlist, regfunc)
}
