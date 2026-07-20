package accounts

import (
	"context"
	"strconv"
	"testing"
	"time"

	"encore.app/internal/api_errors"
	"encore.app/internal/jwt"
	"encore.app/internal/lang"
	"encore.app/internal/validation"
	"encore.app/services/accounts/db"
	ah "encore.app/services/accounts/handlers/auth"
	user "encore.app/services/accounts/handlers/user"
	"encore.dev/beta/auth"
	"encore.dev/beta/errs"
	"encore.dev/et"
)

func adminAuthContext(adminID int64) context.Context {
	uid := auth.UID(strconv.FormatInt(adminID, 10))
	return auth.WithContext(context.Background(), uid, &AuthData{
		UserID: adminID,
		Role:   UserRoleAdmin,
	})
}

func agentAuthContext(agentID int64, adminRefID *int64) context.Context {
	uid := auth.UID(strconv.FormatInt(agentID, 10))
	return auth.WithContext(context.Background(), uid, &AuthData{
		UserID:     agentID,
		Role:       UserRoleAgent,
		AdminRefID: adminRefID,
	})
}

const (
	testPassword = "ValidPass123!"
	testEmail    = "valid_email@example.com"
)

func TestLogin(t *testing.T) {
	ctx := context.Background()

	t.Run("Invalid email", func(t *testing.T) {
		cases := []ah.LoginParams{
			{Email: "", Password: testPassword},
			{Email: "ab", Password: testPassword},
			{Email: "xsxs@@dd.com", Password: testPassword},
		}

		for _, p := range cases {
			err := p.Validate()
			expectedErr := api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
				Code:  api_errors.CodeInvalidValue,
				Field: "email",
			})

			api_errors.AssertApiError(t, expectedErr, err)
		}
	})

	t.Run("Invalid password", func(t *testing.T) {
		p := ah.LoginParams{
			Email: testEmail,
		}
		err := p.Validate()
		expectedErr := api_errors.NewErrorWithDetail(errs.InvalidArgument, validation.InvalidValueMsg, api_errors.ErrorDetails{
			Code:  api_errors.CodeInvalidValue,
			Field: "password",
		})

		api_errors.AssertApiError(t, expectedErr, err)
	})

	t.Run("User not found", func(t *testing.T) {
		_, err := Login(ctx, ah.LoginParams{Email: testEmail, Password: testPassword})
		api_errors.AssertApiError(t, ah.ErrInvalidCredentials, err)
	})

	t.Run("Incorrect password", func(t *testing.T) {
		user, err := CreateAdmin(ctx, user.CreateAdminParams{
			FirstName: "Test",
			LastName:  "Admin",
			Email:     testEmail,
			Password:  testPassword,
		})
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}
		defer query.DeleteUser(ctx, user.ID)

		_, err = Login(ctx, ah.LoginParams{Email: testEmail, Password: "WrongPass123!"})
		api_errors.AssertApiError(t, ah.ErrInvalidCredentials, err)
	})

	t.Run("Successful login", func(t *testing.T) {
		adminEmail := "admin_" + testEmail
		_, delAdmin, err := registerAdmin(ctx, adminEmail, testPassword)

		if err != nil {
			t.Fatalf("Failed to create test admin: %v", err)
		}

		defer delAdmin()

		agentEmail := "agent_" + testEmail
		_, delAgent, err := createAgent(ctx, user.CreateAgentParams{
			FirstName:   "Test",
			LastName:    "Agent",
			Email:       agentEmail,
			Password:    testPassword,
			PhoneNumber: "0505050505",
		})

		if err != nil {
			t.Fatalf("Failed to create test agent: %v", err)
		}

		defer delAgent()

		cases := []struct {
			name  string
			email string
		}{
			{name: "Admin user", email: adminEmail},
			{name: "Agent user", email: agentEmail},
		}

		for _, c := range cases {

			resp, err := Login(ctx, ah.LoginParams{Email: c.email, Password: testPassword})
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if resp.AccessToken == "" {
				t.Fatal("Expected access token, got empty string")
			}
			refreshToken := refreshCookieToken(t, resp)
			if refreshToken == "" {
				t.Fatal("Expected refresh token, got empty string")
			}

			accessClaims, err := jwt.ValidateAccessToken(resp.AccessToken)
			if err != nil {
				t.Fatalf("Failed to validate access token: %v", err)
			}

			user, err := query.GetUserByEmail(ctx, c.email)
			if err != nil {
				t.Fatalf("Failed to query user: %v, user: %s", err, c.email)
			}

			var orgCtx *jwt.OrganizationContext
			if user.OfficeID != nil && user.OrganizationID != nil && user.IsOrganic != nil {
				orgCtx = &jwt.OrganizationContext{
					OfficeID:       *user.OfficeID,
					OrganizationID: *user.OrganizationID,
					IsOrganic:      *user.IsOrganic,
				}
			}

			assertAccessClaims(t, accessClaims, jwt.AccessTokenData{
				UserID:              user.ID,
				Role:                user.Role,
				OrganizationContext: orgCtx,
			})
			if time.Until(accessClaims.ExpiresAt.Time) <= 0 {
				t.Error("Access token already expired")
			}

			refreshClaims, err := jwt.ValidateRefreshToken(refreshToken)
			if err != nil {
				t.Fatalf("Failed to validate refresh token: %v", err)
			}

			assertRefreshClaims(t, refreshClaims, user.ID)
			if time.Until(refreshClaims.ExpiresAt.Time) <= 0 {
				t.Error("Refresh token already expired")
			}

			// Verify stored refresh token in DB
			rt, err := query.GetRefreshToken(ctx, refreshClaims.ID)
			if err != nil {
				t.Fatalf("Failed to retrieve refresh token from DB: %v", err)
			}
			assertTimeAlmostEqual(t, rt.ExpiresAt.Time, refreshClaims.ExpiresAt.Time)
			if rt.UserID != user.ID {
				t.Errorf("Expected token.UserID %d, got %d", user.ID, rt.UserID)
			}
			if rt.Jti != refreshClaims.ID {
				t.Errorf("Expected token.JTI %s, got %s", refreshClaims.ID, rt.Jti)
			}
		}

	})
}

