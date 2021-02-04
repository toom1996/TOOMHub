// @Description
// @Author: toom1996 <1023150697@qq.com>
// @dateTime: 2021/1/29 16:26
package service

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"
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
