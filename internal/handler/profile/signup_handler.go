package profilehdlr

import (
	"net/http"
	"stori/internal/core/domain"
	"stori/internal/ports"
	"stori/pkg/response"
	"stori/pkg/validator"

	"github.com/gin-gonic/gin"
)

type ProfileHdlr struct {
	service ports.ProfileServicePort
}

// register params
type Parameters struct {
	// Full name
	Name string `json:"name" validate:"required"`

	// For authentication process
	Email string `json:"email" validate:"required,email"`
}

func ProvideProfileHandler(srv ports.ProfileServicePort) *ProfileHdlr {
	return &ProfileHdlr{
		service: srv,
	}
}

func (hdl *ProfileHdlr) SignUpHandler(ctx *gin.Context) {
	var input Parameters
	if err := ctx.BindJSON(&input); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Failure(err.Error()))
		return
	}

	errs := validator.ValidateStructure(input)
	if len(errs) > 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.FailureMappingErrors(errs))
		return
	}

	request := domain.ProfileDTO{
		Name:  input.Name,
		Email: input.Email,
	}

	profile, err := hdl.service.Create(&request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.Failure(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(profile))
}
