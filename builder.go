package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudspannerecosystem/memefish/ast"
	"github.com/k0kubun/pp"
	"github.com/open-policy-agent/opa/rego"
)

type pair struct{ a, d int64 }

type actionTypes struct {
	add, drop map[int64]struct{}

	replaceTable  []pair
	replaceColumn []pair
}

func newBaseCategory() *actionTypes {
	return &actionTypes{
		add:  map[int64]struct{}{},
		drop: map[int64]struct{}{},
	}
}

var mark = struct{}{}

func buildPlan(ddls []ast.DDL, prefix string) *MigrationPlan {
	info := newBaseCategory()
	aTable := NewRegoRuntime(addQuery)
	rs, err := aTable.Eval(context.Background(), rego.EvalInput(ddls))
	if err != nil {
		pp.Println(err)
		panic(1)
	}
	for _, pat := range rs {
		id := int64Id(pat.Bindings["id"])
		info.add[id] = mark
	}

	dTable := NewRegoRuntime(delQuery)
	rs, err = dTable.Eval(context.Background(), rego.EvalInput(ddls))
	if err != nil {
		pp.Println(err)
		panic(1)
	}
	for _, pat := range rs {
		id := int64Id(pat.Bindings["id"])
		info.drop[id] = mark
	}

	qTable := NewRegoRuntime(swapTableQuery)
	rs, err = qTable.Eval(context.Background(), rego.EvalInput(ddls))
	if err != nil {
		pp.Println(err)
		panic(1)
	}
	for _, pat := range rs {
		cID := int64Id(pat.Bindings["c_id"])
		dID := int64Id(pat.Bindings["d_id"])

		delete(info.add, cID)
		delete(info.drop, dID)
		info.replaceTable = append(info.replaceTable, pair{cID, dID})
	}

	qCol := NewRegoRuntime(swapColumnQuery)
	rs, err = qCol.Eval(context.Background(), rego.EvalInput(ddls))
	if err != nil {
		pp.Println(err)
		panic(1)
	}
	for _, pat := range rs {
		cID := int64Id(pat.Bindings["a_id"])
		dID := int64Id(pat.Bindings["d_id"])

		delete(info.add, cID)
		delete(info.drop, dID)
		info.replaceColumn = append(info.replaceColumn, pair{cID, dID})
	}

	result := &MigrationPlan{}
	// make prepare queries
	for id := range info.add {
		result.Prepare = append(result.Prepare, ddls[int(id)].SQL())
	}
	for id := range info.drop {
		result.Temporal = append(result.Temporal, ddls[int(id)].SQL())
	}
	for _, pair := range info.replaceTable {
		result.Temporal = append(result.Temporal, ddls[int(pair.d)].SQL())
		result.Next = append(result.Next, ddls[int(pair.a)].SQL())

		drop := ddls[int(pair.d)]
		drop.(*ast.DropTable).Name.Name = prefix + drop.(*ast.DropTable).Name.Name
		result.Cleanup = append(result.Cleanup, drop.SQL())
		create := ddls[int(pair.a)]
		create.(*ast.CreateTable).Name.Name = prefix + create.(*ast.CreateTable).Name.Name
		result.Prepare = append(result.Prepare, create.SQL())
	}

	for _, pair := range info.replaceColumn {
		result.Temporal = append(result.Temporal, ddls[int(pair.d)].SQL())
		result.Next = append(result.Next, ddls[int(pair.a)].SQL())

		drop := ddls[int(pair.d)]
		drop.(*ast.AlterTable).TableAlteration.(*ast.DropColumn).Name.Name = prefix + drop.(*ast.AlterTable).TableAlteration.(*ast.DropColumn).Name.Name
		result.Cleanup = append(result.Cleanup, drop.SQL())

		create := ddls[int(pair.a)]
		create.(*ast.AlterTable).TableAlteration.(*ast.AddColumn).Column.Name.Name = prefix + create.(*ast.AlterTable).TableAlteration.(*ast.AddColumn).Column.Name.Name
		result.Prepare = append(result.Prepare, create.SQL())
	}
	return result
}

func printDDLs(ddls []ast.DDL) {
	fmt.Println("--------")
	b, _ := json.Marshal(ddls)
	fmt.Println(string(b))
	pp.Println(ddls)
	fmt.Println("--------")
}
