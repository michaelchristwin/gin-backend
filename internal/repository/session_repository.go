package repository

type SessionRepository interface {
	Create()
	Get()
	Delete()
	DeleteUserSessions()
	DeleteExpiredSessions()
	GetSessionWithUser()
}
