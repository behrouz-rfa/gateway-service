package db

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/behrouz-rfa/gateway-service/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	dialector := postgres.New(postgres.Config{
		Conn: db,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	return gormDB, mock, nil
}

func TestCreateUser(t *testing.T) {
	gormDB, mock, err := setupMockDB()
	assert.NoError(t, err)
	defer func() {
		sqlDB, _ := gormDB.DB()
		sqlDB.Close()
	}()

	repo := &DbRepository{db: gormDB}
	tu := time.Now()
	ct := time.Now()
	userInput := &domain.User{
		Base: domain.Base{
			ID:        "0fcd28d8-7f43-4543-a92f-d7eb585f8d06",
			CreatedAt: ct,
			UpdatedAt: &tu,
		},
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "users" \("id","created_at","updated_at","name","email","password"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6\)`).
		WithArgs(
			sqlmock.AnyArg(), // Match any UUID
			userInput.CreatedAt,
			userInput.UpdatedAt,
			userInput.Name,
			userInput.Email,
			userInput.Password,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	userID, err := repo.CreateUser(context.Background(), userInput)
	assert.NoError(t, err)
	assert.NotEmpty(t, userID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUser(t *testing.T) {
	gormDB, mock, err := setupMockDB()
	assert.NoError(t, err)
	defer mock.ExpectClose()

	repo := &DbRepository{db: gormDB}

	userID := "1"
	expectedUser := &domain.User{
		Base:     domain.Base{ID: userID},
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	users := sqlmock.NewRows([]string{"id", "name", "email", "password"}).
		AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Password)
	expectedSQL := "SELECT (.+) FROM \"users\" WHERE id =(.+)"
	mock.ExpectQuery(expectedSQL).WillReturnRows(users)

	user, err := repo.GetUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserPlan(t *testing.T) {
	gormDB, mock, err := setupMockDB()
	assert.NoError(t, err)
	defer mock.ExpectClose()

	repo := &DbRepository{db: gormDB}

	userID := "1"
	expectedPlan := &domain.Plan{
		UserID:      userID,
		CreditsUsed: 10,
	}
	expectedPlan.ID = "1"
	rows := sqlmock.NewRows([]string{"id", "user_id", "credits_used"}).
		AddRow(expectedPlan.ID, expectedPlan.UserID, expectedPlan.CreditsUsed)

	mock.ExpectQuery("SELECT (.+) FROM \"palns\" WHERE user_id =(.+)").
		WillReturnRows(rows)

	plan, err := repo.GetUserPlan(userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedPlan, plan)
	assert.NoError(t, mock.ExpectationsWereMet())
}
