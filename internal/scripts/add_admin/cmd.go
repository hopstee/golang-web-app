package add_admin

import (
	"fmt"
	"log"
	"log/slog"
	"mobile-backend-boilerplate/internal/config"
	"mobile-backend-boilerplate/internal/repository"
	"mobile-backend-boilerplate/internal/repository/sqlite"
	cmdHelper "mobile-backend-boilerplate/pkg/helper/cmd"
	passwordHelper "mobile-backend-boilerplate/pkg/helper/password"
	customLogger "mobile-backend-boilerplate/pkg/logger"
	"strings"

	"github.com/spf13/cobra"
)

var logger = customLogger.New(slog.LevelDebug)

var Command = &cobra.Command{
	Use:   "create",
	Short: "Create a new admin",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		force, _ := cmd.Flags().GetBool("force")

		if username == "" {
			username = cmdHelper.PromtLine("Username: ")
		}

		if password == "" {
			pw1 := cmdHelper.PromtPassword("Password: ")
			pw2 := cmdHelper.PromtPassword("Confirm password: ")
			if pw1 != pw2 {
				log.Fatal("Passwords do not match")
			}
			password = pw1
		}

		if len(password) < 6 {
			log.Fatal("Passwords must be at least 6 characters")
		}

		hash, err := passwordHelper.HashPassword(password)
		if err != nil {
			log.Fatalf("bcrypt error: %v", err)
		}

		config, err := config.LoadConfig("config/config.yml")
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		var repo repository.Repository
		var adminRepo repository.AdminRepository

		switch config.Database.Driver {
		case "sqlite":
			repo, err = sqlite.NewSQLiteRepository(config.Database.DataSource, logger)
			if err != nil {
				log.Fatalf("failed to init sqlite repository: %v", err)
			}

			db := repo.(*sqlite.SQLiteRepository).DB
			adminRepo = sqlite.NewAdminRepo(db, logger)

			fmt.Println("SQLite repository initialized with DSN:", config.Database.DataSource)
		default:
			log.Fatalf("unsupported database driver: %s", config.Database.Driver)
		}

		newAdmin := repository.Admin{
			Username: username,
			Password: hash,
			Role:     "admin",
		}

		id, err := adminRepo.Create(newAdmin)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "UNIQUE") && force {
				newAdmin.ID = id
				err = adminRepo.Update(newAdmin)
				if err != nil {
					log.Fatalf("admin update failed: %v", err)
				}
				log.Printf("password for %s updated.\n", username)
				return
			}
			log.Fatalf("admin create failed: %v", err)
		}

		log.Printf("admin %s created successfully.\n", username)
	},
}

func init() {
	Command.Flags().String("username", "", "Admin username")
	Command.Flags().String("password", "", "Admin password (unsafe, better use prompt)")
	Command.Flags().Bool("force", false, "Update password if admin already exists")
}
