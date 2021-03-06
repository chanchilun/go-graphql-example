package schema

import (
	"github.com/OscarYuen/go-graphql-example/model"
)

func (r *Resolver) CreateUser(args *struct {
	Email    string
	Password string
}) (*userResolver, error) {
	user := &model.User{
		Email:    args.Email,
		Password: args.Password,
	}
	user.HashedPassword()
	if err := DB.Create(user).Error; err != nil {
		return nil, err
	}
	return &userResolver{user}, nil
}
