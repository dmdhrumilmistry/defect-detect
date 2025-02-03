package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/dmdhrumilmistry/defect-detect/pkg/config"
	"github.com/dmdhrumilmistry/defect-detect/pkg/db"
	"github.com/dmdhrumilmistry/defect-detect/pkg/service/auth"
	"github.com/dmdhrumilmistry/defect-detect/pkg/types"
	"github.com/dmdhrumilmistry/defect-detect/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func createUser(authStore *auth.AuthStore, name string, email string, isSuperUser bool) {
	if name == "" {
		log.Fatal().Msg("invalid name")
	}

	// validate email
	if !utils.IsValidEmail(email) {
		log.Fatal().Msg("Invalid email")
	}

	now := time.Now()
	user := types.User{
		Name:        name,
		Email:       email,
		AvatarUrl:   "https://lh3.googleusercontent.com/-XdUIqdMkCWA/AAAAAAAAAAI/AAAAAAAAAAA/4252rscbv5M/photo.jpg",
		IsActive:    true,
		IsSuperUser: isSuperUser,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	user_id, err := authStore.CreateUser(user)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to create user")
		return
	}

	log.Info().Msgf("User created successfully: %s", user_id)
}

func createToken(userId, email string, authStore *auth.AuthStore) {
	var (
		user types.User
		err  error
	)

	if userId != "" {
		user, err = authStore.GetUserById(userId, config.DefaultConfig.DbQueryTimeout)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to fetch user by id: %s", userId)
		}
	} else if utils.IsValidEmail(email) {
		user, err = authStore.GetUserByEmail(email, config.DefaultConfig.DbQueryTimeout)
		if err != nil {
			log.Fatal().Err(err).Msgf("failed to fetch user by email: %s", email)
		}
	}
	token, err := auth.CreateJWT(user.Id)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create token for user: %s", user.Id)
	}
	log.Info().Msg(token)
}

func main() {
	// Check if at least one argument is provided
	if len(os.Args) < 2 {
		log.Fatal().Msg("valid subcommand 'user'/'token'")
	}

	mgo, err := db.NewMongo(config.DefaultConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get db connection")
	}
	defer mgo.Client.Disconnect(context.TODO())
	authStore := auth.NewAuthStore(mgo.Db)

	if !config.DefaultConfig.IsDevEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	subcommand := os.Args[1]
	args := os.Args[2:]

	userFlag := flag.NewFlagSet("user", flag.ExitOnError)
	name := userFlag.String("name", "", "user full name")
	email := userFlag.String("email", "", "user email")
	superUser := userFlag.Bool("superUser", false, "provides superuser perms if value is true. Default value: false")

	tokenFlag := flag.NewFlagSet("token", flag.ExitOnError)
	userId := tokenFlag.String("id", "", "mongo db user id")
	userEmail := tokenFlag.String("email", "", "user email id")

	switch subcommand {
	case "user":
		userFlag.Parse(args)
		createUser(authStore, *name, *email, *superUser)
	case "token":
		tokenFlag.Parse(args)
		log.Info().Msgf("%s %s", *userId, *userEmail)

		createToken(*userId, *userEmail, authStore)

	default:
		log.Fatal().Msgf("invalid command: %s", subcommand)
	}

}
