package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"sledgehammer.echo-mesh.com/internal/api"
	"sledgehammer.echo-mesh.com/internal/clientbans"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	db, err := sqlx.Open("sqlite3", "./sledgehammer.db")
	if err != nil {
		sugar.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	clientBanStore := clientbans.NewStore(db)
	if err := clientBanStore.InitSchema(); err != nil {
		sugar.Fatalf("failed to init schema: %v", err)
	}

	apiHandler := &api.API{Store: clientBanStore}

	// --- Prometheus metrics ---
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of request latencies",
			Buckets: prometheus.DefBuckets, // default: 0.005s → 10s
		},
		[]string{"path", "method", "status"},
	)

	inFlightRequests := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_in_flight_requests",
			Help: "Current number of requests being processed",
		},
	)

	totalBans := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ban_submissions_total",
			Help: "Total number of ban submissions",
		},
	)

	prometheus.MustRegister(requestCounter, requestDuration, inFlightRequests, totalBans)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

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

	r.Get("/bans", apiHandler.ListBans)
	r.Post("/bans", func(w http.ResponseWriter, r *http.Request) {
		totalBans.Inc()
		apiHandler.AddClientBan(w, r)
	})

	// Prometheus metrics
	r.Handle("/metrics", promhttp.Handler())

	sugar.Info("Echo Sledgehammer server running on :8080")
	sugar.Fatal(http.ListenAndServe(":8080", r))
}
