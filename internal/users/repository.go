package users

import (
	"context"

	"github.com/courage173/quiz-api/internal/models"
	"github.com/courage173/quiz-api/pkg/dbcontext"
	dbx "github.com/go-ozzo/ozzo-dbx"

	"github.com/courage173/quiz-api/pkg/log"
)

type Repository interface {
	Create(ctx context.Context, user models.User) error
	GetByID(ctx context.Context, id int) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	Query(ctx context.Context, offset, limit int) ([]models.User, error)
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id int) error
	EmailExist(ctx context.Context, email string) bool
}

type repository struct {
	db *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
    return repository{db, logger}
}

func (r repository) Create(ctx context.Context, user models.User) error {
    return r.db.With(ctx).Model(&user).Insert()
}

func (r repository) GetByID(ctx context.Context, id int) (models.User, error) {
	var user models.User
    err := r.db.With(ctx).Select().Model(id, &user)
	return user, err
}

func (r repository) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
     err := r.db.With(ctx).Select().Where(dbx.HashExp{"email": email}).One(&user)
    return user, err
}

func (r repository) Query(ctx context.Context, offset, limit int) ([]models.User, error) {
	var users []models.User

    err := r.db.With(ctx).Select().Offset(int64(offset)).Limit(int64(limit)).All(&users)
    return users, err
}

func (r repository) EmailExist(ctx context.Context, email string) bool {
	var count int64
    err := r.db.With(ctx).Select("COUNT(*)").Where(dbx.HashExp{"email": email}).One(&count)
    return count > 0 && err == nil
}

func (r repository) Update(ctx context.Context, user models.User) error {
	db := *r.db.With(ctx).Model(&user)
	return db.Update()
}

func (r repository) Delete(ctx context.Context, id int) error {
	user, err := r.GetByID(ctx, id)
	if err!= nil {
        return err
    }
	r.logger.Infof("Deleting user with ID: %d", id)
	return r.db.With(ctx).Model(&user).Delete()
}

