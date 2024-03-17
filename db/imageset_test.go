package db

import (
	"image/color"
	"os"
	"testing"

	"github.com/pokemonpower92/collagecommon/types"
)

var testDb *ImageSetDB

func TestMain(m *testing.M) {
	testDB := SetupTestISDB()
	testDb = testDB.DB
	defer testDB.TearDown()
	os.Exit(m.Run())
}

func TestCreateImageSet(t *testing.T) {
	is := &types.ImageSet{
		Name:        "test",
		Description: "test",
	}
	err := testDb.CreateImageSet(is)
	if err != nil {
		t.Errorf("Error creating imageset: %v", err)
	}
}

func TestSetAverageColors(t *testing.T) {
	aveColors := []*color.RGBA{
		{R: 0, G: 0, B: 0, A: 0},
		{R: 255, G: 255, B: 255, A: 255},
	}
	err := testDb.SetAverageColors(1, aveColors)
	if err != nil {
		t.Errorf("Error setting average colors: %v", err)
	}
}

func TestGetImageSet(t *testing.T) {
	_, err := testDb.GetImageSet(1)
	if err != nil {
		t.Errorf("Error getting imageset: %v", err)
	}
}
