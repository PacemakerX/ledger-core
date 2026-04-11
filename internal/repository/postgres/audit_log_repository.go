package postgres

import (
	"context"

	"github.com/PacemakerX/ledger-core/internal/models"
	"github.com/PacemakerX/ledger-core/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type auditLogRepository struct {
	pool *pgxpool.Pool
}

func NewAuditLogRepository(pool *pgxpool.Pool) repository.AuditLogRepository {
	return &auditLogRepository{pool: pool}
}

func (r *auditLogRepository) Create(ctx context.Context, log *models.AuditLog) error {

	query := `INSERT INTO audit_logs(entity_type, entity_id, action, actor_id, actor_type, old_value, new_value, ip_address)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, entity_type, entity_id, action, actor_id, actor_type, old_value, new_value, ip_address, created_at`

	return r.pool.QueryRow(ctx, query,
		log.EntityType,
		log.EntityID,
		log.Action,
		log.ActorID,
		log.ActorType,
		log.OldValue,
		log.NewValue,
		log.IPAddress,
	).Scan(
		&log.ID,
		&log.EntityType,
		&log.EntityID,
		&log.Action,
		&log.ActorID,
		&log.ActorType,
		&log.OldValue,
		&log.NewValue,
		&log.IPAddress,
		&log.CreatedAt,
	)
}
