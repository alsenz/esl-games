package lesson2

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"html/template"
	"strings"
)

type Act struct {
	Scenes []Scene
	ActRepeaterExpr ActRepeaterExpr
}

type Acts []Act

type ActRepeaterExpr string

func (expr ActRepeaterExpr) EvalBool(checkpoint Checkpoint) (bool, error) {
	// Just translate the expression into a parseable if expression
	tpl, err := template.New("base").Funcs(sprig.HermeticHtmlFuncMap()).Parse(
		string("{{if " + expr + "}}TRUE{{else}}FALSE{{end}}"))
	if err != nil {
		return false, err
	}
	buf := &bytes.Buffer{}
	if err = tpl.Execute(buf, checkpoint); err != nil {
		return false, err
	}
	return strings.TrimSpace(buf.String()) == "TRUE", nil
}
