package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"multi-site-dashboard-go/internal"
	"multi-site-dashboard-go/internal/config"
	"multi-site-dashboard-go/internal/database"
	"multi-site-dashboard-go/internal/delivery/kafka"
	ws "multi-site-dashboard-go/internal/delivery/websocket"
	"multi-site-dashboard-go/internal/domain"
	"multi-site-dashboard-go/internal/repository"
	uc "multi-site-dashboard-go/internal/usecase"
	cv "multi-site-dashboard-go/internal/validator"

	"github.com/jackc/pgx/v5"
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

func CleanupDatabase(cfg *config.Config) error {
	conn, err := database.ProvidePgConn(cfg, true)
	if err != nil {
		return err
	}

	_, err = conn.Exec("DROP DATABASE IF EXISTS " + pgx.Identifier{cfg.Postgres.DBName}.Sanitize())
	if err != nil {
		return err
	}
	return nil
}

func ProvideTestBedServer(t *testing.T) (*TestBedServer, error) {
	// For integration testing.
	var syncErr error
	syncOnce.Do(func() {
		ctx := context.Background()

		cfg, err := config.ProvideConfig()
		if err != nil {
			syncErr = fmt.Errorf("error reading config file: %w", err)
			return
		}

		// setup db.
		if err := database.SetupDatabase(cfg); err != nil {
			syncErr = fmt.Errorf("error setting up db: %w", err)
			return
		}

		// Migrate db.
		m, err := internal.WirePgMigrateInstance(true)
		if err != nil {
			syncErr = fmt.Errorf("error creating db migration instance: %w", err)
			return
		}
		if err := m.Up(); err != nil && err.Error() != "no change" {
			syncErr = fmt.Errorf("error migration db: %w", err)
			return
		}

		// Setup usecase and handler.
		db, err := internal.WirePgConnPool(ctx)
		if err != nil {
			syncErr = fmt.Errorf("error creating db pool: %w", err)
			return
		}
		repo := repository.New(db)
		kw, err := kafka.New(cfg)
		if err != nil {
			syncErr = fmt.Errorf("error creating kafka topics: %w", err)
			return
		}
		ws := ws.New()
		uc := uc.NewUseCaseService(repo, kw, ws)
		rh := NewRestHandler(uc)

		// Setup server.
		e := echo.New()
		e.Validator = cv.ProvideValidator()

		tbs = &TestBedServer{
			server: e,
			handler: rh,
		}
	})
	return tbs, syncErr
}

func TestHandlerCreateMachineResourceUsageRT(t *testing.T) {
	tbs, err := ProvideTestBedServer(t)
	if err != nil {
		t.Fatalf("unable to setup test bed: %v", err)
	}

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

func TestHandlerGetAggMachineResourceUsageRT(t *testing.T) {
	tbs, err := ProvideTestBedServer(t)
	if err != nil {
		t.Fatalf("unable to setup test bed: %v", err)
	}

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