package handler

import "context"

func (h *Handler) CreateChat(ctx context.Context, usernames []string) (string, error) {
	chatID, err := h.chatClient.CreateChat(ctx, usernames)
	if err != nil {
		return "", err
	}

	return chatID, nil
}
