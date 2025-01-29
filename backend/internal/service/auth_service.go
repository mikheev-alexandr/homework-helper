package service

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mikheev-alexandr/pet-project/backend/internal/models"
	"github.com/mikheev-alexandr/pet-project/backend/internal/repository"
	"github.com/mikheev-alexandr/pet-project/backend/pkg/codegen"
	"golang.org/x/crypto/argon2"
)

type AuthService struct {
	repos *repository.Repository

	emailSender EmailSender
}

func NewAuthService(repos *repository.Repository, emailSender EmailSender) *AuthService {
	return &AuthService{
		repos:       repos,
		emailSender: emailSender,
	}
}

type authClaims struct {
	Id   int `json:"id"`
	Role int `json:"role"` // 0 - teacher, 1 - student
	jwt.RegisteredClaims
}

type confirmationClaims struct {
	Id int `json:"id"`
	jwt.RegisteredClaims
}

const (
	teacherRole = 0
	sruedntRole = 1
)

func (s *AuthService) CreateTeacher(teacher models.Teacher) (string, error) {
	teacher.Password = generatePasswordHash(teacher.Password)
	teacherId, err := s.repos.Authorization.CreateTeacher(teacher)
	if err != nil {
		return "", err
	}

	claims := confirmationClaims{
		teacherId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signingKey, err := getECDSAPrivateKeyFromEnv()
	if err != nil {
		return "", err
	}

	return token.SignedString(signingKey)
}

func (s *AuthService) CreateStudent(teacherId int, name string, classNum int) (models.Student, error) {
	var student models.Student

	secretKey := os.Getenv("SYMMETRICK_KEY")

	codeWord, password, err := s.repos.Authorization.GetCodeWord()
	if err != nil {
		return student, err
	}

	password, err = codegen.Decrypt(password, []byte(secretKey))
	if err != nil {
		return student, err
	}

	student = models.Student{
		Name:        name,
		Code:        codeWord,
		Password:    generatePasswordHash(password),
		ClassNumber: classNum,
	}

	student.Id, err = s.repos.Authorization.CreateStudent(teacherId, student)
	if err != nil {
		return models.Student{}, err
	}

	student.Password = password

	return student, nil
}

func (s *AuthService) GetTeacherByEmail(email string) (models.Teacher, error) {
	return s.repos.Authorization.GetTeacherByEmail(email)
}

func (s *AuthService) UpdateStudentPassword(studentId int, oldPassword, newPassword string) error {
	oldPassword = generatePasswordHash(oldPassword)
	newPassword = generatePasswordHash(newPassword)

	if err := s.repos.Authorization.UpdateStudentPassword(studentId, oldPassword, newPassword); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) UpdateTeacherPassword(teacherId int, newPassword string) error {
	newPassword = generatePasswordHash(newPassword)

	if err := s.repos.Authorization.UpdateTeacherPassword(teacherId, newPassword); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) SendConfirmationEmail(email, token string) error {
	link := fmt.Sprintf("http://localhost:3000/auth/confirm?token=%s", token)
	subject := "Confirm your email"
	body := fmt.Sprintf("Click the link to confirm your email: %s", link)
	return s.emailSender.SendEmail(email, subject, body)
}

func (s *AuthService) SendResetEmail(email, token string) error {
	link := fmt.Sprintf("http://localhost:3000/auth/teacher/update-password?token=%s", token)
	subject := "Reset your password"
	body := fmt.Sprintf("Click the link to reset your password: %s", link)
	return s.emailSender.SendEmail(email, subject, body)
}

func (s *AuthService) ConfirmEmail(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &confirmationClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getECDSAPublicKeyFromEnv()
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*confirmationClaims)
	if !ok {
		return 0, errors.New("wrong type of token claims")
	}

	return claims.Id, nil
}

func (s *AuthService) GenerateResetToken(id int) (string, error) {
	claims := confirmationClaims{
		id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signingKey, err := getECDSAPrivateKeyFromEnv()
	if err != nil {
		return "", err
	}

	return token.SignedString(signingKey)
}

func (s *AuthService) ActivateUser(userId int) error {
	return s.repos.Authorization.ActivateUser(userId)
}

func (s *AuthService) GenerateTeacherToken(email, password string) (string, error) {
	password = generatePasswordHash(password)
	teacher, err := s.repos.Authorization.GetTeacher(email, password)
	if err != nil {
		return "", err
	}

	claims := authClaims{
		teacher.Id,
		teacherRole,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signingKey, err := getECDSAPrivateKeyFromEnv()
	if err != nil {
		return "", err
	}

	return token.SignedString(signingKey)
}

func (s *AuthService) GenerateStudentToken(codeWord, password string) (string, error) {
	password = generatePasswordHash(password)
	student, err := s.repos.Authorization.GetStudent(codeWord, password)
	if err != nil {
		return "", err
	}

	claims := authClaims{
		student.Id,
		sruedntRole,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signingKey, err := getECDSAPrivateKeyFromEnv()
	if err != nil {
		return "", err
	}

	return token.SignedString(signingKey)
}

func (s *AuthService) ParseToken(accessToken string) (int, int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &authClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getECDSAPublicKeyFromEnv()
	})
	if err != nil {
		return 0, 0, err
	}

	claims, ok := token.Claims.(*authClaims)
	if !ok {
		return 0, 0, errors.New("wrong type of token claims")
	}

	if claims.Role == 1 {
		err = s.repos.Authorization.GetStudentById(claims.Id)
	} else {
		err = s.repos.Authorization.GetTeacherById(claims.Id)
	}

	if err != nil {
		return 0, 0, errors.New("user does nit exist")
	}

	return claims.Id, claims.Role, nil
}

func (s *AuthService) ParseResetToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &confirmationClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getECDSAPublicKeyFromEnv()
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*confirmationClaims)
	if !ok {
		return 0, errors.New("wrong type of token claims")
	}

	return claims.Id, nil
}

func generatePasswordHash(password string) string {
	salt := os.Getenv("SALT")

	key := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)

	return fmt.Sprintf("%x", key)
}

func getECDSAPrivateKeyFromEnv() (*ecdsa.PrivateKey, error) {
	signingKey := os.Getenv("PRIVATE_KEY")
	block, _ := pem.Decode([]byte(signingKey))
	if block == nil {
		return nil, errors.New("signing private key is not correct")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func getECDSAPublicKeyFromEnv() (*ecdsa.PublicKey, error) {
	signingKey := os.Getenv("PUBLIC_KEY")
	block, _ := pem.Decode([]byte(signingKey))
	if block == nil {
		return nil, errors.New("signing public key is not correct")
	}

	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := parsedKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("key is not of type *ecdsa.PublicKey")
	}

	return key, nil
}
