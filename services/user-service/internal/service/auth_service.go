package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/config"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/domain"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/repository"
)

type IAuthService interface {
	ValidateToken(ctx context.Context, token string) (*domain.ValidateTokenResponse, error)
	SyncUser(ctx context.Context, keycloakID string) (*domain.User, error)
}

type authService struct {
	userRepo repository.IUserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.IUserRepository, cfg *config.Config) IAuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

type keycloakIntrospectResponse struct {
	Active bool   `json:"active"`
	Sub    string `json:"sub"`
	Email  string `json:"email"`
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*domain.ValidateTokenResponse, error) {
	keycloak := s.cfg.Keycloak
	introspectURL := keycloak.Url + "/realms/" + keycloak.Realm + "/protocol/openid-connect/token/introspect"

	data := url.Values{}
	data.Set("client_id", keycloak.ClientID)
	data.Set("client_secret", keycloak.ClientSecret)
	data.Set("token", token)

	req, err := http.NewRequestWithContext(ctx, "POST", introspectURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, domain.ErrKeycloakUnreachable
	}
	defer response.Body.Close()

	var introspect keycloakIntrospectResponse
	if err := json.NewDecoder(response.Body).Decode(&introspect); err != nil {
		return nil, err
	}

	if !introspect.Active {
		return nil, domain.ErrInvalidToken
	}

	user, err := s.SyncUser(ctx, introspect.Sub)
	if err != nil {
		return nil, err
	}

	return &domain.ValidateTokenResponse{
		UserId:     user.Id,
		Email:      user.Email,
		Role:       user.Role,
		KeycloakId: user.KeycloakId,
	}, nil
}

func (s *authService) SyncUser(ctx context.Context, keycloakID string) (*domain.User, error) {
	user, err := s.userRepo.FindByKeycloakID(ctx, keycloakID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	// kullanıcı yoksa oluştur
	newUser := &domain.User{
		KeycloakId: keycloakID,
		Role:       domain.RoleCustomer,
		IsActive:   true,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
