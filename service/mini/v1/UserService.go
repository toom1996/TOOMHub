// @Description
// @Author    2020/8/19 14:38
package ServiceMiniV1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
	ModelMiniV1 "toomhub/model/mini/v1"
	"toomhub/util"
	"toomhub/validatorRules"
)

// @title	通过OPENID获取用户信息
// @description
// @auth	toom <1023150697@qq.com>
func GetUser(openid string, validator *validatorRules.Login) (interface{}, error) {
	//db := util.DB
	//user, err := HasUser(openid)
	//
	//if err != nil {
	//	//如果是没有这个用户
	//	if gorm.IsRecordNotFoundError(err) {
	//		res, err := UserCreate(openid, db, validator)
	//		return res, err
	//	}
	//
	//	// 查询失败
	//	return "", err
	//}
	//
	////重新登陆就会刷新token 和 refreshToken
	//_, _ = UpdateUserInfoToRedis(user.MiniId)
	//
	//info, err := GetUserInfoByRedis(user.MiniId)
	//
	//if err != nil {
	//	userModel := ModelMiniV1.ToomhubUserMini{}
	//	query := db.Where("mini_id = ?", user.MiniId).Find(&userModel)
	//	if query.Error != nil {
	//		fmt.Println(query.Error)
	//	}
	//
	//	profileModel := ModelMiniV1.ToomhubUserMiniProfile{
	//		MiniId:    userModel.MiniId,
	//		NickName:  validator.RawData.NickName,
	//		Gender:    validator.RawData.Gender,
	//		City:      validator.RawData.City,
	//		Province:  validator.RawData.Province,
	//		Country:   validator.RawData.Country,
	//		AvatarUrl: validator.RawData.AvatarUrl,
	//	}
	//
	//	miniTokenModel := ModelMiniV1.ToomhubUserMiniToken{}
	//	query = db.Where("mini_id = ?", user.MiniId).Find(&miniTokenModel)
	//	if query.Error != nil {
	//		fmt.Println(query.Error)
	//	}
	//	_, _ = SetUserInfoToRedis(userModel, profileModel, miniTokenModel)
	//}
	//info, err = GetUserInfoByRedis(user.MiniId)
	//return info, nil
	return gin.H{}, nil
}

// @title	判断是否为新用户
func HasUser(openid string) (ModelMiniV1.ToomhubUserMini, error) {
	db := util.DB
	tableModel := ModelMiniV1.ToomhubUserMini{}
	//根据openid 查找用户, 如果有这个用户, 就从缓存中获取信息
	query := db.Where("open_id = ?", openid).Take(&tableModel)

	//查询失败
	if query.Error != nil {
		return tableModel, query.Error
	}

	return tableModel, nil
}

type UserInfo struct {
	Id           string
	AccessToken  string
	RefreshToken string
}

// @title	创建一个新的小程序用户
// @desc
// @auth toom <1023150697@qq.com>
func UserCreate(openid string, DB *gorm.DB, validator *validatorRules.Login) (interface{}, error) {
	fmt.Println("not found this user")
	createTime := time.Now().Unix()
	//开启事务
	transaction := DB.Begin()

	//插入到用户表
	userModel := ModelMiniV1.ToomhubUserMini{
		OpenId:    openid,
		CreatedAt: createTime,
	}
	userQuery := transaction.Create(&userModel).Scan(&userModel)
	if userQuery.Error != nil {
		transaction.Rollback()
		return "", userQuery.Error
	}

	profileModel := ModelMiniV1.ToomhubUserMiniProfile{
		MiniId:    userModel.MiniId,
		NickName:  validator.RawData.NickName,
		Gender:    validator.RawData.Gender,
		City:      validator.RawData.City,
		Province:  validator.RawData.Province,
		Country:   validator.RawData.Country,
		AvatarUrl: validator.RawData.AvatarUrl,
	}

	profileQuery := transaction.Create(&profileModel)
	if profileQuery.Error != nil {
		transaction.Rollback()
		return "", profileQuery.Error
	}

	//生成token
	token, err := util.GenerateToken(userModel.MiniId)
	if err != nil {
		fmt.Println(err)
	}

	refreshToken := util.GetRandomString(64)
	tokenModel := ModelMiniV1.ToomhubUserMiniToken{
		MiniId:       userModel.MiniId,
		AccessToken:  token,
		RefreshToken: refreshToken,
		CreatedAt:    createTime,
		UpdatedAt:    createTime,
	}
	err = transaction.Create(&tokenModel).Error
	if err != nil {
		fmt.Println(err)
	}

	_, _ = SetUserInfoToRedis(userModel, profileModel, tokenModel)

	//提交事务
	transaction.Commit()
	if transaction.Error != nil {
		fmt.Println(transaction.Error)
	}

	info, _ := GetUserInfoByRedis(userModel.MiniId)

	return info, err
}

