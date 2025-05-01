package auth

import (
	"time"

	"github.com/Chinedum-Ogbuagu/guardianGo-backend.git/internal/user"
	uuid "github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key") // use env var in production

func GenerateJWT(userID uuid.UUID, role user.Role) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
