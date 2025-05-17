package request

type RolePaginateParams struct {
	Page     int     `query:"page"`
	PageSize int     `query:"pageSize"`
	Name     *string `query:"name"`
}
