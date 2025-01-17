// Code generated by hertz generator.

package user

import (
	"context"
	user "offer_tiktok/biz/model/basic/user"
	"offer_tiktok/biz/mw/jwt"
	"offer_tiktok/biz/pack"
	service "offer_tiktok/biz/service/user"
	"offer_tiktok/pkg/errno"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// UserRegister .
// @router /douyin/user/register/  [POST]
func UserRegister(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.DouyinUserRegisterRequest
	err = c.BindAndValidate(&req)
	hlog.CtxInfof(ctx, "OK")
	if err != nil {
		resp := pack.BuildBaseResp(err)
		c.JSON(consts.StatusOK, user.DouyinUserRegisterResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		})
		return
	}

	_, err = service.NewUserService(ctx, c).UserRegister(&req)
	if err != nil {
		resp := pack.BuildBaseResp(err)
		c.JSON(consts.StatusOK, user.DouyinUserRegisterResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		})
		return
	}
	jwt.JwtMiddleware.LoginHandler(ctx, c)
	token := c.GetString("token")
	v, _ := c.Get("user_id")
	user_id := v.(int64)
	c.JSON(consts.StatusOK, user.DouyinUserRegisterResponse{
		StatusCode: errno.SuccessCode,
		StatusMsg:  errno.SuccessMsg,
		Token:      token,
		UserId:     user_id,
	})
}

// UserLogin .
// @router /douyin/user/login/  [POST]
func UserLogin(ctx context.Context, c *app.RequestContext) {
	v, _ := c.Get("user_id")
	user_id := v.(int64)
	token := c.GetString("token")
	c.JSON(consts.StatusOK, user.DouyinUserLoginResponse{
		StatusCode: errno.SuccessCode,
		StatusMsg:  errno.SuccessMsg,
		Token:      token,
		UserId:     user_id,
	})
}

// User .
// @router /douyin/user/ [GET]
func User(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.DouyinUserRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		resp := pack.BuildBaseResp(err)
		c.JSON(consts.StatusOK, user.DouyinUserResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		})
		return
	}

	u, err := service.NewUserService(ctx, c).UserInfo(&req)

	resp := pack.BuildBaseResp(err)
	c.JSON(consts.StatusOK, user.DouyinUserResponse{
		StatusCode: resp.StatusCode,
		StatusMsg:  resp.StatusMsg,
		User:       u,
	})
}
