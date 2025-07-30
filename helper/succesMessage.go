package helper

import "github.com/gofiber/fiber/v2"

func Message200(c *fiber.Ctx, data interface{}, msg string) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": msg,
		"data":    data,
	})
}

func Message201(c *fiber.Ctx, data interface{}, msg string) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": msg,
		"data":    data,
	})
}
