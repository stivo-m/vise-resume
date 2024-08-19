package middleware

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stivo-m/vise-resume/internal/core/utils"
)

var validate *validator.Validate

type ValidationErrorDto struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

func init() {
	validate = validator.New()
}

func ValidationMiddleware(dto interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Create a new instance of the DTO
		dtoInstance := reflect.New(reflect.TypeOf(dto).Elem()).Interface()
		if err := c.BodyParser(dtoInstance); err != nil {
			res := utils.FormatApiResponse(
				"The request body is invalid",
				err.Error(),
			)
			return c.Status(fiber.StatusBadRequest).JSON(res)
		}

		if err := validate.Struct(dtoInstance); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				res := utils.FormatApiResponse(
					"Validation error",
					err.Error(),
				)
				return c.Status(fiber.StatusInternalServerError).JSON(res)

			}

			var errorList []ValidationErrorDto
			for _, err := range err.(validator.ValidationErrors) {
				field := utils.GetJSONFieldName(dtoInstance, err.StructField())
				item := ValidationErrorDto{
					Field:   field,
					Rule:    err.Tag(),
					Message: utils.GetValidationMessage(err, field),
				}
				errorList = append(errorList, item)
			}

			if len(errorList) > 0 {
				data := utils.FormatApiResponse(
					"one or more of the required fields are invalid or missing",
					errorList,
				)
				return c.Status(fiber.StatusUnprocessableEntity).JSON(data)
			}
		}

		return c.Next()
	}
}
