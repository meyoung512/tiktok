package db

import (
	"offer_tiktok/pkg/constants"
	"offer_tiktok/pkg/errno"
	"time"
)

type Messages struct {
	ID         int64     `json:"id"`
	ToUserId   int64     `json:"to_user_id"`
	FromUserId int64     `json:"from_user_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

func (m *Messages) TableName() string {
	return constants.MessageTableName
}

func AddNewMessage(message *Messages) (bool, error) {
	exist, err := QueryUserById(message.FromUserId)
	if exist == nil || err != nil {
		return false, errno.UserIsNotExistErr
	}
	exist, err = QueryUserById(message.ToUserId)
	if exist == nil || err != nil {
		return false, errno.UserIsNotExistErr
	}
	err = DB.Create(message).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetMessageByIdPair(user_id1 int64, user_id2 int64, pre_msg_time time.Time) ([]Messages, error) {
	exist, err := QueryUserById(user_id1)
	if exist == nil || err != nil {
		return nil, errno.UserIsNotExistErr
	}
	exist, err = QueryUserById(user_id2)
	if exist == nil || err != nil {
		return nil, errno.UserIsNotExistErr
	}

	var messages []Messages
	err = DB.Where("to_user_id = ? AND from_user_id = ? AND created_at > ?", user_id1, user_id2, pre_msg_time).Or("to_user_id = ? AND from_user_id = ? AND created_at > ?", user_id2, user_id1, pre_msg_time).Find(&messages).Error

	if err != nil {
		return nil, err
	}
	return messages, nil
}

func GetLatestMessageByIdPair(user_id1 int64, user_id2 int64) (*Messages, error) {
	exist, err := QueryUserById(user_id1)
	if exist == nil || err != nil {
		return nil, errno.UserIsNotExistErr
	}
	exist, err = QueryUserById(user_id2)
	if exist == nil || err != nil {
		return nil, errno.UserIsNotExistErr
	}

	var message Messages
	err = DB.Where("to_user_id = ? AND from_user_id = ?", user_id1, user_id2).Or("to_user_id = ? AND from_user_id = ?", user_id2, user_id1).Last(&message).Error
	if err == nil {
		return &message, nil
	} else {
		if err.Error() == "record not found" {
			return nil, nil
		} else {
			return nil, err
		}
	}
}
