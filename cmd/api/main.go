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

	"transaction-service/internal/config"
	"transaction-service/internal/database"
	"transaction-service/internal/handler"
	"transaction-service/internal/repository"
	"transaction-service/internal/service"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting transaction service...")

	dbConfig := config.NewDBConfig()

	db, err := database.Initialize(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepo)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	router := setupRouter(transactionHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	log.Printf("Server is starting on port 8080...")
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()
	log.Println("Server stopped gracefully")
}

func setupRouter(transactionHandler *handler.TransactionHandler) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(handler.Logger())
	router.Use(handler.ErrorHandler())
	router.Use(handler.CORS())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	transactionHandler.RegisterRoutes(router)

	return router
}
