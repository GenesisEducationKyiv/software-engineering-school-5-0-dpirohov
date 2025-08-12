package tests

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func Test_WeatherService_ShouldOnlyUseCacheRepository(t *testing.T) {
	archtest.Package(t, "weatherApi/internal/service/weather").
		ShouldNotDependOn(
			"weatherApi/internal/repository/subscription",
			"weatherApi/internal/repository/user",
		)
}

func Test_SubscriptionService_ShouldOnlyUseSQLRepositories(t *testing.T) {
	archtest.Package(t, "weatherApi/internal/service/subscription").
		ShouldNotDependOn(
			"weatherApi/internal/repository/weather",
		)
}

func Test_Repositories_ShouldNotDependOnServiceOrExternalLayers(t *testing.T) {
	archtest.Package(t, "weatherApi/internal/repository/...").
		ShouldNotDependOn(
			"weatherApi/internal/service/...",
			"weatherApi/internal/server/...",
			"weatherApi/internal/provider/...",
			"weatherApi/internal/broker/...",
			"weatherApi/internal/service/...",
		)
}

func Test_Config_ShouldNotDependOnAnyInternalLogic(t *testing.T) {
	archtest.Package(t, "weatherApi/internal/config/...").
		ShouldNotDependOn(
			"weatherApi/internal/service/...",
			"weatherApi/internal/server/...",
			"weatherApi/internal/provider/...",
			"weatherApi/internal/broker/...",
			"weatherApi/internal/repository/...",
			"weatherApi/internal/tests/...",
			"weatherApi/internal/service/...",
		)
}

func Test_EventPublisher_ShouldNotDependOnBusinessLogic(t *testing.T) {
	archtest.Package(t, "weatherApi/internal/broker/...").
		ShouldNotDependOn(
			"weatherApi/internal/service/...",
			"weatherApi/internal/server/...",
			"weatherApi/internal/provider/...",
			"weatherApi/internal/repository/...",
			"weatherApi/internal/tests/...",
			"weatherApi/internal/service/...",
		)
}

func Test_Metrics_ShouldNotDependOnInternalLayer(t *testing.T) {
	archtest.Package(t, "weatherApi/internal/metrics/...").
		ShouldNotDependOn(
			"weatherApi/internal/service/...",
			"weatherApi/internal/server/...",
			"weatherApi/internal/provider/...",
			"weatherApi/internal/repository/...",
			"weatherApi/internal/tests/...",
			"weatherApi/internal/service/...",
		)
}
