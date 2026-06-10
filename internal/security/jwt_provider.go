package security

import (
	"context"
	"fmt"
	"sandbox-api-gin/internal/domain/apperror"
	"sandbox-api-gin/internal/domain/model"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type JwtProvider struct {
	cache            *jwk.Cache
	jwksURL          string
	issuer           string
	allowedAudiences []string
}

func NewJwtProvider(issuer string, audiences []string) (*JwtProvider, error) {
	ctx := context.Background()
	cache := jwk.NewCache(ctx)

	jwksURL := issuer + "/.well-known/jwks.json"

	if err := cache.Register(jwksURL, jwk.WithMinRefreshInterval(15*time.Minute)); err != nil {
		return nil, fmt.Errorf("JWKS登録失敗: %w", err)
	}

	// 起動時に初回取得
	if _, err := cache.Refresh(ctx, jwksURL); err != nil {
		return nil, fmt.Errorf("JWKS初回取得失敗: %w", err)
	}

	return &JwtProvider{
		cache:            cache,
		jwksURL:          jwksURL,
		issuer:           issuer,
		allowedAudiences: audiences,
	}, nil
}

// Parse はJWTを検証してAuthUserを返す。検証失敗時はAuthenticationErrorを返す。
func (p *JwtProvider) Parse(tokenStr string) (*model.AuthUser, error) {
	ctx := context.Background()

	keyset, err := p.cache.Get(ctx, p.jwksURL)
	if err != nil {
		return nil, apperror.NewAuthenticationError("JWKSの取得に失敗しました: " + err.Error())
	}

	tok, err := jwt.Parse(
		[]byte(tokenStr),
		jwt.WithKeySet(keyset),
		jwt.WithValidate(true),
		jwt.WithIssuer(p.issuer),
	)
	if err != nil {
		return nil, apperror.NewAuthenticationError("JWT検証失敗: " + err.Error())
	}

	// audience検証（許可リストのいずれかと一致すればOK）
	tokenAudiences := tok.Audience()
	found := false
outer:
	for _, ta := range tokenAudiences {
		for _, allowed := range p.allowedAudiences {
			if ta == allowed {
				found = true
				break outer
			}
		}
	}
	if !found {
		return nil, apperror.NewAuthenticationError("JWT audienceが不正です")
	}

	sub := tok.Subject()
	emailVal, ok := tok.Get("email")
	if !ok {
		return nil, apperror.NewAuthenticationError("JWTの必須クレームが不足しています")
	}
	email, ok := emailVal.(string)
	if !ok || email == "" {
		return nil, apperror.NewAuthenticationError("JWTの必須クレームが不足しています")
	}

	emailVerifiedVal, _ := tok.Get("email_verified")
	emailVerified, _ := emailVerifiedVal.(bool)

	if sub == "" {
		return nil, apperror.NewAuthenticationError("JWTの必須クレームが不足しています")
	}

	return &model.AuthUser{
		Sub:           sub,
		Email:         email,
		EmailVerified: emailVerified,
		Admin:         false,
	}, nil
}
