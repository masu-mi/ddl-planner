package main

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"path"

	"github.com/k0kubun/pp"
)

func main() {
	input := flag.String("input", "./input.ddl", "input file")
	output := flag.String("output", "./dist", "output directory")
	prefix := flag.String("prefix", "Tmp_", "prefix for temporal table/column name")
	debug := flag.Bool("debug", false, "debug mode")
	showSteps := flag.Bool("show_steps", false, "show steps in plan")
	flag.Parse()

	f, e := os.Open(*input)
	if e != nil {
		pp.Println(e)
		return
	}

	buf, e := io.ReadAll(f)
	if e != nil {
		pp.Println(e)
		return
	}
	ddlParser := GetSQLParser("input.ddl", string(buf))
	ddls, e := ddlParser.ParseDDLs()
	if e != nil {
		pp.Println(e)
		return
	}
	if *debug {
		printDDLs(ddls)
	}
	plan := buildPlan(ddls, *prefix)
	if *debug || *showSteps {
		pp.Println(plan)
	}
	plan.Generate(*output)
}

func int64Id(v interface{}) int64 {
	r, _ := v.(json.Number).Int64()
	return r
}

type MigrationPlan struct {
	Prepare  []string
	Temporal []string
	Next     []string
	Cleanup  []string
}

func (p *MigrationPlan) Generate(dir string) error {
	f, _ := os.Create(path.Join(dir, "0_prepare.ddl"))
	for _, s := range p.Prepare {
		f.WriteString(s + ";\n")
	}
	f.Close()

	f, _ = os.Create(path.Join(dir, "1_recreate.ddl"))
	for _, s := range p.Temporal {
		f.WriteString(s + ";\n")
	}
	for _, s := range p.Next {
		f.WriteString(s + ";\n")
	}
	f.Close()

	f, _ = os.Create(path.Join(dir, "2_cleanup.ddl"))
	for _, s := range p.Cleanup {
		f.WriteString(s + ";\n")
	}
	f.Close()

	return nil
}
