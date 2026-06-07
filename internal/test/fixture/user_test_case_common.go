package fixture

import (
	"github.com/PhanNam1501/bookmark-management/internal/model"
	"gorm.io/gorm"
)

type UserCommonTestDB struct {
	db *gorm.DB
}

func (f *UserCommonTestDB) SetupDB(db *gorm.DB) {
	f.db = db
}

func (f *UserCommonTestDB) Migrate() error {
	return f.db.AutoMigrate(&model.User{})
}

func (f *UserCommonTestDB) DB() *gorm.DB {
	return f.db
}

func (f *UserCommonTestDB) GenerateData() error {
	db := f.db.Session(&gorm.Session{})

	users := []*model.User{
		{
			Base: model.Base{
				ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			},
			UserName:    "johndoe",
			Password:    "hashed_password_123",
			DisplayName: "John Doe",
			Email:       "johndoe@gmail.com",
		},
		{
			Base: model.Base{
				ID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
			},
			UserName:    "janedoe",
			Password:    "hashed_password_456",
			DisplayName: "Jane Doe",
			Email:       "janedoe@gmail.com",
		},
	}

	return db.CreateInBatches(users, 10).Error
}
