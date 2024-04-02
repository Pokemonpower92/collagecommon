package db

import (
	"context"
	"fmt"
	"image/color"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pokemonpower92/collagecommon/types"
)

type ImageSetDB struct {
	client *pgxpool.Pool
	l      *log.Logger
}

func NewImageSetDB(conf DBConfig) (*ImageSetDB, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DbName,
	)
	client, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	return &ImageSetDB{
		client: client,
		l:      log.New(log.Writer(), "imageset-db ", log.LstdFlags),
	}, nil
}

func (isdb *ImageSetDB) GetImageSet(id int) (*types.ImageSet, error) {
	imCtx, imCancel := context.WithCancel(context.Background())
	defer imCancel()

	avCtx, avCancel := context.WithCancel(context.Background())
	defer avCancel()

	transaction, err := isdb.client.Begin(imCtx)
	if err != nil {
		isdb.l.Printf("Error beginning transaction: %v\n", err)
		return nil, err
	}

	is := &types.ImageSet{}

	imageSet := transaction.QueryRow(imCtx,
		"SELECT * FROM imagesets WHERE id = $1",
		id,
	)

	err = imageSet.Scan(&is.ID, &is.Name, &is.Description)
	if err != nil {
		isdb.l.Printf("Error scanning imageSet: %v\n", err)
		return nil, err
	}

	err = transaction.Commit(context.Background())
	if err != nil {
		isdb.l.Printf("Error committing transaction: %v\n", err)
		return nil, err
	}

	transaction, err = isdb.client.Begin(avCtx)
	if err != nil {
		isdb.l.Printf("Error beginning transaction: %v\n", err)
		return nil, err
	}

	ave_colors, err := transaction.Query(avCtx,
		"SELECT r, g, b, a FROM average_colors WHERE imageset_id = $1",
		id,
	)
	if err != nil {
		isdb.l.Printf("Error querying average colors: %v\n", err)
		return nil, err
	}

	for ave_colors.Next() {
		var r, g, b, a int
		err = ave_colors.Scan(&r, &g, &b, &a)
		if err != nil {
			isdb.l.Printf("Error scanning ave_colors: %v\n", err)
			return nil, err
		}
		is.AverageColors = append(is.AverageColors, &color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
	}

	err = transaction.Commit(context.Background())
	if err != nil {
		isdb.l.Printf("Error committing transaction: %v\n", err)
		return nil, err
	}

	return is, nil
}

func (isdb *ImageSetDB) CreateImageSet(is *types.ImageSet) error {
	transaction, err := isdb.client.Begin(context.Background())
	if err != nil {
		isdb.l.Printf("Error beginning transaction: %v\n", err)
		return err
	}

	_, err = transaction.Exec(
		context.Background(),
		"INSERT INTO imagesets (name, description) VALUES ($1, $2);",
		is.Name,
		is.Description,
	)
	if err != nil {
		return err
	}

	err = transaction.Commit(context.Background())
	if err != nil {
		isdb.l.Printf("Error committing transaction: %v\n", err)
		return err
	}

	return nil
}

func (isdb *ImageSetDB) SetAverageColors(id int, averageColors []*color.RGBA) error {
	transaction, err := isdb.client.Begin(context.Background())
	if err != nil {
		isdb.l.Printf("Error beginning transaction: %v\n", err)
		return err
	}

	for _, color := range averageColors {
		_, err = transaction.Exec(
			context.Background(),
			"INSERT INTO average_colors (imageset_id, r, g, b, a) VALUES ($1, $2, $3, $4, $5);",
			id,
			color.R,
			color.G,
			color.B,
			color.A,
		)
		if err != nil {
			isdb.l.Printf("Error inserting average color: %v\n", err)
			return err
		}
	}

	err = transaction.Commit(context.Background())
	if err != nil {
		isdb.l.Printf("Error committing transaction: %v\n", err)
		return err
	}

	return nil
}
