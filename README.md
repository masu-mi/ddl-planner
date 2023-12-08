DDL-planner: build plan of applying DDLs
------
I wrote this command for study.

```sh
Usage of ./ddl-planner:
  -debug
    	debug mode
  -input string
    	input file (default "./example.ddl")
  -output string
    	output directory (default "./dist")
  -prefix string
    	prefix for temporal table/column name (default "Tmp_")
  -show_steps
    	show steps in plan
```

## example

```sh
hammer diff ./example/old.ddl ./example/new.ddl > ./example/hammer-gen.ddl
```
