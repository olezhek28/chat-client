package root

import (
	"log"
	"os"
	"strings"

	"github.com/olezhek28/chat-client/internal/app"
	"github.com/olezhek28/chat-client/internal/model"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chat-client",
	Short: "Клиент лучшего в мире чата",
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Что-то создает",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Авторизует на сервере",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalf("failed to get username: %s\n", err.Error())
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			log.Fatalf("failed to get password: %s\n", err.Error())
		}

		serviceProvider := app.NewServiceProvider()
		handlerService := serviceProvider.GetHandlerService(ctx)

		err = handlerService.Login(ctx, &model.AuthInfo{
			Username: username,
			Password: password,
		})
		if err != nil {
			log.Fatalf("failed to login: %s\n", err.Error())
		}

		log.Println("login success")
	},
}

var createChatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Создает новый чат",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		usernamesStr, err := cmd.Flags().GetString("usernames")
		if err != nil {
			log.Fatalf("failed to get usernames: %s\n", err.Error())
		}

		usernames := strings.Split(usernamesStr, ",")
		if len(usernames) == 0 {
			log.Fatalf("usernames must be not empty")
		}

		serviceProvider := app.NewServiceProvider()
		handlerService := serviceProvider.GetHandlerService(ctx)

		chatID, err := handlerService.CreateChat(ctx, usernames)
		if err != nil {
			log.Fatalf("failed to create chat: %s\n", err.Error())
		}

		log.Printf("chat created with id: %s\n", chatID)
	},
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Подключается к чат-серверу",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		chatID, err := cmd.Flags().GetString("chat-id")
		if err != nil {
			log.Fatalf("failed to get chat id: %s\n", err.Error())
		}

		serviceProvider := app.NewServiceProvider()
		handlerService := serviceProvider.GetHandlerService(ctx)

		err = handlerService.ConnectChat(ctx, chatID)
		if err != nil {
			log.Fatalf("failed to connect: %s\n", err.Error())
		}

		log.Println("chat finished")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(connectCmd)
	createCmd.AddCommand(createChatCmd)

	loginCmd.Flags().StringP("username", "u", "", "Имя пользователя")
	err := loginCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark username flag as required: %s\n", err.Error())
	}

	loginCmd.Flags().StringP("password", "p", "", "Пароль пользователя")
	err = loginCmd.MarkFlagRequired("password")
	if err != nil {
		log.Fatalf("failed to mark password flag as required: %s\n", err.Error())
	}

	connectCmd.Flags().StringP("chat-id", "c", "", "Chat id")
	err = connectCmd.MarkFlagRequired("chat-id")
	if err != nil {
		log.Fatalf("failed to mark chat-id flag required: %s", err.Error())
	}

	createChatCmd.Flags().StringP("usernames", "u", "", "List of usernames for chat")
	err = createChatCmd.MarkFlagRequired("usernames")
	if err != nil {
		log.Fatalf("failed to mark usernames flag required: %s", err.Error())
	}
}
