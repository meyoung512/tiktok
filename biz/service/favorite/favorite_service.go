package service

import (
	"context"
	"offer_tiktok/biz/dal/db"
	favorite "offer_tiktok/biz/model/interact/favorite"
	"offer_tiktok/pkg/errno"

	"github.com/cloudwego/hertz/pkg/app"
)

const (
	FAVORITE   = 1
	UNFAVORITE = 2
)

type InteractService struct {
	ctx context.Context
	c   *app.RequestContext
}

// new InteractService
func NewInteractService(ctx context.Context, c *app.RequestContext) *InteractService {
	return &InteractService{ctx: ctx, c: c}
}

// like action, include like and unlike
// request parameters:
// string token = 1;       // 用户鉴权token
// int64 to_user_id = 2;   // 对方用户id
// int32 action_type = 3;  // 1-点赞，2-取消点赞
func (r *InteractService) FavoriteAction(req *favorite.DouyinFavoriteActionRequest) (flag bool, err error) {
	// 颁发和验证token的行为均交给jwt处理，当发送到handler层时，默认已通过验证
	// 只需要检查参数VideoID的合法性

	_, err = db.CheckVideoExistById(*req.VideoId) //zheli
	if err != nil {
		return false, err
	}
	if req.ActionType != FAVORITE && req.ActionType != UNFAVORITE {
		return false, errno.ParamErr
	}
	// 获取current_user_id
	current_user_id, _ := r.c.Get("current_user_id")
	// // 不准自己关注自己
	// if req.ToUserId == current_user_id.(int64) {
	// 	return false, errno.ParamErr
	// }
	new_favorite_relation := &db.Favorites{
		UserId:  current_user_id.(int64),
		VideoId: req.VideoId,
	}
	// 请求参数校验完毕，检查favorite表中是否已经存在这两者的关系
	favorite_exist, _ := db.CheckFavoriteRelationExist(new_favorite_relation)
	if req.ActionType == FAVORITE {
		if favorite_exist {
			return false, errno.FavoriteRelationAlreadyExistErr
		}
		flag, err = db.AddNewFavorite(new_favorite_relation)
	} else {
		if !favorite_exist {
			return false, errno.FavoriteRelationNotExistErr
		}
		flag, err = db.DeleteFavorite(new_favorite_relation)
	}
	return flag, err
}

// 获取用户点赞的所有视频列表，需要注意的是这里的token是客户端当前用户，而user_id可以是任意用户//zheli
// request parameters:
// string token;       // 用户鉴权token
// int64  user_id;     // 用户id
func (r *InteractService) GetFavoriteList(req *favorite.DouyinFavoriteListRequest) (favoritelist []favorite.Video, err error) {
	_, err = db.CheckUserExistById(req.UserId)
	if err != nil {
		return nil, err
	}
	// 获取current_user_id
	current_user_id, _ := r.c.Get("current_user_id")
	return db.GetFavoriteInfo(current_user_id.(int64), *req.UserId)
}