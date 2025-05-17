package repository

import (
	"github.com/andrianprasetya/eventHub/internal/user/dto/request"
	"github.com/andrianprasetya/eventHub/internal/user/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	CreateWithTx(tx *gorm.DB, user *model.User) error
	GetByEmail(email string) (*model.User, error)
	GetByID(id string) (*model.User, error)
	GetAll(query request.UserPaginateParams, tenantID *string) ([]*model.User, int64, error)
	Update(user *model.User) error
	Delete(id string) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *userRepository) CreateWithTx(tx *gorm.DB, user *model.User) error {
	return tx.Create(user).Error
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	var count int64

	//return kalo data tidak ada 0
	r.DB.Model(&model.User{}).Count(&count)
	if count == 0 {
		return nil, nil
	}

	if err := r.DB.Preload("Tenant").Preload("Role").First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetAll(query request.UserPaginateParams, tenantID *string) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	db := r.DB.Model(&model.User{}).Scopes(FilterUserQuery(query, tenantID))

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize

	if err := db.Preload("Role").Limit(query.PageSize).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	if err := r.DB.Preload("Tenant").Preload("Role").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.DB.Save(user).Error
}

func (r *userRepository) Delete(id string) error {
	return r.DB.Where("id = ?", id).Delete(&model.User{}).Error
}

func FilterUserQuery(query request.UserPaginateParams, tenantID *string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		//tenant params if super admin want filter the list
		if query.TenantID != nil {
			db = db.Where("tenant_id = ?", *query.TenantID)
		}
		if query.RoleID != nil {
			db = db.Where("role_id = ?", *query.RoleID)
		}
		if query.Name != nil {
			db = db.Where("name ILIKE ?", "%"+*query.Name+"%")
		}
		if query.Email != nil {
			db = db.Where("email ILIKE ?", "%"+*query.Email+"%")
		}
		if query.IsActive != nil {
			db = db.Where("is_active = ?", *query.IsActive)
		}
		//tenant id for filter if logged is tenant and select the users of self tenant
		if tenantID != nil {
			db = db.Where("tenant_id = ?", *tenantID)
		}
		return db
	}
}
