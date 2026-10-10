package response

import "github.com/gofiber/fiber/v2"

// ErrorVm represents standard RFC7807/Spring style error view model
type ErrorVm struct {
	StatusCode string `json:"statusCode"`
	Title      string `json:"title"`
	Detail     string `json:"detail"`
}

func NotFound(c *fiber.Ctx, detail string) error {
	return c.Status(fiber.StatusNotFound).JSON(ErrorVm{
		StatusCode: "404",
		Title:      "Not found",
		Detail:     detail,
	})
}

func BadRequest(c *fiber.Ctx, detail string) error {
	return c.Status(fiber.StatusBadRequest).JSON(ErrorVm{
		StatusCode: "400",
		Title:      "Bad request",
		Detail:     detail,
	})
}

func InternalError(c *fiber.Ctx, detail string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(ErrorVm{
		StatusCode: "500",
		Title:      "Internal server error",
		Detail:     detail,
	})
}
