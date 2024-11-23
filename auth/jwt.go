package auth

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"time"

	"github.com/gitwub5/go_todo_app/clock"
	"github.com/gitwub5/go_todo_app/entity"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const (
	RoleKey     = "role"
	UserNameKey = "user_name"
)

//go:embed cert/secret.pem
var rawPrivKey []byte

//go:embed cert/public.pem
var rawPubKey []byte

// JWTer 구조체는 JWT를 생성하고 검증할 때 사용하는 키와 스토어를 포함함
type JWTer struct {
	PrivateKey, PublicKey jwk.Key       // 개인 키 및 공개 키를 JWK 형식으로 저장
	Store                 Store         // JWT와 사용자 데이터를 저장 및 로드하기 위한 인터페이스
	Clocker               clock.Clocker // 시간 관련 기능을 수행하는 clock.Clocker 객체
}

//go:generate go run github.com/matryer/moq -out moq_test.go . Store
type Store interface {
	Save(ctx context.Context, key string, userID entity.UserID) error // 사용자 ID를 특정 키와 함께 저장
	Load(ctx context.Context, key string) (entity.UserID, error)      // 특정 키를 통해 사용자 ID를 불러옴
}

// NewJWTer 함수는 JWTer 구조체를 초기화하는 생성자 함수
func NewJWTer(s Store, c clock.Clocker) (*JWTer, error) {
	j := &JWTer{Store: s}

	// rawPrivKey를 사용하여 개인 키를 JWK 형식으로 파싱
	privkey, err := parse(rawPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed in NewJWTer: private key: %w", err)
	}

	// rawPubKey를 사용하여 공개 키를 JWK 형식으로 파싱
	pubkey, err := parse(rawPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed in NewJWTer: public key: %w", err)
	}

	// PrivateKey 및 PublicKey 필드에 파싱된 키를 할당
	j.PrivateKey = privkey
	j.PublicKey = pubkey
	j.Clocker = c
	return j, nil
}

// parse 함수는 PEM 형식의 키 데이터를 JWK 형식으로 변환
func parse(rawKey []byte) (jwk.Key, error) {
	key, err := jwk.ParseKey(rawKey, jwk.WithPEM(true))
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (j *JWTer) GenerateToken(ctx context.Context, u entity.User) ([]byte, error) {
	tok, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer(`github.com/gitwub5/go_todo_app`).
		Subject("access_token").
		IssuedAt(j.Clocker.Now()).
		// Redis의 expire(만료 시간) 설정에는 아래의 링크를 사용.
		// https://pkg.go.dev/github.com/go-redis/redis/v8#Client.Set
		// clock.Duration이기 때문에 Sub(뺄셈) 메서드를 사용하여 만료 시간을 설정
		Expiration(j.Clocker.Now().Add(30*time.Minute)).
		Claim(RoleKey, u.Role).
		Claim(UserNameKey, u.Name).
		Build()
	if err != nil {
		return nil, fmt.Errorf("GenerateToken: failed to build token: %w", err)
	}
	if err := j.Store.Save(ctx, tok.JwtID(), u.ID); err != nil {
		return nil, err
	}

	// Sign a JWT!
	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256, j.PrivateKey))
	if err != nil {
		return nil, err
	}
	return signed, nil
}

func (j *JWTer) GetToken(ctx context.Context, r *http.Request) (jwt.Token, error) {
	token, err := jwt.ParseRequest(
		r,
		jwt.WithKey(jwa.RS256, j.PublicKey),
		jwt.WithValidate(false),
	)
	if err != nil {
		return nil, err
	}

	// 토큰 검증 (만료 시간 등)
	if err := jwt.Validate(token, jwt.WithClock(j.Clocker)); err != nil {
		return nil, fmt.Errorf("GetToken: failed to validate token: %w", err)
	}
	// 레디스에서 삭제해서 수동으로 expire 시키는 경우도 있다.
	if _, err := j.Store.Load(ctx, token.JwtID()); err != nil {
		return nil, fmt.Errorf("GetToken: %q expired: %w", token.JwtID(), err)
	}
	return token, nil
}

/* 애플리케이션 코드에서 매번 jwt를 생성하지 않고, context에 jwt에서 가져온 사용자 ID와 권한을 설정한다. */

type userIDKey struct{} // context에 사용자 ID를 저장하기 위한 키
type roleKey struct{}   // context에 권한을 저장하기 위한 키

// FillContext 함수는 context에 사용자 ID와 권한을 설정
func (j *JWTer) FillContext(r *http.Request) (*http.Request, error) {
	token, err := j.GetToken(r.Context(), r)
	if err != nil {
		return nil, err
	}
	uid, err := j.Store.Load(r.Context(), token.JwtID())
	if err != nil {
		return nil, err
	}
	ctx := SetUserID(r.Context(), uid)

	ctx = SetRole(ctx, token)
	clone := r.Clone(ctx)
	return clone, nil
}

// SetUserID 함수는 context에 사용자 ID를 설정
func SetUserID(ctx context.Context, uid entity.UserID) context.Context {
	return context.WithValue(ctx, userIDKey{}, uid)
}

// GetUserID 함수는 context에서 사용자 ID를 가져옴
func GetUserID(ctx context.Context) (entity.UserID, bool) {
	id, ok := ctx.Value(userIDKey{}).(entity.UserID)
	return id, ok
}

// SetRole 함수는 context에 권한을 설정
func SetRole(ctx context.Context, tok jwt.Token) context.Context {
	get, ok := tok.Get(RoleKey)
	if !ok {
		return context.WithValue(ctx, roleKey{}, "")
	}
	return context.WithValue(ctx, roleKey{}, get)
}

// GetRole 함수는 context에서 권한을 가져옴
func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(roleKey{}).(string)
	return role, ok
}

// IsAdmin 함수는 context에서 권한을 가져와서 admin인지 확인
func IsAdmin(ctx context.Context) bool {
	role, ok := GetRole(ctx)
	if !ok {
		return false
	}
	return role == "admin"
}
