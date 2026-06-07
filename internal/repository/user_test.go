package repository

import (
	"testing"

	"github.com/PhanNam1501/bookmark-management/internal/model"
	"github.com/PhanNam1501/bookmark-management/internal/test/fixture"
	"github.com/go-openapi/testify/v2/assert"
	"gorm.io/gorm"
)

func TestUser_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB func(t *testing.T) *gorm.DB

		inputUser *model.User

		expectedErrStr string
		expectedOutput *model.User

		verifyFunc func(db *gorm.DB, user *model.User)
	}{
		{
			name: "normal case",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputUser: &model.User{
				Base: model.Base{
					ID: "550e8400-e29b-41d4-a716-446655440000",
				},
				UserName:    "johndoe2",
				Password:    "hashed_password_123",
				DisplayName: "John Doe 2",
				Email:       "johndoe2@gmail.com",
			},

			expectedOutput: &model.User{
				Base: model.Base{
					ID: "550e8400-e29b-41d4-a716-446655440000",
				},
				UserName:    "johndoe2",
				Password:    "hashed_password_123",
				DisplayName: "John Doe 2",
				Email:       "johndoe2@gmail.com",
			},

			verifyFunc: func(db *gorm.DB, user *model.User) {
				checkUser := &model.User{}
				err := db.Where("username = ?", user.UserName).First(checkUser).Error
				assert.Nil(t, err)
				assert.Equal(t, checkUser.UserName, user.UserName)
			},
		},
		{
			name: "err case",

			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			},

			inputUser: &model.User{
				Base: model.Base{
					ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
				},
				UserName:    "johndoe",
				Password:    "hashed_password_123",
				DisplayName: "John Doe",
				Email:       "johndoe@gmail.com",
			},

			expectedErrStr: "UNIQUE constraint failed: users.email",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			testRepo := NewUser(db)

			res, err := testRepo.CreateUser(ctx, tc.inputUser)

			if tc.expectedErrStr == "" {
				assert.Nil(t, err)
			} else {
				assert.ErrorContains(t, err, tc.expectedErrStr)
			}
			assert.Equal(t, res, tc.expectedOutput)

			if err == nil && tc.verifyFunc != nil {
				tc.verifyFunc(db, tc.inputUser)
			}
		})
	}
}
