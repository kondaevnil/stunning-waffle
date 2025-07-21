package domain

import (
	"testing"
	"vk/ecom/internal/domain"
	"vk/ecom/internal/pkg/jwt"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	t.Run("should generate valid token for user", func(t *testing.T) {
		user := &domain.User{
			ID:    1,
			Login: "testuser",
		}

		token, err := jwt.GenerateToken(user)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.IsType(t, "", token)
	})

	t.Run("should generate different tokens for different users", func(t *testing.T) {
		user1 := &domain.User{
			ID:    1,
			Login: "testuser1",
		}
		user2 := &domain.User{
			ID:    2,
			Login: "testuser2",
		}

		token1, err1 := jwt.GenerateToken(user1)
		token2, err2 := jwt.GenerateToken(user2)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, token1, token2)
	})
}

func TestParseToken(t *testing.T) {
	t.Run("should parse valid token correctly", func(t *testing.T) {
		originalUser := &domain.User{
			ID:    1,
			Login: "testuser",
		}

		token, err := jwt.GenerateToken(originalUser)
		assert.NoError(t, err)

		parsedUser, err := jwt.ParseToken(token)

		assert.NoError(t, err)
		assert.NotNil(t, parsedUser)
		assert.Equal(t, originalUser.ID, parsedUser.ID)
		assert.Equal(t, originalUser.Login, parsedUser.Login)
	})

	t.Run("should fail with invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.string"

		parsedUser, err := jwt.ParseToken(invalidToken)

		assert.Error(t, err)
		assert.Nil(t, parsedUser)
	})

	t.Run("should fail with empty token", func(t *testing.T) {
		parsedUser, err := jwt.ParseToken("")

		assert.Error(t, err)
		assert.Nil(t, parsedUser)
	})

	t.Run("should fail with malformed token", func(t *testing.T) {
		malformedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.malformed"

		parsedUser, err := jwt.ParseToken(malformedToken)

		assert.Error(t, err)
		assert.Nil(t, parsedUser)
	})
}

func TestTokenRoundTrip(t *testing.T) {
	t.Run("should maintain user data through generate and parse cycle", func(t *testing.T) {
		testCases := []struct {
			name string
			user *domain.User
		}{
			{
				name: "regular user",
				user: &domain.User{ID: 1, Login: "testuser"},
			},
			{
				name: "user with long login",
				user: &domain.User{ID: 999, Login: "very_long_username_with_underscores"},
			},
			{
				name: "user with special characters in login",
				user: &domain.User{ID: 42, Login: "user@example.com"},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				token, err := jwt.GenerateToken(tc.user)
				assert.NoError(t, err)

				parsedUser, err := jwt.ParseToken(token)
				assert.NoError(t, err)
				assert.Equal(t, tc.user.ID, parsedUser.ID)
				assert.Equal(t, tc.user.Login, parsedUser.Login)
			})
		}
	})
}
