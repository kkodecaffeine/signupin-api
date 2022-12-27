package persistence

import (
	"signupin-api/internal/app/api/dto"
	"signupin-api/internal/pkg/user"

	"github.com/kamva/mgm/v3"

	"github.com/kkodecaffeine/go-common/core/database/mongo/errortype"
	"github.com/kkodecaffeine/go-common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepo struct {
	client *mongo.Client
	mapper entityMapper
}

var _ user.Repository = &userRepo{}

func (r *userRepo) SaveOne(model *user.User) (string, error) {
	coll := mgm.Coll(model)
	err := coll.Create(model)
	if err != nil {
		return "", errortype.ParseAndReturnDBError(err, coll.Name(), nil, nil, nil)
	}

	insertedID := utils.MapToStringID(model.ID)
	return insertedID, nil
}

func (r *userRepo) GetOne(email string, password ...string) (*dto.PostSignUpResponse, error) {
	found := &user.User{}
	var filter primitive.M

	if len(password) == 0 {
		filter = bson.M{"email": email}
	} else {
		filter = bson.M{
			"$and": []bson.M{
				{"email": email},
				{"password": password[0]},
			},
		}
	}

	err := mgm.Coll(found).FindOne(mgm.Ctx(), filter).Decode(found)
	if err != nil {
		return nil, errortype.ParseAndReturnDBError(err, mgm.CollName(found), filter, nil, nil)
	}

	result := r.mapper.toDomainProps(found.ID, found)

	return result, nil
}

func (r *userRepo) GetOneByID(ID string) (*dto.PostSignUpResponse, error) {
	objectID, err := utils.MapToObjectID(ID)
	if err != nil {
		return nil, err
	}

	found := &user.User{}
	filter := bson.M{"_id": objectID}

	coll := mgm.Coll(found)
	err = coll.FindOne(mgm.Ctx(), filter).Decode(found)
	if err != nil {
		return nil, errortype.ParseAndReturnDBError(err, coll.Name(), filter, nil, nil)
	}

	result := r.mapper.toDomainProps(found.ID, found)

	return result, nil
}

func New(client *mongo.Client) user.Repository {
	return &userRepo{client, entityMapper{}}
}