func TestLoginAsUser(t *testing.T) {
	ctx := context.Background()

	t.Run("Validation: missing user ID", func(t *testing.T) {
		err := (ah.LoginAsUserParams{}).Validate()
		expectedErr := invalidValueErr("userId")
		api_errors.AssertApiError(t, expectedErr, err)
	})

	t.Run("User not found", func(t *testing.T) {
		adminCtx := adminAuthContext(999)
		_, err := LoginAsUser(adminCtx, ah.LoginAsUserParams{UserID: 99999})
		api_errors.AssertApiError(t, ah.ErrInvalidCredentials, err)
	})

	t.Run("Disallow non-agent/customer role", func(t *testing.T) {
		adminEmail := generateTestEmail()
		admin, delAdmin, err := registerAdmin(ctx, adminEmail, testPassword)
		if err != nil {
			t.Fatalf("Failed to create admin: %v", err)
		}
		defer delAdmin()

		adminCtx := adminAuthContext(admin.ID)
		// Try to login as the admin themselves (role=admin) — must be rejected
		_, err = LoginAsUser(adminCtx, ah.LoginAsUserParams{UserID: admin.ID})
		api_errors.AssertApiError(t, api_errors.ErrUnauthorized, err)
	})

	t.Run("Successful login as agent", func(t *testing.T) {
		adminEmail := generateTestEmail()
		admin, delAdmin, err := registerAdmin(ctx, adminEmail, testPassword)
		if err != nil {
			t.Fatalf("Failed to create admin: %v", err)
		}
		defer delAdmin()

		agentEmail := generateTestEmail()
		agent, delAgent, err := createAgent(ctx, user.CreateAgentParams{
			FirstName:   "Test",
			LastName:    "Agent",
			Email:       agentEmail,
			Password:    testPassword,
			PhoneNumber: randomIsraeliPhoneNumber(),
		})
		if err != nil {
			t.Fatalf("Failed to create agent: %v", err)
		}
		defer delAgent()

		adminCtx := adminAuthContext(admin.ID)
		resp, err := LoginAsUser(adminCtx, ah.LoginAsUserParams{UserID: agent.ID})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if resp.AccessToken == "" {
			t.Fatal("Expected access token")
		}
		refreshToken := refreshCookieToken(t, resp)
		if refreshToken == "" {
			t.Fatal("Expected refresh token")
		}
		if resp.ID != agent.ID {
			t.Errorf("Expected ID %d, got %d", agent.ID, resp.ID)
		}

		// Verify access token includes adminRefID
		accessClaims, err := jwt.ValidateAccessToken(resp.AccessToken)
		if err != nil {
			t.Fatalf("Failed to validate access token: %v", err)
		}
		if accessClaims.AdminRefID == nil {
			t.Fatal("Expected AdminRefID in access token claims")
		}
		if *accessClaims.AdminRefID != admin.ID {
			t.Errorf("Expected AdminRefID %d, got %d", admin.ID, *accessClaims.AdminRefID)
		}

		agentUser, err := query.GetUserByEmail(ctx, agentEmail)
		if err != nil {
			t.Fatalf("Failed to query agent: %v", err)
		}

		var orgCtx *jwt.OrganizationContext
		if agentUser.OfficeID != nil && agentUser.OrganizationID != nil && agentUser.IsOrganic != nil {
			orgCtx = &jwt.OrganizationContext{
				OfficeID:       *agentUser.OfficeID,
				OrganizationID: *agentUser.OrganizationID,
				IsOrganic:      *agentUser.IsOrganic,
			}
		}

		assertAccessClaims(t, accessClaims, jwt.AccessTokenData{
			UserID:              agentUser.ID,
			Role:                agentUser.Role,
			OrganizationContext: orgCtx,
		})

		// Verify refresh token stored with admin ref
		refreshClaims, err := jwt.ValidateRefreshToken(refreshToken)
		if err != nil {
			t.Fatalf("Failed to validate refresh token: %v", err)
		}
		rt, err := query.GetRefreshToken(ctx, refreshClaims.ID)
		if err != nil {
			t.Fatalf("Failed to get stored refresh token: %v", err)
		}
		if rt.AdminRefID == nil {
			t.Fatal("Expected AdminRefID in stored refresh token")
		}
		if *rt.AdminRefID != admin.ID {
			t.Errorf("Expected stored AdminRefID %d, got %d", admin.ID, *rt.AdminRefID)
		}
	})

	t.Run("Successful login as customer", func(t *testing.T) {
		adminEmail := generateTestEmail()
		admin, delAdmin, err := registerAdmin(ctx, adminEmail, testPassword)
		if err != nil {
			t.Fatalf("Failed to create admin: %v", err)
		}
		defer delAdmin()

		phone := randomIsraeliPhoneNumber()
		customer, delCustomer, err := createCustomer(ctx, phone, nil)
		if err != nil {
			t.Fatalf("Failed to create customer: %v", err)
		}
		defer delCustomer()

		adminCtx := adminAuthContext(admin.ID)
		resp, err := LoginAsUser(adminCtx, ah.LoginAsUserParams{UserID: customer.ID})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if resp.AccessToken == "" {
			t.Fatal("Expected access token")
		}
		if resp.ID != customer.ID {
			t.Errorf("Expected ID %d, got %d", customer.ID, resp.ID)
		}

		// Verify OrganizationContext is nil for customers
		accessClaims, err := jwt.ValidateAccessToken(resp.AccessToken)
		if err != nil {
			t.Fatalf("Failed to validate access token: %v", err)
		}
		if accessClaims.OrganizationContext != nil {
			t.Error("Expected OrganizationContext to be nil for customer")
		}
		if accessClaims.AdminRefID == nil {
			t.Fatal("Expected AdminRefID in access token claims")
		}
		if *accessClaims.AdminRefID != admin.ID {
			t.Errorf("Expected AdminRefID %d, got %d", admin.ID, *accessClaims.AdminRefID)
		}
	})
}

