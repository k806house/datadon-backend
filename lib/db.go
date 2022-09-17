package lib

import (
	"context"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func ConnectToDB(ctx context.Context) (*sqlx.DB, error) {
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
	log.Println(err)

	return conn, err
}
