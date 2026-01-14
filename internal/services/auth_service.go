package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/ScriptVandal/backend-go/internal/config"
	"github.com/ScriptVandal/backend-go/internal/models"
	"github.com/ScriptVandal/backend-go/internal/repositories"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

type AuthService struct {
	userRepo         repositories.UserRepository
	refreshTokenRepo repositories.RefreshTokenRepository
	config           *config.Config
}

func NewAuthService(userRepo repositories.UserRepository, refreshTokenRepo repositories.RefreshTokenRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		config:           cfg,
	}
}

// Register creates a new user
func (s *AuthService) Register(email, password string) (*models.User, error) {
	// Check if user already exists
	existing, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	passwordHash, err := s.hashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		ID:           generateID(),
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(email, password string) (*models.User, string, string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", "", err
	}
	if user == nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	// Verify password
	if !s.verifyPassword(password, user.PasswordHash) {
		return nil, "", "", errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, jti, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	// Store refresh token
	token := &models.RefreshToken{
		JTI:       jti,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(s.config.RefreshTTL),
		CreatedAt: time.Now(),
	}
	if err := s.refreshTokenRepo.Create(token); err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

// Refresh generates new access token from refresh token
func (s *AuthService) Refresh(refreshTokenString string) (string, error) {
	// Parse refresh token
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.config.JWTRefreshSecret), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid user ID in token")
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return "", errors.New("invalid JTI in token")
	}

	// Check if token is revoked
	storedToken, err := s.refreshTokenRepo.GetByJTI(jti)
	if err != nil {
		return "", err
	}
	if storedToken == nil {
		return "", errors.New("token not found")
	}
	if storedToken.RevokedAt != nil {
		return "", errors.New("token has been revoked")
	}
	if time.Now().After(storedToken.ExpiresAt) {
		return "", errors.New("token has expired")
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(userID)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// Logout revokes the refresh token
func (s *AuthService) Logout(refreshTokenString string) error {
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.config.JWTRefreshSecret), nil
	})

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return errors.New("invalid JTI in token")
	}

	return s.refreshTokenRepo.Revoke(jti)
}

// ValidateAccessToken validates an access token and returns user ID
func (s *AuthService) ValidateAccessToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid access token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid user ID in token")
	}

	return userID, nil
}

// Helper functions

func (s *AuthService) hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	
	// Encode salt and hash together
	encoded := base64.StdEncoding.EncodeToString(salt) + ":" + base64.StdEncoding.EncodeToString(hash)
	return encoded, nil
}

func (s *AuthService) verifyPassword(password, encoded string) bool {
	parts := splitPasswordHash(encoded)
	if len(parts) != 2 {
		return false
	}

	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	hash, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	testHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return subtle(hash, testHash)
}

func (s *AuthService) generateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(s.config.AccessTTL).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

func (s *AuthService) generateRefreshToken(userID string) (string, string, error) {
	jti := generateID()
	
	claims := jwt.MapClaims{
		"sub": userID,
		"jti": jti,
		"exp": time.Now().Add(s.config.RefreshTTL).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWTRefreshSecret))
	return tokenString, jti, err
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func splitPasswordHash(encoded string) []string {
	result := []string{}
	current := ""
	for _, c := range encoded {
		if c == ':' {
			result = append(result, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func subtle(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}
