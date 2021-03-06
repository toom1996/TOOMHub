// @Description
// @Author: toom1996 <1023150697@qq.com>
// @dateTime: 2021/1/29 16:26
package service

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"time"
	"toomhub/model"
	rules "toomhub/rules/user/v1"
	"toomhub/setting"
	"toomhub/util"
)

type UserService struct {
}

type token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"` // 这个字段没用到
	Scope       string `json:"scope"`      // 这个字段也没用到
	State       string `json:"state"`      // 这个字段也没用到
}

func (service *UserService) GetGithubOAuthInfo(validator *rules.V1UserGithubOAuth) (model.ZawazawaUserProfileGithub, error) {

	//编译好链接
	s := fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		setting.ZConfig.GithubOAuth.ClientId, setting.ZConfig.GithubOAuth.ClientSecret, validator.Code,
	)
	var err error
	// 形成请求
	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, s, nil); err != nil {
		return model.ZawazawaUserProfileGithub{}, err
	}

	req.Header.Set("accept", "application/json")

	//发送请求并获得响应
	var httpClient = http.Client{}
	var res *http.Response
	if res, err = httpClient.Do(req); err != nil {
		return model.ZawazawaUserProfileGithub{}, err
	}

	// 将响应体解析为 token，并返回
	var token token
	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		return model.ZawazawaUserProfileGithub{}, err
	}
	fmt.Println(&token)

	var userInfo = model.ZawazawaUserProfileGithub{}
	// 形成请求
	var userInfoUrl = "https://api.github.com/user" // github用户信息获取接口

	if req, err = http.NewRequest(http.MethodGet, userInfoUrl, nil); err != nil {
		return userInfo, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token.AccessToken))

	// 发送请求并获取响应
	var client = http.Client{}
	if res, err = client.Do(req); err != nil {
		return userInfo, err
	}

	// 将响应的数据写入 userInfo 中，并返回
	if err = json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return userInfo, err
	}
	fmt.Println(userInfo.GitOauthID)
	return userInfo, nil
}

//是否为新用户
//查数据库
func (service *UserService) IsNewUser(gitId uint) (*gorm.DB, bool) {
	fmt.Println("gitid", gitId)
	result := util.DB.Select("id").Debug().Where(&model.ZawazawaUserProfileGithub{GitOauthID: gitId}).Find(&model.ZawazawaUserProfileGithub{})
	if result.RowsAffected != 0 {
		return result, false
	}
	return nil, true
}

//存储github信息
func (service *UserService) SaveGithubOAuthInfo(info *model.ZawazawaUserProfileGithub) (map[string]interface{}, error) {
	fmt.Println("SaveGithubOAuthInfo")
	db := util.DB
	transaction := util.DB.Begin()
	//存入数据库
	fmt.Println(info)
	err := model.ZawazawaUserProfileGithubMgr(db).Create(&info).Error

	if err != nil {
		fmt.Println(err)
		transaction.Rollback()
		return nil, err
	}

	user := model.ZawazawaUser{
		Nickname:  info.Name,
		OauthID:   info.GitOauthID,
		OauthType: util.OAuthGithub,
	}
	err = model.ZawazawaUserMgr(db).Create(&user).Error

	if err != nil {
		fmt.Println(err)
		transaction.Rollback()
		return nil, err

	}

	transaction.Commit()
	//TODO: 存入redis
	return map[string]interface{}{
		"avatar":   info.AvatarURL,
		"username": user.Nickname,
	}, nil
}

//更新github信息
func (service *UserService) UpdateGithubOAuthInfo(p *gorm.DB, info *model.ZawazawaUserProfileGithub) (map[string]interface{}, error) {
	fmt.Println("update")
	info.Name = "xx"
	p.Debug().Save(&info)
	//TODO: 存入redis
	return map[string]interface{}{
		"avatar":   info.AvatarURL,
		"username": "ddddddd",
	}, nil
}

// 是否已经注册
func (service *UserService) IsRegister(mobile string) interface{} {

	// TODO 判断REDIS是否存在

	// TODO 查数据库
	s := &model.ZawazawaUser{}
	err := model.ZawazawaUserMgr(util.DB).Select("id", "mobile").Debug().Where(&model.ZawazawaUser{
		Mobile: mobile,
	}).Find(s)
	if err.RowsAffected == 0 {
		return false
	}
	return s.ID
}

// 存储通过手机验证码登陆的新用户
func (service *UserService) SaveMobileUser(validator *rules.V1UserRegister) (interface{}, error) {
	t := util.DB.Begin()
	// 注册用户
	info := model.ZawazawaUser{
		Nickname:   "咋哇咋哇用户",
		Mobile:     validator.Mobile,
		ZawazawaID: "zawazawa_" + validator.Mobile,
	}
	err := util.DB.Create(&info).Error
	if err != nil {
		t.Rollback()
		return false, err
	}

	// 注册token
	g, err := util.GenerateToken(info.ID)

	if err != nil {
		t.Rollback()
		return false, err
	}

	token := model.ZawazawaUserToken{
		UId:          info.ID,
		Token:        g,
		RefreshToken: "zawazawa_" + validator.Mobile,
	}

	err = util.DB.Create(&token).Error
	if err != nil {
		t.Rollback()
		return false, err
	}

	t.Commit()

	return map[string]interface{}{
		"username":      "咋哇咋哇用户",
		"avatar":        "http://v.bootstrapmb.com/2019/6/mmjod5239/img/avatar7-sm.jpg",
		"expire":        setting.ZConfig.Jwt.JwtExpire,
		"issuing_time":  time.Now().Unix(),
		"token":         token.Token,
		"refresh_token": token.RefreshToken,
	}, nil
}

// 存储通过手机验证码登陆的新用户
func (service *UserService) GetMobileUser(id uint) (interface{}, error) {
	t := util.DB.Begin()

	// 生成token
	g, err := util.GenerateToken(id)
	if err != nil {
		t.Rollback()
		return false, err
	}

	token := model.ZawazawaUserToken{
		Token: g,
	}
	// 生成refresh token
	refreshToken, _ := util.GenerateRefreshToken(id)

	err = util.DB.Debug().Model(&model.ZawazawaUserToken{}).Where(&model.ZawazawaUserToken{
		UId: id,
	}).Updates(&token).Error
	if err != nil {
		t.Rollback()
		return false, err
	}

	t.Commit()

	setUserCache(uint(id), cacheUser{
		Id:          id,
		AccessToken: g,
	})

	return map[string]interface{}{
		"username":      "咋哇咋哇用户",
		"avatar":        "http://v.bootstrapmb.com/2019/6/mmjod5239/img/avatar7-sm.jpg",
		"expire":        setting.ZConfig.Jwt.JwtExpire,
		"issuing_time":  time.Now().Unix(),
		"token":         token.Token,
		"refresh_token": refreshToken,
	}, nil
}

func (service *UserService) getUsrCache() {

}

type cacheUser struct {
	Id          uint
	AccessToken string
	Mobile      string
}

func setUserCache(uid uint, data cacheUser) {
	r, err := util.Rdb.HMSet(util.Ctx, util.CacheUserKey+fmt.Sprintf("%d", uid), data).Result()

	fmt.Println(err)
	fmt.Println(r)
}
