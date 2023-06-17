package root

import (
	"log"
	"os"

	"github.com/olezhek28/chat-client/internal/app"
	"github.com/olezhek28/chat-client/internal/model"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chat-client",
	Short: "Клиент лучшего в мире чата",
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

	loginCmd.Flags().StringP("username", "u", "", "Имя пользователя")
	err := loginCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark username flag as required: %s\n", err.Error())
	}
}
