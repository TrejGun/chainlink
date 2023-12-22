package loadfunctions

import (
	"testing"
	"time"

	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/require"

	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
)

func TestFunctionsLoad(t *testing.T) {
	generalConfig, err := tc.GetConfig(t.Name(), tc.Load, tc.Functions)
	require.NoError(t, err, "failed to get config")

	ft, err := SetupLocalLoadTestEnv(&generalConfig)
	require.NoError(t, err)
	ft.EVMClient.ParallelTransactions(false)

	labels := map[string]string{
		"branch": "functions_healthcheck",
		"commit": "functions_healthcheck",
	}

	MonitorLoadStats(t, ft, labels, &generalConfig)

	t.Run("mumbai functions soak test http", func(t *testing.T) {
		config, err := tc.GetConfig(t.Name(), tc.Soak, tc.Functions)
		require.NoError(t, err, "failed to get config")
		cfg := config.Functions
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:                     t,
				LoadType:              wasp.RPS,
				GenName:               "functions_soak_gen",
				RateLimitUnitDuration: 5 * time.Second,
				CallTimeout:           3 * time.Minute,
				Schedule: wasp.Plain(
					*cfg.Performance.RPS,
					cfg.Performance.Duration.Duration(),
				),
				Gun: NewSingleFunctionCallGun(
					ft,
					ModeHTTPPayload,
					*cfg.Performance.RequestsPerCall,
					*cfg.Common.FunctionsCallPayloadHTTP,
					*cfg.Common.SecretsSlotID,
					*cfg.Common.SecretsVersionID,
					[]string{},
					*cfg.Common.SubscriptionID,
					StringToByte32(*cfg.Common.DONID),
				),
				Labels:     labels,
				LokiConfig: tc.LokiConfigFromToml(&config),
			})).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("mumbai functions stress test http", func(t *testing.T) {
		config, err := tc.GetConfig(t.Name(), tc.Stress, tc.Functions)
		require.NoError(t, err, "failed to get config")
		cfg := config.Functions
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:                     t,
				LoadType:              wasp.RPS,
				GenName:               "functions_stress_gen",
				RateLimitUnitDuration: 5 * time.Second,
				CallTimeout:           3 * time.Minute,
				Schedule: wasp.Plain(
					*cfg.Performance.RPS,
					cfg.Performance.Duration.Duration(),
				),
				Gun: NewSingleFunctionCallGun(
					ft,
					ModeHTTPPayload,
					*cfg.Performance.RequestsPerCall,
					*cfg.Common.FunctionsCallPayloadHTTP,
					*cfg.Common.SecretsSlotID,
					*cfg.Common.SecretsVersionID,
					[]string{},
					*cfg.Common.SubscriptionID,
					StringToByte32(*cfg.Common.DONID),
				),
				Labels:     labels,
				LokiConfig: tc.LokiConfigFromToml(&config),
			})).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("mumbai functions soak test only secrets", func(t *testing.T) {
		config, err := tc.GetConfig(t.Name(), "secrets_soak", tc.Functions)
		require.NoError(t, err, "failed to get config")
		cfg := config.Functions
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:                     t,
				LoadType:              wasp.RPS,
				GenName:               "functions_soak_gen",
				RateLimitUnitDuration: 5 * time.Second,
				CallTimeout:           3 * time.Minute,
				Schedule: wasp.Plain(
					*cfg.Performance.RPS,
					cfg.Performance.Duration.Duration(),
				),
				Gun: NewSingleFunctionCallGun(
					ft,
					ModeSecretsOnlyPayload,
					*cfg.Performance.RequestsPerCall,
					*cfg.Common.FunctionsCallPayloadWithSecrets,
					*cfg.Common.SecretsSlotID,
					*cfg.Common.SecretsVersionID,
					[]string{},
					*cfg.Common.SubscriptionID,
					StringToByte32(*cfg.Common.DONID),
				),
				Labels:     labels,
				LokiConfig: tc.LokiConfigFromToml(&config),
			})).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("mumbai functions stress test only secrets", func(t *testing.T) {
		config, err := tc.GetConfig(t.Name(), "secrets_stress", tc.Functions)
		require.NoError(t, err, "failed to get config")
		cfg := config.Functions
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:                     t,
				LoadType:              wasp.RPS,
				GenName:               "functions_stress_gen",
				RateLimitUnitDuration: 5 * time.Second,
				CallTimeout:           3 * time.Minute,
				Schedule: wasp.Plain(
					*cfg.Performance.RPS,
					cfg.Performance.Duration.Duration(),
				),
				Gun: NewSingleFunctionCallGun(
					ft,
					ModeSecretsOnlyPayload,
					*cfg.Performance.RequestsPerCall,
					*cfg.Common.FunctionsCallPayloadWithSecrets,
					*cfg.Common.SecretsSlotID,
					*cfg.Common.SecretsVersionID,
					[]string{},
					*cfg.Common.SubscriptionID,
					StringToByte32(*cfg.Common.DONID),
				),
				Labels:     labels,
				LokiConfig: tc.LokiConfigFromToml(&config),
			})).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("mumbai functions soak test real", func(t *testing.T) {
		config, err := tc.GetConfig(t.Name(), "real_soak", tc.Functions)
		require.NoError(t, err, "failed to get config")
		cfg := config.Functions
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:                     t,
				LoadType:              wasp.RPS,
				GenName:               "functions_soak_gen",
				RateLimitUnitDuration: 5 * time.Second,
				CallTimeout:           3 * time.Minute,
				Schedule: wasp.Plain(
					*cfg.Performance.RPS,
					cfg.Performance.Duration.Duration(),
				),
				Gun: NewSingleFunctionCallGun(
					ft,
					ModeReal,
					*cfg.Performance.RequestsPerCall,
					*cfg.Common.FunctionsCallPayloadReal,
					*cfg.Common.SecretsSlotID,
					*cfg.Common.SecretsVersionID,
					[]string{"1", "2", "3", "4"},
					*cfg.Common.SubscriptionID,
					StringToByte32(*cfg.Common.DONID),
				),
				Labels:     labels,
				LokiConfig: tc.LokiConfigFromToml(&config),
			})).
			Run(true)
		require.NoError(t, err)
	})

	t.Run("mumbai functions stress test real", func(t *testing.T) {
		config, err := tc.GetConfig(t.Name(), "real_stress", tc.Functions)
		require.NoError(t, err, "failed to get config")
		cfg := config.Functions
		_, err = wasp.NewProfile().
			Add(wasp.NewGenerator(&wasp.Config{
				T:                     t,
				LoadType:              wasp.RPS,
				GenName:               "functions_stress_gen",
				RateLimitUnitDuration: 5 * time.Second,
				CallTimeout:           3 * time.Minute,
				Schedule: wasp.Plain(
					*cfg.Performance.RPS,
					cfg.Performance.Duration.Duration(),
				),
				Gun: NewSingleFunctionCallGun(
					ft,
					ModeReal,
					*cfg.Performance.RequestsPerCall,
					*cfg.Common.FunctionsCallPayloadReal,
					*cfg.Common.SecretsSlotID,
					*cfg.Common.SecretsVersionID,
					[]string{"1", "2", "3", "4"},
					*cfg.Common.SubscriptionID,
					StringToByte32(*cfg.Common.DONID),
				),
				Labels:     labels,
				LokiConfig: tc.LokiConfigFromToml(&config),
			})).
			Run(true)
		require.NoError(t, err)
	})
}
