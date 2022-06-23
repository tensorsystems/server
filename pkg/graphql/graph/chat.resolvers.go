package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/tensoremr/server/pkg/graphql/graph/generated"
	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/middleware"
	"github.com/tensoremr/server/pkg/repository"
)

func (r *mutationResolver) CreateChat(ctx context.Context, input model.ChatInput) (*repository.Chat, error) {
	// Get current user
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var user repository.User
	if err := user.GetByEmail(email); err != nil {
		return nil, err
	}

	var recipient repository.User
	if err := recipient.Get(input.RecipientID); err != nil {
		return nil, err
	}

	var chatMember repository.ChatMember
	chatID, err := chatMember.FindCommonChatID(user.ID, recipient.ID)
	if err != nil {
		return nil, err
	}

	var chat repository.Chat
	chat.RecentMessage = input.Message

	if chatID != 0 {
		if err := chat.Get(chatID); err != nil {
			return nil, err
		}
	} else {
		chat = repository.Chat{
			RecentMessage: input.Message,
			ChatMembers: []repository.ChatMember{
				{UserID: user.ID, DisplayName: user.FirstName + " " + user.LastName},
				{UserID: recipient.ID, DisplayName: recipient.FirstName + " " + recipient.LastName},
			},
		}
	}

	if err := chat.Save(); err != nil {
		return nil, err
	}

	// Save message
	var chatMessage repository.ChatMessage
	chatMessage.Body = input.Message
	chatMessage.ChatID = chat.ID
	chatMessage.UserID = user.ID
	if err := chatMessage.Save(); err != nil {
		return nil, err
	}

	// Save unread message
	var unreadMessage repository.ChatUnreadMessage
	unreadMessage.ChatID = chat.ID
	unreadMessage.UserID = recipient.ID
	if err := unreadMessage.Save(); err != nil {
		return nil, err
	}

	return &chat, nil
}

func (r *mutationResolver) SendMessage(ctx context.Context, input model.ChatMessageInput) (*repository.ChatMessage, error) {
	// Get current user
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var user repository.User
	if err := user.GetByEmail(email); err != nil {
		return nil, err
	}

	// Save message
	var chatMessage repository.ChatMessage
	chatMessage.Body = input.Body
	chatMessage.ChatID = input.ChatID
	chatMessage.UserID = user.ID
	if err := chatMessage.Save(); err != nil {
		return nil, err
	}

	// Update chat
	var chat repository.Chat
	chat.ID = input.ChatID
	chat.RecentMessage = chatMessage.Body
	if err := chat.Update(); err != nil {
		return nil, err
	}

	var chatMemberRepo repository.ChatMember
	chatMembers, err := chatMemberRepo.GetByChatID(input.ChatID)
	if err != nil {
		return nil, err
	}

	var recipient *repository.ChatMember
	for _, member := range chatMembers {
		if member.UserID != user.ID {
			recipient = member
		}
	}

	// Save unread message
	var unreadMessage repository.ChatUnreadMessage
	unreadMessage.ChatID = input.ChatID
	unreadMessage.UserID = recipient.UserID
	if err := unreadMessage.Save(); err != nil {
		return nil, err
	}

	return &chatMessage, nil
}

func (r *mutationResolver) MuteChat(ctx context.Context, id int) (*repository.ChatMute, error) {
	// Get current user
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var user repository.User
	if err := user.GetByEmail(email); err != nil {
		return nil, err
	}

	var chatMute repository.ChatMute
	chatMute.ChatID = id
	chatMute.UserID = user.ID

	if err := chatMute.Save(); err != nil {
		return nil, err
	}

	return &chatMute, nil
}

func (r *mutationResolver) UnmuteChat(ctx context.Context, id int) (bool, error) {
	// Get current user
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return false, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return false, errors.New("Cannot find user")
	}

	var user repository.User
	if err := user.GetByEmail(email); err != nil {
		return false, err
	}

	var chatMute repository.ChatMute
	if err := chatMute.Delete(user.ID, id); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) DeleteChat(ctx context.Context, id int) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteUnreadMessages(ctx context.Context, userID int, chatID int) (bool, error) {
	var repository repository.ChatUnreadMessage
	if err := repository.DeleteForUserChat(userID, chatID); err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) GetChat(ctx context.Context) (*repository.Chat, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetUserChats(ctx context.Context) ([]*repository.Chat, error) {
	// Get current user
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var user repository.User
	if err := user.GetByEmail(email); err != nil {
		return nil, err
	}

	var chatRepository repository.Chat
	results, err := chatRepository.GetUserChats(user.ID)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *queryResolver) GetChatMessages(ctx context.Context, id int) ([]*repository.ChatMessage, error) {
	// Get current user
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var user repository.User
	if err := user.GetByEmail(email); err != nil {
		return nil, err
	}

	var chatMessage repository.ChatMessage
	results, err := chatMessage.GetByChatID(id)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *queryResolver) GetChatMembers(ctx context.Context, id int) ([]*repository.ChatMember, error) {
	// Get current user
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var user repository.User
	if err := user.GetByEmail(email); err != nil {
		return nil, err
	}

	var chatMember repository.ChatMember
	results, err := chatMember.GetByChatID(id)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *queryResolver) GetUnreadMessages(ctx context.Context) ([]*repository.ChatUnreadMessage, error) {
	// Get current user
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var user repository.User
	if err := user.GetByEmail(email); err != nil {
		return nil, err
	}

	var unreadMessages repository.ChatUnreadMessage
	results, err := unreadMessages.GetByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *queryResolver) GetCommonChat(ctx context.Context, recipientID int) (*repository.Chat, error) {
	// Get current user
	gc, err := middleware.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email := gc.GetString("email")
	if len(email) == 0 {
		return nil, errors.New("Cannot find user")
	}

	var user repository.User
	if err := user.GetByEmail(email); err != nil {
		return nil, err
	}

	var recipient repository.User
	if err := recipient.Get(recipientID); err != nil {
		return nil, err
	}

	var chatMember repository.ChatMember
	chatID, err := chatMember.FindCommonChatID(user.ID, recipient.ID)
	if err != nil {
		return nil, err
	}

	var chat repository.Chat
	if err := chat.Get(chatID); err != nil {
		return nil, err
	}

	return &chat, nil
}

func (r *subscriptionResolver) Notification(ctx context.Context) (<-chan *model.Notification, error) {
	panic(fmt.Errorf("not implemented"))
}

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
