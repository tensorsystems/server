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

import "gorm.io/gorm"

// Chat ...
type Chat struct {
	gorm.Model
	ID                 int                 `gorm:"primaryKey"`
	RecentMessage      string              `json:"recentMessage"`
	ChatMembers        []ChatMember        `json:"chatMembers"`
	ChatMutes          []ChatMute          `json:"chatMutes"`
	ChatDeletes        []ChatDelete        `json:"chatDeletes"`
	ChatUnreadMessages []ChatUnreadMessage `json:"chatUnreadMessages"`
	ChatMessages       []ChatMessage       `json:"chatMessages"`
}

// Save ...
func (r *Chat) Save() error {
	return DB.Create(&r).Error
}

// Get ...
func (r *Chat) Get(ID int) error {
	return DB.Where("id = ?", ID).Preload("ChatMembers").Take(&r).Error
}

// GetUserChats ...
func (r *Chat) GetUserChats(userID int) ([]*Chat, error) {
	var chats []*Chat
	err := DB.Joins("inner join chat_members on chat_members.chat_id = chats.id").Where("chat_members.user_id = ?", userID).Order("updated_at desc").Preload("ChatMembers").Preload("ChatMutes").Preload("ChatUnreadMessages").Find(&chats).Error
	return chats, err
}

// Update ...
func (r *Chat) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *Chat) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