func TestLoginBackToAdmin(t *testing.T) {
	ctx := context.Background()

	t.Run("No admin ref in auth data", func(t *testing.T) {
		agentCtx := agentAuthContext(1, nil)
		_, err := LoginBackToAdmin(agentCtx)
		api_errors.AssertApiError(t, ah.ErrInvalidCredentials, err)
	})

	t.Run("Admin not found", func(t *testing.T) {
		adminRefID := int64(99999)
		agentCtx := agentAuthContext(1, &adminRefID)
		_, err := LoginBackToAdmin(agentCtx)
		api_errors.AssertApiError(t, ah.ErrInvalidCredentials, err)
	})

	t.Run("Successful login back to admin", func(t *testing.T) {
		adminEmail := generateTestEmail()
		admin, delAdmin, err := registerAdmin(ctx, adminEmail, testPassword)
		if err != nil {
			t.Fatalf("Failed to create admin: %v", err)
		}
		defer delAdmin()

		agentEmail := generateTestEmail()
		agent, delAgent, err := createAgent(ctx, user.CreateAgentParams{
			FirstName:   "Test",
			LastName:    "Agent",
			Email:       agentEmail,
			Password:    testPassword,
			PhoneNumber: randomIsraeliPhoneNumber(),
		})
		if err != nil {
			t.Fatalf("Failed to create agent: %v", err)
		}
		defer delAgent()

		agentCtx := agentAuthContext(agent.ID, &admin.ID)
		resp, err := LoginBackToAdmin(agentCtx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if resp.AccessToken == "" {
			t.Fatal("Expected access token")
		}
		if refreshCookieToken(t, resp) == "" {
			t.Fatal("Expected refresh token")
		}
		if resp.ID != admin.ID {
			t.Errorf("Expected ID %d, got %d", admin.ID, resp.ID)
		}
		if resp.Role != db.UserRoleAdmin {
			t.Errorf("Expected role admin, got %s", resp.Role)
		}

		// Verify access token does NOT include adminRefID (back to admin session)
		accessClaims, err := jwt.ValidateAccessToken(resp.AccessToken)
		if err != nil {
			t.Fatalf("Failed to validate access token: %v", err)
		}
		if accessClaims.AdminRefID != nil {
			t.Error("Expected no AdminRefID in access token after returning to admin")
		}

		adminUser, err := query.GetUserByEmail(ctx, adminEmail)
		if err != nil {
			t.Fatalf("Failed to query admin: %v", err)
		}
		assertAccessClaims(t, accessClaims, jwt.AccessTokenData{
			UserID:              adminUser.ID,
			Role:                adminUser.Role,
			OrganizationContext: nil,
		})
	})
}

func TestSendCustomerLoginOTP(t *testing.T) {
	ctx := context.Background()

	t.Run("Validation: missing phone number", func(t *testing.T) {
		err := (ah.SendCustomerLoginOTPParams{}).Validate()
		api_errors.AssertApiError(t, invalidValueErr("phoneNumber"), err)
	})

	t.Run("User not found", func(t *testing.T) {
		err := SendCustomerLoginOTP(ctx, ah.SendCustomerLoginOTPParams{PhoneNumber: randomIsraeliPhoneNumber()})
		api_errors.AssertApiError(t, ah.ErrInvalidCredentials, err)
	})

	t.Run("User found but role is not customer", func(t *testing.T) {
		agentPhone := randomIsraeliPhoneNumber()
		agentEmail := generateTestEmail()
		_, delAgent, err := createAgent(ctx, user.CreateAgentParams{
			FirstName:   "Test",
			LastName:    "Agent",
			Email:       agentEmail,
			Password:    testPassword,
			PhoneNumber: agentPhone,
		})
		if err != nil {
			t.Fatalf("Failed to create agent: %v", err)
		}
		defer delAgent()

		err = SendCustomerLoginOTP(ctx, ah.SendCustomerLoginOTPParams{PhoneNumber: agentPhone})
		api_errors.AssertApiError(t, ah.ErrInvalidCredentials, err)
	})

	t.Run("Success: OTP saved and event published", func(t *testing.T) {
		phoneNumber := randomIsraeliPhoneNumber()
		customer, cleanup, err := createCustomer(ctx, phoneNumber, nil)
		if err != nil {
			t.Fatalf("Failed to create customer: %v", err)
		}
		defer cleanup()

		publishedBefore := len(et.Topic(ah.CustomerLoginOTPRequestedTopic).PublishedMessages())

		err = SendCustomerLoginOTP(ctx, ah.SendCustomerLoginOTPParams{PhoneNumber: phoneNumber})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		updatedUser, err := query.GetUserByPhone(ctx, &phoneNumber)
		if err != nil {
			t.Fatalf("Failed to fetch customer after send OTP: %v", err)
		}
		if updatedUser.ID != customer.ID {
			t.Fatalf("Expected customer ID %d, got %d", customer.ID, updatedUser.ID)
		}
		if updatedUser.Otp == nil {
			t.Fatal("Expected OTP to be saved")
		}
		if len(*updatedUser.Otp) != 6 {
			t.Fatalf("Expected OTP with length 6, got %q", *updatedUser.Otp)
		}

		publishedAfter := len(et.Topic(ah.CustomerLoginOTPRequestedTopic).PublishedMessages())
		if publishedAfter != publishedBefore+1 {
			t.Fatalf("Expected one published OTP event, before=%d after=%d", publishedBefore, publishedAfter)
		}

		published := et.Topic(ah.CustomerLoginOTPRequestedTopic).PublishedMessages()
		last := published[len(published)-1]
		if last == nil {
			t.Fatal("Expected published OTP event payload")
		}
		if last.LangCode != "he" {
			t.Fatalf("Expected default lang code 'he', got %q", last.LangCode)
		}
	})

	t.Run("Success: uses lang from context (en)", func(t *testing.T) {
		phoneNumber := randomIsraeliPhoneNumber()
		_, cleanup, err := createCustomer(ctx, phoneNumber, nil)
		if err != nil {
			t.Fatalf("Failed to create customer: %v", err)
		}
		defer cleanup()

		ctxWithLang := context.WithValue(ctx, lang.ContextKey, "en")

		publishedBefore := len(et.Topic(ah.CustomerLoginOTPRequestedTopic).PublishedMessages())

		err = SendCustomerLoginOTP(ctxWithLang, ah.SendCustomerLoginOTPParams{PhoneNumber: phoneNumber})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		published := et.Topic(ah.CustomerLoginOTPRequestedTopic).PublishedMessages()
		if len(published) != publishedBefore+1 {
			t.Fatalf("Expected one published OTP event, before=%d after=%d", publishedBefore, len(published))
		}

		last := published[len(published)-1]
		if last == nil {
			t.Fatal("Expected published OTP event payload")
		}
		if last.LangCode != "en" {
			t.Fatalf("Expected lang code 'en', got %q", last.LangCode)
		}
	})
}

func TestValidateCustomerLoginOTP(t *testing.T) {
	ctx := context.Background()

	t.Run("Validation: missing otp", func(t *testing.T) {
		err := (ah.ValidateCustomerLoginOTPParams{PhoneNumber: randomIsraeliPhoneNumber()}).Validate()
		api_errors.AssertApiError(t, invalidValueErr("otp"), err)
	})

	t.Run("User not found", func(t *testing.T) {
		_, err := ValidateCustomerLoginOTP(ctx, ah.ValidateCustomerLoginOTPParams{
			PhoneNumber: randomIsraeliPhoneNumber(),
			OTP:         "123456",
		})
		api_errors.AssertApiError(t, ah.ErrInvalidCredentials, err)
	})

	t.Run("Incorrect OTP", func(t *testing.T) {
		phoneNumber := randomIsraeliPhoneNumber()
		otp := "123456"
		_, cleanup, err := createCustomer(ctx, phoneNumber, &otp)
		if err != nil {
			t.Fatalf("Failed to create customer: %v", err)
		}
		defer cleanup()

		_, err = ValidateCustomerLoginOTP(ctx, ah.ValidateCustomerLoginOTPParams{
			PhoneNumber: phoneNumber,
			OTP:         "654321",
		})
		api_errors.AssertApiError(t, ah.ErrInvalidCredentials, err)
	})

	t.Run("Success: valid response and OTP cleared", func(t *testing.T) {
		phoneNumber := randomIsraeliPhoneNumber()
		otp := "123456"
		customer, cleanup, err := createCustomer(ctx, phoneNumber, &otp)
		if err != nil {
			t.Fatalf("Failed to create customer: %v", err)
		}
		defer cleanup()

		resp, err := ValidateCustomerLoginOTP(ctx, ah.ValidateCustomerLoginOTPParams{
			PhoneNumber: phoneNumber,
			OTP:         otp,
		})
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if resp.ID != customer.ID {
			t.Fatalf("Expected response ID %d, got %d", customer.ID, resp.ID)
		}
		if resp.Role != db.UserRoleCustomer {
			t.Fatalf("Expected role customer, got %s", resp.Role)
		}
		if resp.AccessToken == "" {
			t.Fatal("Expected access token")
		}
		if refreshCookieToken(t, resp) == "" {
			t.Fatal("Expected refresh token")
		}

		updatedUser, err := query.GetUserByPhone(ctx, &phoneNumber)
		if err != nil {
			t.Fatalf("Failed to fetch customer after validate OTP: %v", err)
		}
		if updatedUser.Otp != nil {
			t.Fatal("Expected OTP to be cleared after successful validation")
		}
	})
}

func TestGetCustomerToken(t *testing.T) {
	ctx := context.Background()
	t.Run("User not found", func(t *testing.T) {
		agent, delAgent, err := createAgent(ctx, user.CreateAgentParams{
			FirstName:   "Test",
			LastName:    "Agent",
			Email:       randomName() + "@example.com",
			Password:    testPassword,
			PhoneNumber: randomIsraeliPhoneNumber(),
		})
		if err != nil {
			t.Fatalf("Failed to create agent: %v", err)
		}
		defer delAgent()

		_, err = GetCustomerToken(ctx, ah.GetCustomerTokenParams{UserID: agent.ID})
		if err == nil {
			t.Fatalf("Expected error when getting customer token for agent, got nil")
		}

		api_errors.AssertApiError(t, api_errors.ErrNotFound, err)
	})

	t.Run("Success: valid response without refresh token", func(t *testing.T) {
		phoneNumber := randomIsraeliPhoneNumber()
		otp := "123456"
		customer, cleanup, err := createCustomer(ctx, phoneNumber, &otp)
		if err != nil {
			t.Fatalf("Failed to create customer: %v", err)
		}
		defer cleanup()

		resp, err := GetCustomerToken(ctx, ah.GetCustomerTokenParams{UserID: customer.ID})
		if err != nil {
			t.Fatalf("Failed to get customer token: %v", err)
		}

		if resp.ID != customer.ID {
			t.Fatalf("Expected response ID %d, got %d", customer.ID, resp.ID)
		}

		if resp.SetCookies != nil {
			t.Fatal("Expected no SetCookies in response for GetCustomerToken")
		}
	})
}

func createCustomer(ctx context.Context, phoneNumber string, otp *string) (db.User, func(), error) {
	_, err := query.CreateCustomer(ctx, db.CreateCustomerParams{
		Email:        generateTestEmail(),
		PhoneNumber:  &phoneNumber,
		Otp:          otp,
		PasswordHash: "test-password-hash",
	})
	if err != nil {
		return db.User{}, nil, err
	}

	user, err := query.GetUserByPhone(ctx, &phoneNumber)
	if err != nil {
		return db.User{}, nil, err
	}

	cleanup := func() {
		_ = query.DeleteUser(ctx, user.ID)
	}

	return user, cleanup, nil
}
