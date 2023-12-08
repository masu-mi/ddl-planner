package main

import (
	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/token"
)

func GetSQLParser(name, sql string) *memefish.Parser {
	return &memefish.Parser{
		Lexer: &memefish.Lexer{File: &token.File{
			FilePath: name,
			Buffer:   sql,
		}},
	}
}
