package meta

import "BaiduNetdisk_golang/db"

type UserMeta struct {
	Id       int64
	UserName string
	Phone    string
}

func GetUserMetaDB(username, password string) (UserMeta, error) {
	tUser, err := db.OnLoginHandle(username, password)
	if err != nil {
		return UserMeta{}, err
	}
	fMeta := UserMeta{
		Id:       tUser.Id.Int64,
		UserName: tUser.UserName.String,
		Phone:    tUser.Phone.String,
	}
	return fMeta, nil
}
