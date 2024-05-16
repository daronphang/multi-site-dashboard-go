package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"multi-site-dashboard-go/internal"
	"multi-site-dashboard-go/internal/config"
	"multi-site-dashboard-go/internal/delivery/kafka"
	ws "multi-site-dashboard-go/internal/delivery/websocket"
	"multi-site-dashboard-go/internal/domain"
	"multi-site-dashboard-go/internal/repository"
	uc "multi-site-dashboard-go/internal/usecase"
	cv "multi-site-dashboard-go/internal/validator"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type TestBedServer struct {
	server *echo.Echo
	handler *RestHandler
}

var (
	syncOnce sync.Once
	tbs *TestBedServer
)

func ProvideTestBedServer(t *testing.T) *TestBedServer {
	// For integration testing.
	syncOnce.Do(func() {
		ctx := context.Background()

		cfg, err := config.ProvideConfig()
		if err != nil {
			t.Fatalf("error reading config file: %v", err)
		}

		// Setup server.
		e := echo.New()
		e.Validator = cv.ProvideValidator()

		// Migrate db.
		m, err := internal.WirePgMigrateInstance()
		if err != nil {
			t.Fatalf("error creating DB migration instance: %v", err)
		}
		if err := m.Up(); err != nil && err.Error() != "no change" {
			t.Fatalf("error migrating db: %v", err)
		}

		// Setup usecase and handler.
		db, err := internal.WirePgConnPool(ctx)
		if err != nil {
			t.Fatalf("error creating db pool: %v", err)
		}
		repo := repository.New(db)
		kw := kafka.New(cfg)
		ws := ws.New()
		uc := uc.NewUseCaseService(repo, kw, ws)
		rh := NewRestHandler(uc)

		tbs = &TestBedServer{
			server: e,
			handler: rh,
		}
	})
	return tbs
}

func TestCreateMachineResourceUsageRT(t *testing.T) {
	tbs := ProvideTestBedServer(t)

	t.Run("should return 200 success", func(t *testing.T) {
		// Setup request.
		bodyJSON := `{"machine":"testMachine","metric1":10,"metric2":50,"metric3":64}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(bodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := tbs.server.NewContext(req, rec)

		if assert.NoError(t, tbs.handler.CreateMachineResourceUsageRT(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			rv := new(domain.MachineResourceUsage)
			_ = json.Unmarshal(rec.Body.Bytes(), rv)
			assert.Equal(t, rv.Machine, "testMachine")
			assert.Equal(t, rv.Metric1, int32(10))
			assert.Equal(t, rv.Metric3, int32(64))
		}
	})

	t.Run("should return 400 if missing metric3", func(t *testing.T) {
		// Setup request.
		bodyJSON := `{"machine":"testMachine","metric1":10,"metric2":50}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(bodyJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := tbs.server.NewContext(req, rec)

		if assert.NoError(t, tbs.handler.CreateMachineResourceUsageRT(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "validation error")
		}
	})
}

func TestGetAggMachineResourceUsageRT(t *testing.T) {
	tbs := ProvideTestBedServer(t)

	t.Run("should return 200 success", func(t *testing.T) {
		// Setup request.
		q := make(url.Values)
		q.Set("lookBackPeriod", "1 day")
		q.Set("timeBucket", "1 day")

		req := httptest.NewRequest(http.MethodGet, "/?" + q.Encode(), nil)
		rec := httptest.NewRecorder()
		c := tbs.server.NewContext(req, rec)
		c.SetPath("/machines/:machine")
		c.SetParamNames("machine")
		c.SetParamValues("testMachine")

		if assert.NoError(t, tbs.handler.GetAggMachineResourceUsageRT(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			var rv []domain.AggMachineResourceUsage
			_ = json.Unmarshal(rec.Body.Bytes(), &rv)
			assert.Equal(t, len(rv), 1)
		}
	})

	t.Run("should return 400 error if missing query params", func(t *testing.T) {
		// Setup request.
		req := httptest.NewRequest(http.MethodGet, "/?", nil)
		rec := httptest.NewRecorder()
		c := tbs.server.NewContext(req, rec)
		c.SetPath("/machines/:machine")
		c.SetParamNames("machine")
		c.SetParamValues("testMachine")

		if assert.NoError(t, tbs.handler.GetAggMachineResourceUsageRT(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "validation error")
		}
	})
}