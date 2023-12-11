package main

import (
	"context"

	"github.com/k0kubun/pp"
	"github.com/open-policy-agent/opa/rego"

	_ "embed"
)

//go:embed module.rego
var module string

func NewRegoRuntime(q func(*rego.Rego)) rego.PreparedEvalQuery {
	r := rego.New(q, rego.Module("data.ddl", module))
	p, e := r.PrepareForEval(context.Background())
	if e != nil {
		pp.Println(e)
		panic(1)
	}
	return p
}

var addQuery = rego.Query(`data.ddl.add[id]`)
var delQuery = rego.Query(`data.ddl.del[id]`)

var replaceTableQuery = rego.Query(`data.ddl.replace_table_queries[c_id][d_id]`)

var swapColumnQuery = rego.Query(`data.ddl.alter_add_col[a_id]
data.ddl.alter_drop_col[d_id]
input[a_id].Name.Name == input[d_id].Name.Name
input[a_id].TableAlteration.Column.Name.Name == input[d_id].TableAlteration.Name.Name`)
