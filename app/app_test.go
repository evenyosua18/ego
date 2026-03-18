package app

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/evenyosua18/ego/tracer"
	"github.com/evenyosua18/ego/sqldb"
	"github.com/evenyosua18/ego/cache/redis_adapter"
	"github.com/evenyosua18/ego/http"

	"github.com/evenyosua18/ego/config"
	"github.com/stretchr/testify/assert"
)

func TestGetApp(t *testing.T) {
	appInstance := GetApp()
	assert.NotNil(t, appInstance)
}

func TestRunRest_Success_All(t *testing.T) {
	// Mock external funcs to return success
	oldLoadCodes := loadCodesFunc
	oldRunSentry := runSentryFunc
	oldDbConnect := dbConnectFunc
	oldRedisNew := redisNewFunc
	oldAuthManage := authManageFunc
	oldOsExit := osExit

type mockDB struct { sqldb.ISqlDB }

	loadCodesFunc = func(path string) error { return nil }
	runSentryFunc = func(cfg tracer.Config) (func(string), error) {
		return func(string) {}, nil
	}
	dbConnectFunc = func(driver, dsn string, cfg *sqldb.Config) (sqldb.ISqlDB, error) {
		return &mockDB{}, nil
	}
	redisNewFunc = func(cfg redis_adapter.RedisConfig) (*redis_adapter.RedisAdapter, error) {
		return &redis_adapter.RedisAdapter{}, nil
	}
	authManageFunc = func(ctx context.Context, baseUrl, clientId, clientSecret string) error {
		return nil
	}
	oldNewRouter := newRouterFunc
	newRouterFunc = func(cfg http.RouteConfig) http.IHttpRouter {
		return &mockRouterSuccess{}
	}

	var exitCode int
	osExit = func(code int) {
		exitCode = code
	}

	var capturedQuit chan<- os.Signal
	oldSignalNotify := signalNotify
	signalNotify = func(c chan<- os.Signal, sig ...os.Signal) {
		capturedQuit = c
	}

	defer func() {
		loadCodesFunc = oldLoadCodes
		runSentryFunc = oldRunSentry
		dbConnectFunc = oldDbConnect
		redisNewFunc = oldRedisNew
		authManageFunc = oldAuthManage
		newRouterFunc = oldNewRouter
		osExit = oldOsExit
		signalNotify = oldSignalNotify
	}()

	config.SetTestConfig(map[string]any{
		RouterShowRegistered:   false,
		RouterPort:             ":0",
		CustomCodeFilePath:     "dummy.json",
		TracerDSN:              "dummy-dsn",
		DatabaseName:           "dummy-db",
		CacheRedisAddr:         "dummy-redis:6379",
		AuthSvcClientId:        "dummy-client",
		AuthSvcClientSecret:    "dummy-secret",
	})

	go func() {
		time.Sleep(100 * time.Millisecond)
		if capturedQuit != nil {
			capturedQuit <- os.Interrupt
		}
	}()

	appInst := GetApp()
	appInst.RunRest()

	assert.Equal(t, 0, exitCode)
}

func TestRunRest_ShutdownTimeout(t *testing.T) {
	config.SetTestConfig(map[string]any{
		RouterShowRegistered:   false,
		RouterPort:             ":0",
		ServiceShutdownTimeout: "1ns",
	})

	var exitCode int
	oldOsExit := osExit
	osExit = func(code int) {
		exitCode = code
	}

	oldNewRouter := newRouterFunc
	newRouterFunc = func(cfg http.RouteConfig) http.IHttpRouter {
		return &mockRouterShutdownError{}
	}

	var capturedQuit chan<- os.Signal
	oldSignalNotify := signalNotify
	signalNotify = func(c chan<- os.Signal, sig ...os.Signal) {
		capturedQuit = c
	}
	defer func() {
		osExit = oldOsExit
		signalNotify = oldSignalNotify
		newRouterFunc = oldNewRouter
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		if capturedQuit != nil {
			capturedQuit <- os.Interrupt
		}
	}()

	appInst := GetApp()
	appInst.RunRest()

	assert.Equal(t, 0, exitCode)
}

func TestRunRest_CodePanic(t *testing.T) {
	config.SetTestConfig(map[string]any{
		CustomCodeFilePath: "invalid.json",
	})

	oldLoadCodes := loadCodesFunc
	loadCodesFunc = func(path string) error {
		return fmt.Errorf("load codes error")
	}
	defer func() { loadCodesFunc = oldLoadCodes }()

	appInst := GetApp()
	assert.Panics(t, func() {
		appInst.RunRest()
	})
}

