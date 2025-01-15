package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"transaction-service/internal/config"
	"transaction-service/internal/database"
	"transaction-service/internal/handler"
	"transaction-service/internal/repository"
	"transaction-service/internal/service"
)

const (
	serverPort        = ":8080"
	shutdownTimeout   = 30 * time.Second
	dbMaxIdleConns    = 10
	dbMaxOpenConns    = 100
	dbConnMaxLifetime = time.Hour
)

func main() {
	initLogger()

	db := initDatabase()
	defer cleanupDatabase(db)

	app := initializeApp(db)

	server := setupServer(app.router)
	runServer(server)
}

func initLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting transaction service...")
}

func initDatabase() *gorm.DB {
	dbConfig := config.NewDBConfig()
	db, err := database.Initialize(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(dbMaxIdleConns)
	sqlDB.SetMaxOpenConns(dbMaxOpenConns)
	sqlDB.SetConnMaxLifetime(dbConnMaxLifetime)

	return db
}

func cleanupDatabase(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error getting DB instance during cleanup: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
}

type application struct {
	router             *gin.Engine
	transactionHandler *handler.TransactionHandler
}

func initializeApp(db *gorm.DB) *application {
	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepo)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	router := setupRouter(transactionHandler)

	return &application{
		router:             router,
		transactionHandler: transactionHandler,
	}
}

func setupRouter(transactionHandler *handler.TransactionHandler) *gin.Engine {
	router := gin.New()

	router.Use(
		gin.Recovery(),
		handler.Logger(),
		handler.ErrorHandler(),
		handler.CORS(),
	)

	router.GET("/health", healthCheckHandler)

	transactionHandler.RegisterRoutes(router)

	return router
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func setupServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    serverPort,
		Handler: handler,
	}
}

func runServer(server *http.Server) {
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go handleShutdown(server, serverCtx, serverStopCtx, sig)

	log.Printf("Server is starting on port%s...", serverPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()
	log.Println("Server stopped gracefully")
}

func handleShutdown(server *http.Server, serverCtx context.Context, serverStopCtx context.CancelFunc, sig chan os.Signal) {
	<-sig

	shutdownCtx, cancel := context.WithTimeout(serverCtx, shutdownTimeout)
	defer cancel()

	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			log.Fatal("graceful shutdown timed out.. forcing exit.")
		}
	}()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
	serverStopCtx()
}
