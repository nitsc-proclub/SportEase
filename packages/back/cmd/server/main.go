// @title GYOUJI_HP API
// @version 1.0
// @description This is a sample server for GYOUJI_HP.
// @host localhost:3300
// @BasePath /
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/saku0512/GYOUJI_HP/backend/internal/handler"
	"github.com/saku0512/GYOUJI_HP/backend/internal/repository"
	"github.com/saku0512/GYOUJI_HP/backend/internal/router"
	"github.com/saku0512/GYOUJI_HP/backend/internal/service"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	log.Println("Starting the application...")
	// Database connection
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_ROOT_PASSWORD"),
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		DBName:               os.Getenv("DB_DATABASE"),
		AllowNativePasswords: true, // MySQL 8.0以上で推奨
		ParseTime:            true,
		Collation:            "utf8mb4_unicode_ci",
		// charsetパラメータをParams経由で確実に設定する
		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}

	// cfg.FormatDSN() で安全なDSN文字列を生成
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to the database")

	// Initialize users
	initializeUsers(db)

	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET_KEY is not set")
	}

	// User用DI
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, jwtSecret)
	userHandler := handler.NewUserHandler(userService)

	// Tournament用DI
	tournamentRepository := repository.NewTournamentRepository(db)
	tournamentService := service.NewTournamentService(tournamentRepository)
	tournamentHandler := handler.NewTournamentHandler(tournamentService)

	// Match用DI
	matchRepository := repository.NewMatchRepository(db)
	matchService := service.NewMatchService(matchRepository)
	matchHandler := handler.NewMatchHandler(matchService)

	// Score用DI
	scoreRepository := repository.NewScoreRepository(db)
	scoreService := service.NewScoreService(scoreRepository)
	scoreHandler := handler.NewScoreHandler(scoreService)

	// Relay用DI
	relayRepository := repository.NewRelayRepository(db)
	relayService := service.NewRelayService(relayRepository)
	relayHandler := handler.NewRelayHandler(relayService)

	// Attendance用DI
	attendanceRepository := repository.NewAttendanceRepository(db)
	attendanceService := service.NewAttendanceService(attendanceRepository)
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)

	// Setting用DI
	settingRepository := repository.NewSettingRepository(db)
	settingService := service.NewSettingService(settingRepository)
	settingHandler := handler.NewSettingHandler(settingService)

	r := router.SetupRouter(userHandler, tournamentHandler, matchHandler, jwtSecret, scoreHandler, relayHandler, attendanceHandler, settingHandler)
	log.Println("Server is running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func initializeUsers(db *sql.DB) {
	users := []map[string]string{
		{"username": os.Getenv("SUPERROOT_USER"), "password": os.Getenv("SUPERROOT_PASS"), "role": "superroot", "assigned_sport": ""},
		{"username": os.Getenv("ADMIN_TABLE_TENNIS_USER"), "password": os.Getenv("ADMIN_TABLE_TENNIS_PASS"), "role": "admin", "assigned_sport": "table_tennis"},
		{"username": os.Getenv("ADMIN_VOLLEYBALL_USER"), "password": os.Getenv("ADMIN_VOLLEYBALL_PASS"), "role": "admin", "assigned_sport": "volleyball"},
		{"username": os.Getenv("ADMIN_SOCCER_USER"), "password": os.Getenv("ADMIN_SOCCER_PASS"), "role": "admin", "assigned_sport": "soccer"},
		{"username": os.Getenv("STUDENT_USER"), "password": os.Getenv("STUDENT_PASS"), "role": "student", "assigned_sport": ""},
		{"username": os.Getenv("ADMIN_RELAY_USER"), "password": os.Getenv("ADMIN_RELAY_PASS"), "role": "admin", "assigned_sport": "relay"},
	}

	for _, user := range users {
		log.Printf("Processing user: %s", user["username"])
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", user["username"]).Scan(&exists)
		if err != nil {
			log.Printf("failed to check if user %s exists: %v", user["username"], err)
			continue
		}

		log.Printf("User %s exists: %t", user["username"], exists)

		if !exists {
			log.Printf("User %s does not exist, creating...", user["username"])
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user["password"]), bcrypt.DefaultCost)
			if err != nil {
				log.Printf("failed to hash password for user %s: %v", user["username"], err)
				continue
			}

			_, err = db.Exec("INSERT INTO users (username, hashed_password, role, assigned_sport) VALUES (?, ?, ?, ?)",
				user["username"], string(hashedPassword), user["role"], user["assigned_sport"])
			if err != nil {
				log.Printf("failed to insert user %s: %v", user["username"], err)
				continue
			}
			fmt.Printf("User %s created successfully.\n", user["username"])
		}
	}
}