func TestRunRest_TracerPanic(t *testing.T) {
	config.SetTestConfig(map[string]any{
		TracerDSN: "invalid-dsn",
	})
	oldRunSentry := runSentryFunc
	runSentryFunc = func(cfg tracer.Config) (func(string), error) {
		return nil, fmt.Errorf("tracer error")
	}
	defer func() { runSentryFunc = oldRunSentry }()

	appInst := GetApp()
	assert.Panics(t, func() {
		appInst.RunRest()
	})
}

func TestRunRest_DBPanic(t *testing.T) {
	config.SetTestConfig(map[string]any{
		DatabaseName:   "test",
		DatabaseDriver: "invalid",
	})
	oldDbConnect := dbConnectFunc
	dbConnectFunc = func(driver, dsn string, cfg *sqldb.Config) (sqldb.ISqlDB, error) {
		return nil, fmt.Errorf("db error")
	}
	defer func() { dbConnectFunc = oldDbConnect }()

	appInst := GetApp()
	assert.Panics(t, func() {
		appInst.RunRest()
	})
}

func TestRunRest_RedisPanic(t *testing.T) {
	config.SetTestConfig(map[string]any{
		CacheRedisAddr: "invalid-addr:99999",
	})
	oldRedisNew := redisNewFunc
	redisNewFunc = func(cfg redis_adapter.RedisConfig) (*redis_adapter.RedisAdapter, error) {
		return nil, fmt.Errorf("redis error")
	}
	defer func() { redisNewFunc = oldRedisNew }()

	appInst := GetApp()
	assert.Panics(t, func() {
		appInst.RunRest()
	})
}

func TestRunRest_AuthSvcPanic(t *testing.T) {
	config.SetTestConfig(map[string]any{
		AuthSvcClientId:     "client",
		AuthSvcClientSecret: "secret",
		AuthSvcBaseUrl:      "http://invalid-url:-1",
	})
	oldAuthManage := authManageFunc
	authManageFunc = func(ctx context.Context, baseUrl, clientId, clientSecret string) error {
		return fmt.Errorf("auth error")
	}
	defer func() { authManageFunc = oldAuthManage }()

	appInst := GetApp()
	assert.Panics(t, func() {
		appInst.RunRest()
	})
}

// Mock router tests listen panic
type mockRouter struct {
	http.IHttpRouter
}

func (m *mockRouter) Listen(port string) error {
	return fmt.Errorf("listen error")
}

func (m *mockRouter) ShutdownWithContext(ctx context.Context) error {
	return nil
}

func (m *mockRouter) ActiveConnections() int {
	return 0
}

type mockRouterSuccess struct {
	http.IHttpRouter
}

func (m *mockRouterSuccess) Listen(port string) error {
	return nil
}

func (m *mockRouterSuccess) ShutdownWithContext(ctx context.Context) error {
	return nil
}

func (m *mockRouterSuccess) ActiveConnections() int {
	return 0
}

type mockRouterShutdownError struct {
	http.IHttpRouter
}

func (m *mockRouterShutdownError) Listen(port string) error {
	return nil
}

func (m *mockRouterShutdownError) ShutdownWithContext(ctx context.Context) error {
	return fmt.Errorf("shutdown error")
}

func (m *mockRouterShutdownError) ActiveConnections() int {
	return 0
}

func TestRunRest_ListenPanic(t *testing.T) {
	config.SetTestConfig(map[string]any{})

	oldOsExit := osExit
	osExit = func(code int) {}
	defer func() { osExit = oldOsExit }()

	oldPanic := panicFunc
	defer func() { panicFunc = oldPanic }()

	oldNewRouter := newRouterFunc
	newRouterFunc = func(cfg http.RouteConfig) http.IHttpRouter {
		return &mockRouter{}
	}
	defer func() { newRouterFunc = oldNewRouter }()

	var capturedQuit chan<- os.Signal
	oldSignalNotify := signalNotify
	signalNotify = func(c chan<- os.Signal, sig ...os.Signal) {
		capturedQuit = c
	}
	defer func() { signalNotify = oldSignalNotify }()

	// channel to wait for panic
	panicChan := make(chan bool)
	panicFunc = func(v any) {
		panicChan <- true
	}

	go func() {
		select {
		case <-panicChan:
			// Panicked as expected
			if capturedQuit != nil {
				capturedQuit <- os.Interrupt
			}
		case <-time.After(1 * time.Second):
			panic("expected panic from router listen")
		}
	}()

	appInst := GetApp()
	appInst.RunRest()
}
