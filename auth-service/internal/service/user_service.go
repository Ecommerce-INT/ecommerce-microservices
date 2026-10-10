package service

import (
	"context"
	"errors"
	"fmt"
	"math"

	"com.ecommerce/auth-service/internal/model"
	"com.ecommerce/auth-service/internal/repository"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUsernameExists = errors.New("username already exists")
	ErrEmailExists    = errors.New("email already exists")
	ErrPhoneExists    = errors.New("phone already exists")
	ErrRoleNotFound   = errors.New("role not found")
	ErrRoleConflict   = errors.New("role assignment conflict")
)

type UserService struct {
	repo           *repository.UserRepository
	keycloakClient *KeycloakClient
}

func NewUserService(repo *repository.UserRepository, kc *KeycloakClient) *UserService {
	return &UserService{
		repo:           repo,
		keycloakClient: kc,
	}
}

func (s *UserService) Register(ctx context.Context, req model.RegisterRequest) (*model.User, error) {
	exists, err := s.repo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("%w: %s", ErrUsernameExists, req.Username)
	}

	exists, err = s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("%w: %s", ErrEmailExists, req.Email)
	}

	if req.Phone != "" {
		exists, err = s.repo.ExistsByPhone(ctx, req.Phone)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("%w: %s", ErrPhoneExists, req.Phone)
		}
	}

	roles := req.Roles
	if len(roles) == 0 {
		roles = []string{"USER"}
	}

	keycloakUserID, err := s.keycloakClient.CreateUser(ctx, req.Username, req.Email, req.FullName, req.Password, roles)
	if err != nil {
		return nil, err
	}

	u := &model.User{
		FullName:       req.FullName,
		Username:       req.Username,
		Email:          req.Email,
		Gender:         req.Gender,
		Phone:          req.Phone,
		Avatar:         req.Avatar,
		KeycloakUserID: keycloakUserID,
	}

	created, err := s.repo.Create(ctx, u, roles)
	if err != nil {
		_ = s.keycloakClient.DeleteUser(ctx, keycloakUserID)
		return nil, err
	}

	return created, nil
}

func (s *UserService) Update(ctx context.Context, id int64, req model.UpdateUserRequest) (*model.UserResponse, error) {
	updated, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}
	if updated == nil {
		return nil, ErrUserNotFound
	}
	return toUserResponse(updated), nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	user, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if user.KeycloakUserID != "" {
		_ = s.keycloakClient.DeleteUser(ctx, user.KeycloakUserID)
	}

	return s.repo.Delete(ctx, id)
}

func (s *UserService) FindById(ctx context.Context, id int64) (*model.UserResponse, error) {
	user, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return toUserResponse(user), nil
}

func (s *UserService) FindByUsername(ctx context.Context, username string) (*model.UserResponse, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return toUserResponse(user), nil
}

func (s *UserService) FindAllUsers(ctx context.Context, page, size int, sortBy, sortOrder string) (*model.PageResponse[model.UserResponse], error) {
	if size <= 0 {
		size = 10
	}
	if page < 0 {
		page = 0
	}

	users, total, err := s.repo.FindAll(ctx, page, size, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	var content []model.UserResponse
	for _, u := range users {
		content = append(content, *toUserResponse(&u))
	}

	totalPages := int(math.Ceil(float64(total) / float64(size)))
	if totalPages == 0 && total == 0 {
		totalPages = 0
	}

	return &model.PageResponse[model.UserResponse]{
		Content:          content,
		TotalElements:    total,
		TotalPages:       totalPages,
		Size:             size,
		Number:           page,
		NumberOfElements: len(content),
		First:            page == 0,
		Last:             page >= totalPages-1,
		Empty:            len(content) == 0,
	}, nil
}

func (s *UserService) AssignRole(ctx context.Context, userId int64, roleName string) (bool, error) {
	return s.repo.AssignRole(ctx, userId, roleName)
}

func (s *UserService) RevokeRole(ctx context.Context, userId int64, roleName string) (bool, error) {
	return s.repo.RevokeRole(ctx, userId, roleName)
}

func (s *UserService) GetUserRoles(ctx context.Context, userId int64) ([]string, error) {
	return s.repo.GetUserRoles(ctx, userId)
}

func toUserResponse(u *model.User) *model.UserResponse {
	return &model.UserResponse{
		ID:       u.ID,
		FullName: u.FullName,
		Username: u.Username,
		Email:    u.Email,
		Gender:   u.Gender,
		Phone:    u.Phone,
		Avatar:   u.Avatar,
	}
}
