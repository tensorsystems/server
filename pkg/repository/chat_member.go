/*
  Copyright 2021 Kidus Tiliksew

  This file is part of Tensor EMR.

  Tensor EMR is free software: you can redistribute it and/or modify
  it under the terms of the version 2 of GNU General Public License as published by
  the Free Software Foundation.

  Tensor EMR is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package repository

import (
	"gorm.io/gorm"
)

// ChatMemeber ...
type ChatMember struct {
	gorm.Model
	ID          int     `gorm:"primaryKey"`
	UserID      int     `json:"userId"`
	ChatID      int     `json:"chatId"`
	DisplayName string  `json:"displayName"`
	PhotoURL    *string `json:"photoUrl"`
}

// Save ...
func (r *ChatMember) Save() error {
	return DB.Create(&r).Error
}

// Get ...
func (r *ChatMember) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// FindCommonChatID ...
func (r *ChatMember) FindCommonChatID(userID int, recipientID int) (int, error) {
	var chatMemeber ChatMember
	err := DB.Select("chat_id").Where("user_id = ? or user_id = ?", userID, recipientID).Group("chat_id").Having("count(chat_id) > ?", 1).Order("chat_id desc").Take(&chatMemeber).Error

	if err != nil {

		return 0, nil
	}

	return chatMemeber.ChatID, nil
}

// GetByChatID ...
func (r *ChatMember) GetByChatID(ID int) ([]*ChatMember, error) {
	var members []*ChatMember
	err := DB.Where("chat_id = ?", ID).Order("created_at desc").Find(&members).Error
	return members, err
}

// Update ...
func (r *ChatMember) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *ChatMember) Delete(userID int, chatID int) error {
	return DB.Where("user_id = ? AND chat_id = ?", userID, chatID).Delete(&r).Error
}
