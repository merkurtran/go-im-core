package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/merkurtran/go-im-core/internal/domain/user"
	"github.com/merkurtran/go-im-core/pkg/jwt"
	"github.com/merkurtran/go-im-core/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

// Service-level sentinel errors — re-exported from domain for handler consumption.
var (
	ErrUserAlreadyExists         = user.ErrUserAlreadyExists
	ErrInvalidUsernameOrPassword = user.ErrInvalidUsernameOrPassword
	ErrUserNotFound              = user.ErrUserNotFound
)

type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

func (r *UpdateUserRequest) Validate() error {
	if r.Nickname != "" && (len([]rune(r.Nickname)) < 2 || len([]rune(r.Nickname)) > 20) {
		return errors.New("nickname length must be 2-20 characters")
	}
	return nil
}

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

// TokenGenerator 抽象 JWT 生成，解耦 config 依赖
type TokenGenerator interface {
	Generate(userID string) (string, error)
}

type jwtTokenGenerator struct {
	secret string
	expire int
}

func NewJWTTokenGenerator(secret string, expireSeconds int) TokenGenerator {
	return &jwtTokenGenerator{secret: secret, expire: expireSeconds}
}

func (g *jwtTokenGenerator) Generate(userID string) (string, error) {
	return jwt.GenerateToken(userID, g.secret, g.expire)
}

type UserService struct {
	repo     user.UserRepository
	tokenGen TokenGenerator
}

func NewUserService(repo user.UserRepository, tokenGen TokenGenerator) *UserService {
	return &UserService{
		repo:     repo,
		tokenGen: tokenGen,
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
		return nil, user.ErrUserAlreadyExists
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

	token, err := s.tokenGen.Generate(newUser.ID)
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		ID:    newUser.ID,
		Token: token,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	u, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if u == nil {
		// 统一返回，防止用户枚举攻击
		return nil, user.ErrInvalidUsernameOrPassword
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, user.ErrInvalidUsernameOrPassword
	}

	token, err := s.tokenGen.Generate(u.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User: UserInfo{
			ID:       u.ID,
			Username: u.Username,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Status:   u.Status,
		},
		Token: token,
	}, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateUser only allows updating nickname and avatar
func (s *UserService) UpdateUser(ctx context.Context, id string, req *UpdateUserRequest) error {
	if id == "" {
		return errors.New("user id is required")
	}
	if err := req.Validate(); err != nil {
		return err
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return user.ErrUserNotFound
	}

	if req.Nickname != "" {
		existing.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		existing.Avatar = req.Avatar
	}

	return s.repo.Update(ctx, existing)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) UpdateUserStatus(ctx context.Context, id string, status string) error {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		slog.Error("UpdateUserStatus failed", "error", err)
		return err
	}
	if u == nil {
		return user.ErrUserNotFound
	}
	u.Status = status
	return s.repo.Update(ctx, u)
}
