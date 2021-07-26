package webhook_test

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/multierr"
)

func Test_Authorizer(t *testing.T) {
	db := pgtest.NewGormDB(t)

	eiFoo := cltest.MustInsertExternalInitiator(t, db)
	eiBar := cltest.MustInsertExternalInitiator(t, db)

	jobWithFooAndBarEI, webhookSpecWithFooAndBarEI := cltest.MustInsertWebhookSpec(t, db)
	jobWithBarEI, webhookSpecWithBarEI := cltest.MustInsertWebhookSpec(t, db)
	jobWithNoEI, _ := cltest.MustInsertWebhookSpec(t, db)

	require.NoError(t, multierr.Combine(
		db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, eiFoo.ID, webhookSpecWithFooAndBarEI.ID, `{"ei": "foo", "name": "webhookSpecWithFooAndBarEI"}`).Error,
		db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, eiBar.ID, webhookSpecWithFooAndBarEI.ID, `{"ei": "bar", "name": "webhookSpecWithFooAndBarEI"}`).Error,
		db.Exec(`INSERT INTO external_initiator_webhook_specs (external_initiator_id, webhook_spec_id, spec) VALUES (?,?,?)`, eiBar.ID, webhookSpecWithBarEI.ID, `{"ei": "bar", "name": "webhookSpecTwoEIs"}`).Error,
	))

	t.Run("no user no ei never authorizes", func(t *testing.T) {
		a := webhook.NewAuthorizer(db, nil, nil)

		can, err := a.CanRun(jobWithFooAndBarEI.ExternalJobID)
		require.NoError(t, err)
		assert.False(t, can)
		can, err = a.CanRun(jobWithNoEI.ExternalJobID)
		require.NoError(t, err)
		assert.False(t, can)
		can, err = a.CanRun(uuid.NewV4())
		require.NoError(t, err)
		assert.False(t, can)
	})

	t.Run("with user no ei always authorizes", func(t *testing.T) {
		a := webhook.NewAuthorizer(db, &models.User{}, nil)

		can, err := a.CanRun(jobWithFooAndBarEI.ExternalJobID)
		require.NoError(t, err)
		assert.True(t, can)
		can, err = a.CanRun(jobWithNoEI.ExternalJobID)
		require.NoError(t, err)
		assert.True(t, can)
		can, err = a.CanRun(uuid.NewV4())
		require.NoError(t, err)
		assert.True(t, can)
	})

	t.Run("no user with ei authorizes conditionally", func(t *testing.T) {
		a := webhook.NewAuthorizer(db, nil, &eiFoo)

		can, err := a.CanRun(jobWithFooAndBarEI.ExternalJobID)
		require.NoError(t, err)
		assert.True(t, can)
		can, err = a.CanRun(jobWithBarEI.ExternalJobID)
		require.NoError(t, err)
		assert.False(t, can)
		can, err = a.CanRun(jobWithNoEI.ExternalJobID)
		require.NoError(t, err)
		assert.False(t, can)
		can, err = a.CanRun(uuid.NewV4())
		require.NoError(t, err)
		assert.False(t, can)
	})
}
