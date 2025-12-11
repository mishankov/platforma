package auth

import (
	"github.com/platforma-dev/platforma/httpserver"
)

type Domain struct {
	Repository  *Repository
	Service     *Service
	HandleGroup *httpserver.HandlerGroup
	Middleware  httpserver.Middleware
}

func (d *Domain) GetRepository() any {
	return d.Repository
}

func New(db db, authStorage authStorage, sessionCookieName string, usernameValidator, passwordValidator func(string) error) *Domain {
	repository := NewRepository(db)
	service := NewService(repository, authStorage, sessionCookieName, usernameValidator, passwordValidator)

	authMiddleware := NewAuthenticationMiddleware(service)
	registerHandler := NewRegisterHandler(service)
	loginHandler := NewLoginHandler(service)
	logoutHandler := NewLogoutHandler(service)
	getUserHandler := NewGetHandler(service)
	changePasswordHandler := authMiddleware.Wrap(NewChangePasswordHandler(service))
	deleteHandler := authMiddleware.Wrap(NewDeleteHandler(service))

	authAPI := httpserver.NewHandlerGroup()
	authAPI.Handle("POST /register", registerHandler)
	authAPI.Handle("POST /login", loginHandler)
	authAPI.Handle("POST /logout", logoutHandler)
	authAPI.Handle("GET /me", getUserHandler)
	authAPI.Handle("POST /change-password", changePasswordHandler)
	authAPI.Handle("DELETE /me", deleteHandler)

	return &Domain{
		Repository:  repository,
		Service:     service,
		HandleGroup: authAPI,
		Middleware:  authMiddleware,
	}
}
