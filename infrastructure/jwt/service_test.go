package jwt_test

import (
	"testing"
	"time"

	ihash_mocks "github.com/beka-birhanu/task_manager_final/domain/i_hash/mocks"
	usermodel "github.com/beka-birhanu/task_manager_final/domain/models/user"
	"github.com/beka-birhanu/task_manager_final/infrastructure/jwt"

	jwt_builtin "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/suite"
)

type JWTServiceSuite struct {
	suite.Suite
	jwtService     *jwt.Service
	user           *usermodel.User
	mockHasher     *ihash_mocks.Service
	secretKey      string
	issuer         string
	expirationTime time.Duration
	username       string
	password       string
	hashedPassword string
}

func (suite *JWTServiceSuite) SetupTest() {
	// Initialize property values
	suite.secretKey = "testsecretkey"
	suite.issuer = "testissuer"
	suite.expirationTime = time.Hour
	suite.username = "valid_user"
	suite.password = "&&^_Str0ngP@ssw0rd!@D$"
	suite.hashedPassword = "hashedStrongPassword"

	// Setup mock hasher
	suite.mockHasher = new(ihash_mocks.Service)
	suite.mockHasher.On("Hash", suite.password).Return(suite.hashedPassword, nil)

	// Create a new user using the mocked hasher
	var err error
	suite.user, err = usermodel.New(usermodel.Config{
		Username:       suite.username,
		PlainPassword:  suite.password,
		IsAdmin:        false,
		PasswordHasher: suite.mockHasher,
	})
	suite.NoError(err)
	suite.NotNil(suite.user)

	// Setting up the JWT service
	suite.jwtService = jwt.New(jwt.Config{
		SecretKey: suite.secretKey,
		Issuer:    suite.issuer,
		ExpTime:   suite.expirationTime,
	})
}

func (suite *JWTServiceSuite) TestGenerateToken() {
	// Generate the token
	token, err := suite.jwtService.Generate(suite.user)
	suite.NoError(err)
	suite.NotEmpty(token)

	// Decode the token to validate the claims
	claims := jwt_builtin.MapClaims{}
	_, err = jwt_builtin.ParseWithClaims(token, &claims, func(token *jwt_builtin.Token) (interface{}, error) {
		return []byte(suite.secretKey), nil
	})
	suite.NoError(err)
	suite.Equal(suite.user.ID().String(), claims["user_id"])
	suite.Equal(false, claims["is_admin"])
	suite.Equal(suite.issuer, claims["iss"])
}

func (suite *JWTServiceSuite) TestDecodeToken() {
	// Create a valid token
	token, err := suite.jwtService.Generate(suite.user)
	suite.NoError(err)

	// Decode the token
	claims, err := suite.jwtService.Decode(token)
	suite.NoError(err)
	suite.Equal(suite.user.ID().String(), claims["user_id"])
	suite.Equal(false, claims["is_admin"])
	suite.Equal(suite.issuer, claims["iss"])
}

func TestJWTServiceSuite(t *testing.T) {
	suite.Run(t, new(JWTServiceSuite))
}
