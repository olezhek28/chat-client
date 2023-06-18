package handler

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/olezhek28/chat-client/internal/model"
)

func (h *Handler) ConnectChat(ctx context.Context, chatID string) error {
	username, err := h.redisClient.Get(model.UsernameKey)
	if err != nil {
		return err
	}

	stream, err := h.chatClient.ConnectChat(ctx, chatID, username)
	if err != nil {
		return err
	}

	go func() {
		for {
			message, errRecv := stream.Recv()
			if errRecv == io.EOF {
				return
			}
			if errRecv != nil {
				log.Println("failed to receive message from stream: ", errRecv)
				return
			}

			log.Printf("[%v] - [from: %s]: %s\n", message.GetCreatedAt(), message.GetFrom(), message.GetText())
		}
	}()

	for {
		scanner := bufio.NewScanner(os.Stdin)
		var lines strings.Builder

		for {
			scanner.Scan()
			line := scanner.Text()
			if len(line) == 0 {
				break
			}

			lines.WriteString(line)
			lines.WriteString("\n")
		}

		err = scanner.Err()
		if err != nil {
			log.Println("failed to scan message: ", err)
		}

		err = h.chatClient.SendMessage(ctx, chatID, &model.Message{
			From:      username,
			Text:      lines.String(),
			CreatedAt: time.Now(),
		})
		if err != nil {
			log.Println("failed to send message: ", err)
			return err
		}
	}
}
