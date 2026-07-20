package user

import (
	"context"
	"errors"
	"math/rand"

	"encore.app/internal/api_errors"
	"encore.app/internal/password"
	"encore.app/services/accounts/db"
	"encore.dev/rlog"
	"github.com/jackc/pgx/v5/pgconn"
)

type GetOrCreateCustomerParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type GetOrCreateCustomerResponse struct {
	UserID int64 `json:"userId"`
}

func (s *UserService) GetOrCreateCustomer(ctx context.Context, p GetOrCreateCustomerParams) (*GetOrCreateCustomerResponse, error) {
	customer, err := s.query.GetCustomerByPhoneAndEmail(ctx, db.GetCustomerByPhoneAndEmailParams{
		PhoneNumber: &p.Phone,
		Email:       p.Email,
	})

	if err == nil {
		return &GetOrCreateCustomerResponse{
			UserID: customer.ID,
		}, nil
	}

	if !errors.Is(err, db.ErrNoRows) {
		rlog.Error("failed to get customer by phone and email", "error", err, "phone", p.Phone, "email", p.Email)
		return nil, api_errors.ErrInternalError
	}

	id, err := s.createCustomer(ctx, p)
	if err != nil {
		return nil, err
	}

	return &GetOrCreateCustomerResponse{
		UserID: id,
	}, nil
}

func getRandomPasswordHash() (string, error) {
	randomPassword := generateRandomPassword()
	return password.HashPassword(randomPassword)
}

// generateRandomPassword generates a random password with all the required password complexity: length 8, one uppercase letter, one lowercase letter, one digit, and one special character.
func generateRandomPassword() string {
	var passwordLength = 8
	var uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	var digits = "0123456789"
	var specialCharacters = "!@#$%^&*()-_=+[]{}|;:,.<>?/"

	var allCharacters = uppercaseLetters + lowercaseLetters + digits + specialCharacters

	password := make([]byte, passwordLength)

	// Ensure at least one character from each category is included
	password[0] = uppercaseLetters[rand.Intn(len(uppercaseLetters))]
	password[1] = lowercaseLetters[rand.Intn(len(lowercaseLetters))]
	password[2] = digits[rand.Intn(len(digits))]
	password[3] = specialCharacters[rand.Intn(len(specialCharacters))]

	// Fill the remaining characters randomly from all categories
	for i := 4; i < passwordLength; i++ {
		password[i] = allCharacters[rand.Intn(len(allCharacters))]
	}

	// Shuffle the password to randomize the order of characters
	rand.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password)
}

func (s *UserService) createCustomer(ctx context.Context, p GetOrCreateCustomerParams) (int64, error) {
	passwordHash, err := getRandomPasswordHash()
	if err != nil {
		rlog.Error("failed to generate random password hash", "error", err)
		return 0, api_errors.ErrInternalError
	}

	newCust, err := s.query.CreateCustomer(ctx, db.CreateCustomerParams{
		FirstName:    p.FirstName,
		LastName:     p.LastName,
		Email:        p.Email,
		PhoneNumber:  &p.Phone,
		PasswordHash: passwordHash,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_key":
				return 0, ErrEmailAlreadyExists
			case "users_phone_number_key":
				return 0, ErrPhoneAlreadyExists
			default:
				rlog.Error("failed to create customer due to unique constraint violation", "error", err, "constraint", pgErr.ConstraintName)
				return 0, api_errors.ErrInternalError
			}
		}
		rlog.Error("failed to create customer", "error", err)
		return 0, api_errors.ErrInternalError
	}

	return newCust.ID, nil
}