type RedisUserInfo struct {
	MiniId      int
	OpenId      string
	CreatedAt   int64
	NickName    string
	Gender      int8
	City        string
	Province    string
	Country     string
	AvatarUrl   string
	FansCount   string
	LikesCount  string
	FollowCount string
}

// @title	从REDIS中获取用户信息
func GetUserInfoByRedis(userId int64) (interface{}, error) {
	//从redis中获取
	id := strconv.Itoa(int(userId))
	query := util.Rdb.HMGet(util.Ctx, util.UserCacheKey+id, []string{
		"mini_id",
		"open_id",
		"avatar_url",
		"created_at",
		"nick_name",
		"gender",
		"city",
		"province",
		"country",
		"token",
		"refresh_token",
		"fans_count",
		"likes_count",
		"follow_count",
		"exp",
	}...)
	if query.Err() != nil {
		fmt.Println(query.Err())
		return "", query.Err()
	}
	res, err := query.Result()
	if err != nil {
		fmt.Println(err)
		return "", query.Err()
	}

	if res[0] != nil {
		for index, _ := range res {
			if res[index] != nil {
				p, err := strconv.Atoi(res[index].(string))
				if err == nil {
					res[index] = p
				}
			} else {
				res[index] = 0
			}
		}
		m := map[string]interface{}{
			"MiniId":        res[0],
			"OpenId":        res[1],
			"avatar_url":    res[2],
			"CreatedAt":     res[3],
			"nick_name":     res[4],
			"Gender":        res[5],
			"City":          res[6],
			"Province":      res[7],
			"Country":       res[8],
			"token":         res[9],
			"refresh_token": res[10],
			"fans_count":    res[11],
			"likes_count":   res[12],
			"follow_count":  res[13],
			"exp":           res[14],
		}

		return m, nil
	}

	//没有缓存就重新插入一份
	return "", errors.New("unknown error")
}

// @title	将用户信息塞入REDIS缓存
func SetUserInfoToRedis(userModel ModelMiniV1.ToomhubUserMini, profileModel ModelMiniV1.ToomhubUserMiniProfile, tokenModel ModelMiniV1.ToomhubUserMiniToken) (interface{}, error) {
	key := util.UserCacheKey + strconv.Itoa(int(userModel.MiniId))

	info := map[string]interface{}{
		"mini_id":       userModel.MiniId,
		"open_id":       userModel.OpenId,
		"created_at":    userModel.CreatedAt,
		"nick_name":     profileModel.NickName,
		"gender":        profileModel.Gender,
		"city":          profileModel.City,
		"province":      profileModel.Province,
		"country":       profileModel.Country,
		"avatar_url":    profileModel.AvatarUrl,
		"token":         tokenModel.AccessToken,
		"refresh_token": tokenModel.RefreshToken,
		"fans_count":    0,
		"follow_count":  0,
		"likes_count":   0,
	}

	//塞入redis
	err := util.Rdb.HMSet(util.Ctx, key, info).Err()

	if err != nil {
		fmt.Println(err)
	}
	return info, nil
}

// @title	刷新用户的token 和 refreshToken
func UpdateUserInfoToRedis(miniId int64) (interface{}, error) {
	//db := util.DB
	//
	//profileModel := ModelMiniV1.ToomhubUserMiniProfile{}
	//
	//_ = db.Where("mini_id = ?", miniId).Take(&profileModel)
	//key := util.UserCacheKey + strconv.Itoa(int(miniId))
	////塞入redis
	//token, _ := util.GenerateToken(miniId)
	//fmt.Println("token -> ", token)
	//refreshToken := util.GetRandomString(64)
	//fmt.Println(refreshToken)
	//err := util.Rdb.HMSet(util.Ctx, key, map[string]interface{}{
	//	"nick_name":     profileModel.NickName,
	//	"gender":        profileModel.Gender,
	//	"city":          profileModel.City,
	//	"province":      profileModel.Province,
	//	"country":       profileModel.Country,
	//	"avatar_url":    profileModel.AvatarUrl,
	//	"token":         token,
	//	"refresh_token": refreshToken,
	//}).Err()
	//
	//// 更新数据库 (感觉没啥必要, 因为数据都是从redis取的, 并且redis的数据都已经更新了)
	//updateTime := time.Now().Unix()
	//tokenModel := ModelMiniV1.ToomhubUserMiniToken{
	//	UpdatedAt:    updateTime,
	//	AccessToken:  token,
	//	RefreshToken: refreshToken,
	//}
	//_ = db.Table("toomhub_user_mini_token").Update(&tokenModel).Where("mini_id = ?", miniId)
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//return map[string]string{
	//	"token":         token,
	//	"refresh_token": refreshToken,
	//}, nil
	return map[string]string{
		"token":         "",
		"refresh_token": "",
	}, nil
}

func Refresh(validator *validatorRules.Refresh) (bool, error) {

	db := util.DB

	db.Table("toomhub_user_mini_token").Where("access_token = ?", validator.Token)
	return true, nil
}
