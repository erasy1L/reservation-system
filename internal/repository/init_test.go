package repository

import (
	"context"
	"fmt"
	"log"
	"room-reservation/internal/repository/postgres"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	_ "github.com/jackc/pgx/v5/stdlib"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var db *postgres.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16",
		Env: []string{
			"POSTGRES_USER=test",
			"POSTGRES_PASSWORD=test",
			"POSTGRES_DB=test",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	resource.Expire(180)

	connURL := fmt.Sprintf("postgres://test:test@%s/test?sslmode=disable", resource.GetHostPort("5432/tcp"))

	pool.MaxWait = 180 * time.Second
	if err := pool.Retry(func() error {
		db, err = postgres.New(context.Background(), connURL)
		if err != nil {
			return err
		}
		return db.Ping(context.Background())
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	startMigration(connURL)

	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}

	}()

	m.Run()
}

func startMigration(databaseUrl string) {
	migrate, err := migrate.New("file://postgres/migrations", databaseUrl)
	if err != nil {
		log.Fatalf("could not create migrate instance: %s", err)
	}

	if err = migrate.Up(); err != nil {
		log.Fatalf("could not apply migration: %s", err)
	}
}
