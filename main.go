package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/k0kubun/sqldef"
	"github.com/k0kubun/sqldef/database"
	"github.com/k0kubun/sqldef/database/mysql"
	"github.com/k0kubun/sqldef/parser"
	"github.com/k0kubun/sqldef/schema"
)

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context) (string, error) {
	schema := "./schema/ddl.sql"
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	err = migrate(os.Getenv("DB_HOST"), port, os.Getenv("DB_USER"), os.Getenv("PASSWORD"), os.Getenv("DB_NAME"), schema)
	if err != nil {
		fmt.Print(err)
		return "", err
	}
	return "success", nil
}

func migrate(host string, port int, user string, password string, dbname string, schemaFile string) error {
	db, err := mysql.NewDatabase(database.Config{
		DbName:   dbname,
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
	})
	if err != nil {
		return fmt.Errorf("failed to create a database adapter: %w", err)
	}
	sqlParser := database.NewParser(parser.ParserModeMysql)
	desiredDDLs, err := sqldef.ReadFile(schemaFile)
	fmt.Printf(desiredDDLs)
	if err != nil {
		return fmt.Errorf("Failed to read %s: %w", schemaFile, err)
	}
	options := &sqldef.Options{DesiredDDLs: desiredDDLs}
	sqldef.Run(schema.GeneratorModeMysql, db, sqlParser, options)

	return nil
}
