package tests

import (
	"bytes"
	"encoding/json"
	"kakigori-api/internal/domain"
	"kakigori-api/internal/handler"
	"kakigori-api/internal/middleware"
	"kakigori-api/internal/middleware/cors"
	"kakigori-api/internal/repository"
	"kakigori-api/internal/usecase"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestFullWorkflow(t *testing.T) {
	os.Setenv("KAKIGORI_STORE_IDS", "kakigori-team-001,kakigori-team-002")
	os.Setenv("KAKIGORI_MAX_ORDERS", "10")
	defer func() {
		os.Unsetenv("KAKIGORI_STORE_IDS")
		os.Unsetenv("KAKIGORI_MAX_ORDERS")
	}()

	infoLogger := log.New(os.Stdout, "", log.LstdFlags)
	errorLogger := log.New(os.Stderr, "", log.LstdFlags)
	logger := middleware.NewLogger(infoLogger, errorLogger)

	menuRepo := repository.NewMemoryMenuRepository()
	orderRepo := repository.NewMemoryOrderRepository()
	storeRepo := repository.NewMemoryStoreRepository()
	// テスト用ストアIDを事前登録
	storeRepo.Create(&domain.Store{ID: "kakigori-team-001", MaxOrders: 10})
	storeRepo.Create(&domain.Store{ID: "kakigori-team-002", MaxOrders: 10})
	menuUsecase := usecase.NewMenuUsecase(menuRepo)
	orderUsecase := usecase.NewOrderUsecase(menuRepo, orderRepo, storeRepo)
	router := handler.NewRouter(menuUsecase, orderUsecase, logger)

	corsMiddleware := cors.Middleware(cors.DefaultConfig())

	mux := http.NewServeMux()
	router.Setup(mux)

	var httpHandler http.Handler = mux
	httpHandler = corsMiddleware(httpHandler)
	httpHandler = logger.Middleware(httpHandler)

	server := httptest.NewServer(httpHandler)
	defer server.Close()

	t.Run("メニュー取得", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/v1/stores/kakigori-team-001/menu")
		if err != nil {
			t.Fatalf("Menu request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.Status)
		}

		var menuResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&menuResp); err != nil {
			t.Fatalf("Failed to decode menu response: %v", err)
		}

		menu, ok := menuResp["menu"].([]interface{})
		if !ok {
			t.Fatal("Menu response does not contain menu array")
		}

		if len(menu) != 4 {
			t.Errorf("Expected 4 menu items, got %d", len(menu))
		}
	})

	t.Run("注文作成と取得", func(t *testing.T) {
		orderReq := map[string]interface{}{
			"menu_item_id": "giiku-sai",
		}
		body, _ := json.Marshal(orderReq)

		resp, err := http.Post(
			server.URL+"/v1/stores/kakigori-team-001/orders",
			"application/json",
			bytes.NewReader(body),
		)
		if err != nil {
			t.Fatalf("Order creation failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status Created, got %v", resp.Status)
		}

		var orderResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
			t.Fatalf("Failed to decode order response: %v", err)
		}

		orderID := orderResp["id"].(string)

		time.Sleep(100 * time.Millisecond)
		resp, err = http.Get(server.URL + "/v1/stores/kakigori-team-001/orders")
		if err != nil {
			t.Fatalf("Get orders request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.Status)
		}

		var ordersResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&ordersResp); err != nil {
			t.Fatalf("Failed to decode orders response: %v", err)
		}

		orders, ok := ordersResp["orders"].([]interface{})
		if !ok {
			t.Fatal("Orders response does not contain orders array")
		}

		if len(orders) != 1 {
			t.Errorf("Expected 1 order, got %d", len(orders))
		}

		// 注文を受取待ち状態に更新
		req, _ := http.NewRequest(
			http.MethodPost,
			server.URL+"/v1/stores/kakigori-team-001/orders/"+orderID+"/waiting-pickup",
			nil,
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Order waiting-pickup failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.Status)
		}

		// 注文を完了状態に更新
		req, _ = http.NewRequest(
			http.MethodPost,
			server.URL+"/v1/stores/kakigori-team-001/orders/"+orderID+"/complete",
			nil,
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Order complete failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK, got %v", resp.Status)
		}
	})

	t.Run("FIFOキューの動作確認", func(t *testing.T) {
		os.Setenv("KAKIGORI_MAX_ORDERS", "2")
		defer os.Setenv("KAKIGORI_MAX_ORDERS", "10")

		fifoMenuRepo := repository.NewMemoryMenuRepository()
		fifoOrderRepo := repository.NewMemoryOrderRepository()
		fifoStoreRepo := repository.NewMemoryStoreRepository()
		fifoStoreRepo.Create(&domain.Store{ID: "kakigori-team-001", MaxOrders: 2})
		fifoMenuUsecase := usecase.NewMenuUsecase(fifoMenuRepo)
		fifoOrderUsecase := usecase.NewOrderUsecase(fifoMenuRepo, fifoOrderRepo, fifoStoreRepo)
		fifoRouter := handler.NewRouter(fifoMenuUsecase, fifoOrderUsecase, logger)

		fifoMux := http.NewServeMux()
		fifoRouter.Setup(fifoMux)

		var fifoHandler http.Handler = fifoMux
		fifoHandler = corsMiddleware(fifoHandler)
		fifoHandler = logger.Middleware(fifoHandler)

		fifoServer := httptest.NewServer(fifoHandler)
		defer fifoServer.Close()

		orderReq := map[string]interface{}{"menu_item_id": "giiku-sai"}
		body, _ := json.Marshal(orderReq)

		resp1, _ := http.Post(
			fifoServer.URL+"/v1/stores/kakigori-team-001/orders",
			"application/json",
			bytes.NewReader(body),
		)
		var order1 map[string]interface{}
		json.NewDecoder(resp1.Body).Decode(&order1)
		resp1.Body.Close()

		resp2, _ := http.Post(
			fifoServer.URL+"/v1/stores/kakigori-team-001/orders",
			"application/json",
			bytes.NewReader(body),
		)
		resp2.Body.Close()

		resp3, _ := http.Post(
			fifoServer.URL+"/v1/stores/kakigori-team-001/orders",
			"application/json",
			bytes.NewReader(body),
		)
		resp3.Body.Close()

		resp, err := http.Get(fifoServer.URL + "/v1/stores/kakigori-team-001/orders")
		if err != nil {
			t.Fatalf("Get orders request failed: %v", err)
		}
		defer resp.Body.Close()

		var ordersResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&ordersResp); err != nil {
			t.Fatalf("Failed to decode orders response: %v", err)
		}

		orders, ok := ordersResp["orders"].([]interface{})
		if !ok {
			t.Fatal("Orders response does not contain orders array")
		}

		if len(orders) != 2 {
			t.Errorf("Expected 2 orders in queue, got %d", len(orders))
		}

		for _, order := range orders {
			orderMap := order.(map[string]interface{})
			if orderMap["id"] == order1["id"] {
				t.Error("First order should have been removed due to FIFO")
			}
		}
	})
}
