package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3" // <- required for side-effect registration
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"sledgehammer.echo-mesh.com/internal/api"
	"sledgehammer.echo-mesh.com/internal/database"
	middleware2 "sledgehammer.echo-mesh.com/internal/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := zap.NewProductionConfig()
	if os.Getenv("DEBUG") == "true" {
		cfg.Level.SetLevel(zap.DebugLevel)
	} else {
		cfg.Level.SetLevel(zap.InfoLevel)
	}
	logger, _ := cfg.Build()
	defer logger.Sync()
	sugar := logger.Sugar()

	db, err := sqlx.Open("sqlite3", os.Getenv("DATABASE_PATH"))
	if err != nil {
		sugar.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	clientBanStore := database.NewStore(db)
	if err := clientBanStore.InitSchema(); err != nil {
		sugar.Fatalf("failed to init schema: %v", err)
	}

	apiHandler := &api.API{Store: clientBanStore, Logger: sugar}

	// --- Prometheus metrics ---
	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
		[]string{"path", "method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Histogram of request latencies",
		Buckets: prometheus.DefBuckets, // default: 0.005s → 10s
	},
		[]string{"path", "method", "status"},
	)

	inFlightRequests := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "http_in_flight_requests",
		Help: "Current number of requests being processed",
	})

	totalPlayerBans := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "player_ban_submissions_total",
		Help: "Total number of player ban submissions",
	})

	totalPlayerBansLookups := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "player_ban_lookups_total",
		Help: "Total number of player ban lookups",
	})

	totalFileBans := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "File_ban_submissions_total",
		Help: "Total number of file ban submissions",
	})

	totalFileBansLookups := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "file_ban_lookups_total",
		Help: "Total number of file ban lookups",
	})

	prometheus.MustRegister(requestCounter, requestDuration, inFlightRequests, totalPlayerBans, totalFileBans, totalPlayerBansLookups, totalFileBansLookups)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware2.Ratelimit) // Add all endpoints with this middleware

	// Custom logging + metrics middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			inFlightRequests.Inc()
			start := time.Now()

			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			status := ww.Status()

			sugar.Infow("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", status,
				"duration", duration.Seconds(),
				"remote", r.RemoteAddr,
			)

			labels := []string{r.URL.Path, r.Method, http.StatusText(status)}
			requestCounter.WithLabelValues(labels...).Inc()
			requestDuration.WithLabelValues(labels...).Observe(duration.Seconds())
			inFlightRequests.Dec()
		})
	})

	r.Post("/login", apiHandler.Login)
	r.Group(func(authProtected chi.Router) {
		authProtected.Use(middleware2.AuthMiddleware())

		// Get
		authProtected.Get("/management/ban-history", apiHandler.GetBanHistory)
		authProtected.Get("/management/bans", apiHandler.GetCurrentActiveBans)
		authProtected.Get("/management/waiting-bans", apiHandler.GetWaitingForApprovalBans)

		// Post
		authProtected.Post("/management/approve-ban", apiHandler.ApproveBan)
		authProtected.Post("/management/reject-ban", apiHandler.RejectBan)

		// Delete
		authProtected.Delete("/management/ban", apiHandler.DeleteBan)
	})

	r.Get("/bans/{world}/player/{id}", apiHandler.FetchPlayerBanStatus)
	r.Get("/bans/{world}/player/{id}/info", func(w http.ResponseWriter, r *http.Request) {
		totalPlayerBansLookups.Inc()
		apiHandler.FetchPlayerBanInfo(w, r)
	})

	r.Post("/bans/{world}/player/{id}", func(w http.ResponseWriter, r *http.Request) {
		totalPlayerBans.Inc()
		apiHandler.RequestClientBan(w, r)
	})

	r.Get("/bans/file/{hash}", apiHandler.FetchFileBanstatus)
	r.Get("/bans/file/{hash}/info", func(w http.ResponseWriter, r *http.Request) {
		totalFileBansLookups.Inc()
		apiHandler.FetchFileBanInfo(w, r)
	})
	r.Post("/bans/file/{hash}", func(w http.ResponseWriter, r *http.Request) {
		totalFileBans.Inc()
		apiHandler.RequestFileBan(w, r)
	})

	// Prometheus metrics
	r.Handle("/metrics", promhttp.Handler())

	sugar.Info("Echo Sledgehammer server running on :8080")
	sugar.Fatal(http.ListenAndServe(":8080", r))
}
