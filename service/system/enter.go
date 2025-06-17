package system

type ServiceGroup struct {
	JwtService
	ApiService
	UserService
	CasbinService
	AuthorityService
	OperationRecordService
	AuthorityApiService
}
