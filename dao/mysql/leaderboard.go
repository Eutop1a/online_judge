package mysql

type UserRank struct {
	FinishProblemNum int64  `json:"finish_num"`
	UserName         string `json:"username"`
}

func GetUserLeaderboard() (data *[]UserRank, err error) {
	var retSlice []UserRank
	var userSlice []User
	err = DB.Model(&User{}).Select("finish_num, username").Order("finish_num desc").Find(&userSlice).Error
	if err != nil {
		return nil, err
	}
	for _, user := range userSlice {
		retSlice = append(retSlice, UserRank{
			FinishProblemNum: user.FinishProblemNum,
			UserName:         user.UserName,
		})
	}
	return &retSlice, nil
}
