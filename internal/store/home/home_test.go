package home_test

import (
	"context"
	"testing"

	"github.com/1995parham/fandogh/internal/config"
	"github.com/1995parham/fandogh/internal/db"
	"github.com/1995parham/fandogh/internal/fs"
	"github.com/1995parham/fandogh/internal/model"
	"github.com/1995parham/fandogh/internal/store/home"
	"github.com/1995parham/fandogh/internal/store/user"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"
)

type CommonHomeSuite struct {
	suite.Suite
	Store home.Home
}

func (suite *CommonHomeSuite) TestNoID() {
	require := suite.Require()

	_, err := suite.Store.Get(context.Background(), "invalid_id")
	require.Equal(home.ErrIDNotFound, err)
}

// nolint: funlen
func (suite *CommonHomeSuite) TestSetGet() {
	require := suite.Require()

	cases := []struct {
		name           string
		home           model.Home
		photos         []model.Photo
		expectedSetErr error
		expectedGetErr error
	}{
		{
			name: "Without Error",
			home: model.Home{
				ID:              "",
				Title:           "127.0.0.1",
				Location:        "Iran, Tehran",
				Description:     "Home Sweet Home",
				Peoples:         4,
				Room:            "room_type",
				Bed:             model.Double,
				Rooms:           2,
				Bathrooms:       2,
				Smoking:         false,
				Guest:           false,
				Pet:             false,
				BillsIncluded:   true,
				Contract:        "contract_type",
				SecurityDeposit: 0,
				Photos:          nil,
				Price:           0,
			},
			photos: []model.Photo{
				{
					Name:        "1.png",
					ContentType: "image/png",
					Content:     []byte{'1', '2', '3'},
				},
			},
			expectedSetErr: nil,
			expectedGetErr: nil,
		},
		{
			name: "ID Not Empty",
			home: model.Home{
				ID:              "1378",
				Title:           "127.0.0.1",
				Location:        "Iran, Tehran",
				Description:     "Home Sweet Home",
				Peoples:         4,
				Room:            "room_type",
				Bed:             model.Double,
				Rooms:           2,
				Bathrooms:       2,
				Smoking:         false,
				Guest:           false,
				Pet:             false,
				BillsIncluded:   true,
				Contract:        "contract_type",
				SecurityDeposit: 0,
				Photos:          nil,
				Price:           0,
			},
			photos: []model.Photo{
				{
					Name:        "1.png",
					ContentType: "image/png",
					Content:     []byte{'1', '2', '3'},
				},
			},
			expectedSetErr: home.ErrIDNotEmpty,
			expectedGetErr: nil,
		},
	}

	for _, c := range cases {
		c := c
		suite.Run(c.name, func() {
			require.Equal(c.expectedSetErr, suite.Store.Set(context.Background(), &c.home, c.photos))
			require.NotEmpty(c.home.ID)

			require.NotNil(c.home.Photos)
			for _, key := range c.home.Photos {
				require.NotEmpty(key)
			}

			if c.expectedSetErr == nil {
				home, err := suite.Store.Get(context.Background(), c.home.ID)
				require.Equal(c.expectedGetErr, err)
				if c.expectedGetErr == nil {
					require.Equal(c.home, home)
				}
			}
		})
	}
}

type MongoHomeSuite struct {
	DB *mongo.Database
	CommonHomeSuite
}

func (suite *MongoHomeSuite) SetupSuite() {
	cfg := config.New()

	db, err := db.New(cfg.Database)
	suite.Require().NoError(err)

	minio, err := fs.New(cfg.FileStorage)
	suite.Require().NoError(err)

	suite.DB = db
	suite.Store = home.NewMongoHome(db, minio, trace.NewNoopTracerProvider().Tracer(""))
}

func (suite *MongoHomeSuite) TearDownSuite() {
	_, err := suite.DB.Collection(user.Collection).DeleteMany(context.Background(), bson.D{})
	suite.Require().NoError(err)

	suite.Require().NoError(suite.DB.Client().Disconnect(context.Background()))
}

func TestMongoHomeSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(MongoHomeSuite))
}
