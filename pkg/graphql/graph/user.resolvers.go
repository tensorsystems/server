package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/Tensor-Systems/tensoremr-server/pkg/graphql/graph/model"
	"github.com/Tensor-Systems/tensoremr-server/pkg/jwt"
	"github.com/Tensor-Systems/tensoremr-server/pkg/middleware"
	"github.com/Tensor-Systems/tensoremr-server/pkg/repository"
	deepCopy "github.com/ulule/deepcopier"
)

func (r *mutationResolver) Signup(ctx context.Context, input model.UserInput) (*repository.User, error) {
	var entity repository.User
	deepCopy.Copy(&input).To(&entity)

	entity.Active = true

	if err := entity.HashPassword(); err != nil {
		return nil, err
	}

	if input.Signature != nil {
		fileName, hashedFileName, hash, ext := HashFileName(input.Signature.Name)
		err := WriteFile(input.Signature.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		entity.Signature = &repository.File{
			ContentType: input.Signature.File.ContentType,
			Size:        input.Signature.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		}
	}

	if input.ProfilePic != nil {
		fileName, hashedFileName, hash, ext := HashFileName(input.ProfilePic.Name)
		err := WriteFile(input.ProfilePic.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		entity.ProfilePic = &repository.File{
			ContentType: input.ProfilePic.File.ContentType,
			Size:        input.ProfilePic.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		}
	}

	var userType repository.UserType
	userTypes, err := userType.GetByIds(input.UserTypeIds)
	if err != nil {
		return nil, err
	}

	if err := entity.Save(userTypes); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (string, error) {
	var user repository.User

	// Check if user exists
	err := user.GetByEmail(input.Email)
	if err != nil {
		return "", err
	}

	// Check password validity
	pErr := user.CheckPassword(input.Password)
	if pErr != nil {
		return "", pErr
	}

	// Generate JWT Token
	jwtWrapper := jwt.Wrapper{
		SecretKey:       r.Config.JwtSecret,
		Issuer:          r.Config.JwtIssuer,
		ExpirationHours: 24,
	}

	token, err := jwtWrapper.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *mutationResolver) ResetPassword(ctx context.Context, id int) (*repository.User, error) {
	var entity repository.User
	if err := entity.Get(id); err != nil {
		return nil, err
	}

	entity.Password = "changeme"

	if err := entity.HashPassword(); err != nil {
		return nil, err
	}

	if err := entity.Update(nil); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UserUpdateInput) (*repository.User, error) {
	var entity repository.User
	deepCopy.Copy(&input).To(&entity)

	if input.Signature != nil {
		fileName, hashedFileName, hash, ext := HashFileName(input.Signature.Name)
		err := WriteFile(input.Signature.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		entity.Signature = &repository.File{
			ContentType: input.Signature.File.ContentType,
			Size:        input.Signature.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		}
	}

	if input.ProfilePic != nil {
		fileName, hashedFileName, hash, ext := HashFileName(input.ProfilePic.Name)
		err := WriteFile(input.ProfilePic.File.File, hashedFileName+"."+ext)
		if err != nil {
			return nil, err
		}

		entity.ProfilePic = &repository.File{
			ContentType: input.ProfilePic.File.ContentType,
			Size:        input.ProfilePic.File.Size,
			FileName:    fileName,
			Extension:   ext,
			Hash:        hash,
		}
	}

	var userType repository.UserType
	userTypes, err := userType.GetByIds(input.UserTypeIds)
	if err != nil {
		return nil, err
	}

	if err := entity.Update(userTypes); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) ChangePassword(ctx context.Context, input model.ChangePasswordInput) (*repository.User, error) {
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var user repository.User
	err = user.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	// Check password validity
	if user.CheckPassword(input.PreviousPassword) != nil {
		return nil, err
	}

	// Check password confirmation
	if input.Password != input.ConfirmPassword {
		return nil, errors.New("Passwords do no match")
	}

	user.Password = input.Password
	if user.HashPassword() != nil {
		return nil, err
	}

	if err := user.Update(nil); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *mutationResolver) SaveUserType(ctx context.Context, input model.UserTypeInput) (*repository.UserType, error) {
	var entity repository.UserType
	deepCopy.Copy(&input).To(&entity)

	err := entity.Save()
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *mutationResolver) UpdateUserType(ctx context.Context, input model.UserTypeUpdateInput) (*repository.UserType, error) {
	var userType repository.UserType

	deepCopy.Copy(&input).To(&userType)

	_, err := userType.Update()
	if err != nil {
		return nil, err
	}

	return &userType, nil
}

func (r *mutationResolver) DeleteUserType(ctx context.Context, id int) (bool, error) {
	var userType repository.UserType
	err := userType.Delete(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) User(ctx context.Context, id int) (*repository.User, error) {
	var entity repository.User
	err := entity.Get(id)
	if err != nil {
		return nil, err
	}

	return &entity, err
}

func (r *queryResolver) Users(ctx context.Context, page repository.PaginationInput, filter *model.UserFilter, searchTerm *string) (*model.UserConnection, error) {
	var f repository.User
	if filter != nil {
		deepCopy.Copy(filter).To(&f)
	}

	var entity repository.User
	entities, count, err := entity.Search(page, &f, searchTerm)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.UserEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.UserEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.UserConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) UserTypes(ctx context.Context, page repository.PaginationInput) (*model.UserTypeConnection, error) {
	var entity repository.UserType
	entities, count, err := entity.GetAll(page)

	if err != nil {
		return nil, err
	}

	edges := make([]*model.UserTypeEdge, len(entities))

	for i, entity := range entities {
		e := entity

		edges[i] = &model.UserTypeEdge{
			Node: &e,
		}
	}

	pageInfo, totalCount := GetPageInfo(entities, count, page)
	return &model.UserTypeConnection{PageInfo: pageInfo, Edges: edges, TotalCount: totalCount}, nil
}

func (r *queryResolver) SearchUsers(ctx context.Context, input model.UserSearchInput) ([]*repository.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetByUserTypeTitle(ctx context.Context, input string) ([]*repository.User, error) {
	var user repository.User

	users, err := user.GetByUserTypeTitle(input)
	if err != nil {
		return nil, err
	}

	return users, nil
}
