package horizon

import (
	"github.com/rcrowley/go-metrics"
	"github.com/rs/cors"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

type Web struct {
	router *web.Mux

	requestTimer metrics.Timer
	failureMeter metrics.Meter
	successMeter metrics.Meter
}

func NewWeb(app *App) (*Web, error) {

	api := web.New()
	result := Web{
		router:       api,
		requestTimer: metrics.NewTimer(),
		failureMeter: metrics.NewMeter(),
		successMeter: metrics.NewMeter(),
	}

	app.metrics.Register("requests.total", result.requestTimer)
	app.metrics.Register("requests.succeeded", result.successMeter)
	app.metrics.Register("requests.failed", result.failureMeter)

	installMiddleware(api, app)
	installActions(api)

	return &result, nil
}

func installMiddleware(api *web.Mux, app *App) {
	api.Use(middleware.RequestID)
	api.Use(middleware.Logger)
	api.Use(middleware.Recoverer)
	api.Use(middleware.AutomaticOptions)
	api.Use(appMiddleware(app))
	api.Use(requestMetricsMiddleware)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	api.Use(c.Handler)
}

func installActions(api *web.Mux) {
	api.Get("/", rootAction)
	api.Get("/metrics", metricsAction)

	// ledger actions
	api.Get("/ledgers", ledgerIndexAction)
	api.Get("/ledgers/:id", ledgerShowAction)
	api.Get("/ledgers/:ledger_id/transactions", notImplementedAction)
	api.Get("/ledgers/:ledger_id/operations", notImplementedAction)
	api.Get("/ledgers/:ledger_id/effects", notImplementedAction)

	// account actions
	api.Get("/accounts", notImplementedAction)
	api.Get("/accounts/:id", accountShowAction)
	api.Get("/accounts/:account_id/transactions", notImplementedAction)
	api.Get("/accounts/:account_id/operations", notImplementedAction)
	api.Get("/accounts/:account_id/effects", notImplementedAction)

	// transaction actions
	api.Get("/transactions", notImplementedAction)
	api.Get("/transactions/:id", notImplementedAction)
	api.Get("/transactions/:tx_id/operations", notImplementedAction)
	api.Get("/transactions/:tx_id/effects", notImplementedAction)

	// transaction actions
	api.Get("/operations", notImplementedAction)
	api.Get("/operations/:id", notImplementedAction)
	api.Get("/operations/:tx_id/effects", notImplementedAction)

	api.Get("/stream", streamAction)

	api.NotFound(notFoundAction)
}
