package lib

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

var conn *sqlx.DB

func GetDB(ctx context.Context) *sqlx.DB {
	if conn != nil {
		return conn
	}

	dbName := "datadon"
	dbUser := "postgres"
	dbHost := "datadon.czxkhjdckqtq.eu-west-1.rds.amazonaws.com"
	dbPassword := "s%%a80U2d6b!"
	dbPort := 5432
	//dbEndpoint := fmt.Sprintf("%s:%d", dbHost, dbPort)
	//region := "eu-west-1"

	//cfg, err := config.LoadDefaultConfig(ctx)
	//if err != nil {
	//	return nil, fmt.Errorf("configuration error: %w", err)
	//}

	//authenticationToken, err := auth.BuildAuthToken(
	//	ctx, dbEndpoint, region, dbUser, cfg.Credentials)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to create authentication token:  %w", err)
	//}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	conn, err := sqlx.Connect("pgx", dsn)

	if err != nil {
		panic(err)
	}

	return conn
}
