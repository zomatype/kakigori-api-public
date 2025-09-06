package main

import (
	"context"
	"fmt"
	"kakigori-api/internal/config"
	"kakigori-api/internal/domain"
	"kakigori-api/internal/handler"
	"kakigori-api/internal/middleware"
	"kakigori-api/internal/middleware/cors"
	"kakigori-api/internal/middleware/ratelimit"
	"kakigori-api/internal/repository"
	"kakigori-api/internal/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	infoLogger := log.New(os.Stdout, "", log.LstdFlags)
	errorLogger := log.New(os.Stderr, "", log.LstdFlags)
	logger := middleware.NewLogger(infoLogger, errorLogger)

	// fly.io対応: PORT優先、なければKAKIGORI_PORT、なければ8080
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("KAKIGORI_PORT")
		if port == "" {
			port = "8080"
		}
	}

	storeIDs := config.GetStoreIDs()
	if len(storeIDs) == 0 {
		logger.LogError("No store IDs configured. Set KAKIGORI_STORE_IDS environment variable.")
		os.Exit(1)
	}

	logger.LogInfo(fmt.Sprintf("Found %d store(s):", len(storeIDs)))
	for _, storeID := range storeIDs {
		logger.LogInfo("  " + storeID)
	}

	maxRequests := 10
	if maxReqStr := os.Getenv("KAKIGORI_STORE_MAX_REQUESTS"); maxReqStr != "" {
		if max, err := strconv.Atoi(maxReqStr); err == nil && max > 0 {
			maxRequests = max
		}
	}

	storeRepo := repository.NewMemoryStoreRepository()
	orderRepo := repository.NewMemoryOrderRepository()
	for _, storeID := range storeIDs {
		store := domain.NewStore(storeID, config.GetMaxOrders())
		storeRepo.Create(store)
		orderRepo.InitializeStore(store)
	}
	menuRepo := repository.NewMemoryMenuRepository()

	menuUsecase := usecase.NewMenuUsecase(menuRepo)
	orderUsecase := usecase.NewOrderUsecase(menuRepo, orderRepo, storeRepo)

	router := handler.NewRouter(menuUsecase, orderUsecase, logger)

	corsMiddleware := cors.Middleware(cors.DefaultConfig())
	storeLimiter := ratelimit.NewIpRateLimiter(maxRequests, 1*time.Second, 3*time.Hour)

	mux := http.NewServeMux()
	router.Setup(mux)

	var handler http.Handler = mux
	handler = corsMiddleware(handler)
	handler = storeLimiter.Middleware(handler)
	handler = logger.Middleware(handler)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// サーバーの起動
	go func() {
		logger.LogInfo(fmt.Sprintf("Server starting on port %s", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.LogError(fmt.Sprintf("Server failed to start: %v", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.LogInfo("Server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.LogError(fmt.Sprintf("Server forced to shutdown: %v", err))
		os.Exit(1)
	}

	logger.LogInfo("Server stopped gracefully")
}
