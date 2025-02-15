package db

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"os"
	"strings"
	"time"
)

func generateDbConnectAdminAuthToken(creds *credentials.Credentials, region string, clusterEndpoint string) (string, error) {
	endpoint := "https://" + clusterEndpoint
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	values := req.URL.Query()
	values.Set("Action", "DbConnectAdmin")
	req.URL.RawQuery = values.Encode()

	signer := v4.Signer{
		Credentials: creds,
	}
	_, err = signer.Presign(req, nil, "dsql", region, 15*time.Minute, time.Now())
	if err != nil {
		return "", err
	}
	url := req.URL.String()[len("https://"):]
	return url, nil
}

func GetConnection(ctx context.Context, region string, clusterEndpoint string) (*pgx.Conn, error) {
	var sb strings.Builder
	sb.WriteString("postgres://")
	sb.WriteString(clusterEndpoint)
	sb.WriteString(":5432/postgres?user=admin&sslmode=verify-full")
	url := sb.String()
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		return nil, err
	}
	staticCredentials := credentials.NewStaticCredentials(
		creds.AccessKeyID,
		creds.SecretAccessKey,
		creds.SessionToken,
	)
	token, err := generateDbConnectAdminAuthToken(staticCredentials, region, clusterEndpoint)
	if err != nil {
		return nil, err
	}
	connConfig, err := pgx.ParseConfig(url)
	connConfig.Password = token
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}
	conn, err := pgx.ConnectConfig(ctx, connConfig)
	return conn, err
}

func GetPool(ctx context.Context, region string, clusterEndpoint string) (*pgxpool.Pool, error) {
	var sb strings.Builder
	sb.WriteString("postgres://")
	sb.WriteString(clusterEndpoint)
	sb.WriteString(":5432/postgres?user=admin&sslmode=verify-full")
	url := sb.String()
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		return nil, err
	}
	staticCredentials := credentials.NewStaticCredentials(
		creds.AccessKeyID,
		creds.SecretAccessKey,
		creds.SessionToken,
	)
	token, err := generateDbConnectAdminAuthToken(staticCredentials, region, clusterEndpoint)
	if err != nil {
		return nil, err
	}
	connConfig, err := pgxpool.ParseConfig(url)
	connConfig.ConnConfig.Password = token
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}
	conn, err := pgxpool.NewWithConfig(ctx, connConfig)
	return conn, err
}

func NewDBPool(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://user:password@localhost:5432/mydb?sslmode=disable"
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	return pool, nil
}
