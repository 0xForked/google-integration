package user

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func NewUserModuleProvider(
	rg *gin.RouterGroup,
	db *sql.DB,
) {
	repo := newSQLRepository(db)
	svc := newUserService(repo)
	newUserHandler(svc, rg)
	newBookingHandler(svc, rg)
}
