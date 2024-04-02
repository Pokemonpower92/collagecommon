package db_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // used by migrator
	_ "github.com/golang-migrate/migrate/v4/source/file"       // used by migrator
	_ "github.com/jackc/pgx/v4/stdlib"                         // used by migrator
	"github.com/pokemonpower92/collagecommon/db"
	"github.com/pokemonpower92/collagecommon/types"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestISDB struct {
	DB        *db.ImageSetDB
	container testcontainers.Container
}

func SetupTestISDB() *TestISDB {
	conf := types.DBConfig{
		Host:     "localhost",
		User:     "postgres",
		Password: "postgres",
		Port:     "5432",
		DbName:   "imageset",
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	container, err := createContainer(ctx, conf)
	if err != nil {
		log.Fatal("Failed to setup test imageset db container: ", err)
	}

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatal("Failed to get mapped port: ", err)
	}
	conf.Port = p.Port()

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatal("Failed to get host: ", err)
	}
	conf.Host = host

	err = migrateTestDB(conf)
	if err != nil {
		log.Fatal("Failed to migrate test db: ", err)
	}

	db, err := db.NewImageSetDB(conf)
	if err != nil {
		log.Fatal("Failed to setup test imageset db", err)
	}
	return &TestISDB{
		DB:        db,
		container: container,
	}
}

func (tdb *TestISDB) TearDown() {
	_ = tdb.container.Terminate(context.Background())
}

func createContainer(ctx context.Context, conf types.DBConfig) (testcontainers.Container, error) {
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:14-alpine",
			ExposedPorts: []string{fmt.Sprintf("%s/tcp", conf.Port)},
			Env: map[string]string{
				"POSTGRES_PASSWORD": conf.Password,
				"POSTGRES_USER":     conf.User,
				"POSTGRES_DB":       conf.DbName,
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, err
	}

	log.Println("postgres container ready and running on port: ", conf.Port)
	time.Sleep(time.Second)

	return container, nil
}

func migrateTestDB(conf types.DBConfig) error {
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get path")
	}
	pathToMigrationFiles := filepath.Dir(path) + "/migration"

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DbName,
	)
	m, err := migrate.New(fmt.Sprintf("file:%s", pathToMigrationFiles), connString)
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Println("migration done")

	return nil
}
