package model

type User struct {
	Base
	UserName    string `gorm:"unique;column:username" json:"username"`
	Password    string `gorm:"column:password" json:"-"`
	DisplayName string `gorm:"column:display_name" json:"display_name"`
	Email       string `gorm:"column:email;unique" json:"email"`
}

// func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
// 	if u.ID == "" {
// 		u.ID = uuid.New().String()
// 	}
// 	return nil
// }
