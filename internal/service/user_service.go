package service

import (
	"context"
	"errors"

	"github.com/merkurtran/go-im-core/internal/config"
	"github.com/merkurtran/go-im-core/internal/domain/user"
	"github.com/merkurtran/go-im-core/pkg/jwt"
	"github.com/merkurtran/go-im-core/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type RegisterResponse struct {
	ID    string `json:"user_id"`
	Token string `json:"token"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Status   string `json:"status"`
}

type LoginResponse struct {
	User  UserInfo `json:"user"`
	Token string   `json:"token"`
}

type UserService struct {
	repo      user.UserRepository
	jwtConfig *config.JWTConfig
}

func NewUserService(repo user.UserRepository, jwtConfig *config.JWTConfig) *UserService {
	return &UserService{
		repo:      repo,
		jwtConfig: jwtConfig,
	}
}

func (s *UserService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	if err := validator.ValidateUsername(req.Username); err != nil {
		return nil, err
	}
	if err := validator.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	existing, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &user.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Status:   "offline",
		Avatar:   "https://gravatar.com/avatar/dba6bae8c566f9d4041fb9cd9ada7741?d=identicon&f=y",
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	token, err := jwt.GenerateToken(newUser.ID, s.jwtConfig.Secret, s.jwtConfig.Expire)
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		ID:    newUser.ID,
		Token: token,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	//
	user, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	token, err := jwt.GenerateToken(user.ID, s.jwtConfig.Secret, s.jwtConfig.Expire)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User: UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Status:   user.Status,
		},
		Token: token,
	}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, nickname, avatar string) error {
	if id == "" {
		return errors.New("user id is required")
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	existing.Nickname = nickname
	existing.Avatar = avatar

	return s.repo.Update(ctx, existing)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
