package helper

import "github.com/gofiber/fiber/v2"

func Message400(msg string) error {
	return fiber.NewError(fiber.StatusBadRequest, msg)
}

func Message401(msg string) error {
	return fiber.NewError(fiber.StatusUnauthorized, msg)
}

func Message403(msg string) error {
	return fiber.NewError(fiber.StatusForbidden, msg)
}

func Message404(msg string) error {
	return fiber.NewError(fiber.StatusNotFound, msg)
}

func Message500(msg string) error {
	return fiber.NewError(fiber.StatusInternalServerError, msg)
}
