package infradb

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"sandbox-api-gin/internal/domain/apperror"
	"sandbox-api-gin/internal/domain/model"
	"sandbox-api-gin/internal/domain/repository"
)

// bit(1)カラムはMySQLドライバで[]byteとして返るため、SQLで+0キャストしuint8でスキャンする
type sandboxUserRecord struct {
	ID           int64      `db:"id"`
	UserID       string     `db:"user_id"`
	EmailAddress string     `db:"email_address"`
	NickName     string     `db:"nick_name"`
	Approved     uint8      `db:"approved"`
	ApprovedAt   *time.Time `db:"approved_at"`
	Admin        uint8      `db:"admin"`
	Blocked      uint8      `db:"blocked"`
	Deleted      uint8      `db:"deleted"`
	CreatedAt    time.Time  `db:"created_at"`
	CreatedBy    string     `db:"created_by"`
	UpdatedAt    time.Time  `db:"updated_at"`
	UpdatedBy    string     `db:"updated_by"`
}

type MySQLUserRepository struct {
	db *sqlx.DB
}

func NewMySQLUserRepository(db *sqlx.DB) repository.UserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) Login(ctx context.Context, userID, email string) (*model.User, error) {
	query := `
		SELECT id, user_id, email_address, nick_name,
		       (approved+0) AS approved, approved_at,
		       (admin+0) AS admin, (blocked+0) AS blocked, (deleted+0) AS deleted,
		       created_at, created_by, updated_at, updated_by
		FROM sandbox_user
		WHERE user_id = ? AND email_address = ? AND deleted = 0`

	var rec sandboxUserRecord
	if err := r.db.GetContext(ctx, &rec, query, userID, email); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return toDomain(&rec), nil
}

func (r *MySQLUserRepository) FindByUserID(ctx context.Context, userID string) (*model.User, error) {
	query := `
		SELECT id, user_id, email_address, nick_name,
		       (approved+0) AS approved, approved_at,
		       (admin+0) AS admin, (blocked+0) AS blocked, (deleted+0) AS deleted,
		       created_at, created_by, updated_at, updated_by
		FROM sandbox_user
		WHERE user_id = ? AND deleted = 0`

	var rec sandboxUserRecord
	if err := r.db.GetContext(ctx, &rec, query, userID); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return toDomain(&rec), nil
}

func (r *MySQLUserRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	query := `SELECT COUNT(*) FROM sandbox_user WHERE user_id = ? AND deleted = 0`
	var count int
	if err := r.db.GetContext(ctx, &count, query, userID); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *MySQLUserRepository) InsertUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO sandbox_user
		    (user_id, email_address, nick_name, approved, approved_at, admin, blocked, deleted,
		     created_at, created_by, updated_at, updated_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query,
		user.UserID, user.EmailAddress, user.NickName,
		user.Approved, user.ApprovedAt, user.Admin, user.Blocked, user.Deleted,
		user.CreatedAt, user.CreatedBy, user.UpdatedAt, user.UpdatedBy,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewInsertError("ユーザ新規登録")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (r *MySQLUserRepository) UpdateNickName(ctx context.Context, user *model.User) error {
	query := `UPDATE sandbox_user SET nick_name = ?, updated_at = ?, updated_by = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, user.NickName, user.UpdatedAt, user.UpdatedBy, user.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewUpdateError("ユーザ情報更新")
	}
	return nil
}

func (r *MySQLUserRepository) SearchCount(ctx context.Context, emailAddress string, approved *bool) (int, error) {
	where, args := buildUserSearchWhere(emailAddress, approved)
	query := `SELECT COUNT(*) FROM sandbox_user WHERE deleted = 0` + where
	var count int
	if err := r.db.GetContext(ctx, &count, query, args...); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *MySQLUserRepository) Search(ctx context.Context, emailAddress string, approved *bool, page, size int) ([]*model.User, error) {
	where, args := buildUserSearchWhere(emailAddress, approved)
	query := `
		SELECT id, user_id, email_address, nick_name,
		       (approved+0) AS approved, approved_at,
		       (admin+0) AS admin, (blocked+0) AS blocked, (deleted+0) AS deleted,
		       created_at, created_by, updated_at, updated_by
		FROM sandbox_user
		WHERE deleted = 0` + where + `
		ORDER BY id
		LIMIT ? OFFSET ?`
	args = append(args, size, (page-1)*size)

	var recs []sandboxUserRecord
	if err := r.db.SelectContext(ctx, &recs, query, args...); err != nil {
		return nil, err
	}
	users := make([]*model.User, len(recs))
	for i := range recs {
		users[i] = toDomain(&recs[i])
	}
	return users, nil
}

func buildUserSearchWhere(emailAddress string, approved *bool) (string, []any) {
	var sb strings.Builder
	var args []any
	if emailAddress != "" {
		sb.WriteString(" AND email_address LIKE ?")
		args = append(args, "%"+emailAddress+"%")
	}
	if approved != nil {
		sb.WriteString(" AND (approved+0) = ?")
		if *approved {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}
	return sb.String(), args
}

func (r *MySQLUserRepository) UpdateApproved(ctx context.Context, user *model.User) error {
	query := `UPDATE sandbox_user SET approved = ?, approved_at = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, user.Approved, user.ApprovedAt, user.UpdatedAt, user.UpdatedBy, user.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewUpdateError("ユーザ承認")
	}
	return nil
}

func (r *MySQLUserRepository) UpdateBlocked(ctx context.Context, user *model.User) error {
	query := `UPDATE sandbox_user SET blocked = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, user.Blocked, user.UpdatedAt, user.UpdatedBy, user.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewUpdateError("Block設定")
	}
	return nil
}

func (r *MySQLUserRepository) UpdateAdmin(ctx context.Context, user *model.User) error {
	query := `UPDATE sandbox_user SET admin = ?, updated_at = ?, updated_by = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, user.Admin, user.UpdatedAt, user.UpdatedBy, user.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return apperror.NewUpdateError("管理者権限設定")
	}
	return nil
}

func toDomain(rec *sandboxUserRecord) *model.User {
	return &model.User{
		ID:           rec.ID,
		UserID:       rec.UserID,
		EmailAddress: rec.EmailAddress,
		NickName:     rec.NickName,
		Approved:     rec.Approved != 0,
		ApprovedAt:   rec.ApprovedAt,
		Admin:        rec.Admin != 0,
		Blocked:      rec.Blocked != 0,
		Deleted:      rec.Deleted != 0,
		CreatedAt:    rec.CreatedAt,
		CreatedBy:    rec.CreatedBy,
		UpdatedAt:    rec.UpdatedAt,
		UpdatedBy:    rec.UpdatedBy,
	}
}
