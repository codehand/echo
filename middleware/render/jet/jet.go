/*

   Copyright 2017 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package jet

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sync"

	. "github.com/admpub/jet"
	"github.com/admpub/log"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/logger"
	"github.com/webx-top/echo/middleware/render"
	"github.com/webx-top/echo/middleware/render/driver"
)

func init() {
	render.Reg(`jet`, func(tmplDir string) driver.Driver {
		return New(tmplDir)
	})
}

func New(templateDir string, args ...logger.Logger) driver.Driver {
	var err error
	templateDir, err = filepath.Abs(templateDir)
	if err != nil {
		panic(err.Error())
	}
	a := &Jet{
		NopRenderer: &driver.NopRenderer{},
		templateDir: templateDir,
		set:         NewHTMLSet(templateDir),
	}
	if len(args) > 0 {
		a.logger = args[0]
	} else {
		a.logger = log.New("render-jet")
	}
	return a
}

type Jet struct {
	*driver.NopRenderer
	mutex         sync.RWMutex
	set           *Set
	templateDir   string
	logger        logger.Logger
	debug         bool
	tmplPathFixer func(echo.Context, string) string
}

func (self *Jet) Debug() bool {
	return self.debug
}

func (self *Jet) SetDebug(on bool) {
	self.debug = on
	self.set.SetDevelopmentMode(on)
}

func (self *Jet) SetLogger(l logger.Logger) {
	self.logger = l
}

func (self *Jet) Logger() logger.Logger {
	return self.logger
}

func (self *Jet) TmplDir() string {
	return self.templateDir
}

func (self *Jet) SetTmplPathFixer(fn func(echo.Context, string) string) {
	self.tmplPathFixer = fn
}

func (self *Jet) SetFuncMap(fn func() map[string]interface{}) {
	for name, fn := range fn() {
		self.set.AddGlobal(name, fn)
	}
}

func (self *Jet) Render(w io.Writer, tmpl string, data interface{}, c echo.Context) error {
	if self.tmplPathFixer != nil {
		tmpl = self.tmplPathFixer(c, tmpl)
	}
	t, err := self.set.GetTemplate(tmpl)
	if err != nil {
		return err
	}
	vars := make(VarMap)
	for name, fn := range c.Funcs() {
		vars.Set(name, fn)
	}
	return t.Execute(w, vars, data)
}

func (self *Jet) Fetch(tmpl string, data interface{}, c echo.Context) string {
	if self.tmplPathFixer != nil {
		tmpl = self.tmplPathFixer(c, tmpl)
	}
	w := new(bytes.Buffer)
	t, err := self.set.GetTemplate(tmpl)
	if err != nil {
		return fmt.Sprintf("Parse %v err: %v", tmpl, err)
	}
	vars := make(VarMap)
	for name, fn := range c.Funcs() {
		vars.Set(name, fn)
	}
	err = t.Execute(w, vars, data)
	if err != nil {
		return fmt.Sprintf("Parse %v err: %v", tmpl, err)
	}
	return w.String()
}

func (self *Jet) RawContent(tmpl string) (b []byte, e error) {
	return nil, errors.New(`unsupported`)
}
