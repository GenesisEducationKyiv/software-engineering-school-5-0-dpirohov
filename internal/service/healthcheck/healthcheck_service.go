package healthcheck

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
	"weatherApi/internal/logger"

	_ "github.com/joho/godotenv/autoload"
)

type HealthCheckService interface {
	Health() map[string]string
	Close() error
}

type service struct {
	log   *logger.Logger
	sqlDB *sql.DB
}

func New(log *logger.Logger, db *sql.DB) HealthCheckService {
	return &service{
		log:   log,
		sqlDB: db,
	}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the DB to check connectivity
	err := s.sqlDB.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		s.log.Base().Error().Err(err).Msg("Database down!")
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	dbStats := s.sqlDB.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Warnings
	if dbStats.OpenConnections > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed. Consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime. Consider increasing it or revising usage patterns."
	}

	return stats
}

func (s *service) Close() error {
	s.log.Base().Info().Msg("Disconnected from database")
	return s.sqlDB.Close()
}
