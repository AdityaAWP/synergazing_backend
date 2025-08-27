package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"synergazing.com/synergazing/helper"
	"synergazing.com/synergazing/service"
)

type NotificationController struct {
	notificationService *service.NotificationService
}

func NewNotificationController(ns *service.NotificationService) *NotificationController {
	return &NotificationController{notificationService: ns}
}

// GetNotifications retrieves paginated notifications for the authenticated user
func (ctrl *NotificationController) GetNotifications(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	// Get pagination parameters
	limitStr := c.Query("limit", "20")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	notifications, err := ctrl.notificationService.GetUserNotifications(userID, limit, offset)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, fiber.Map{
		"notifications": notifications,
		"limit":         limit,
		"offset":        offset,
	}, "Notifications retrieved successfully")
}

// GetUnreadNotifications retrieves unread notifications for the authenticated user
func (ctrl *NotificationController) GetUnreadNotifications(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	notifications, err := ctrl.notificationService.GetUnreadNotifications(userID)
	if err != nil {
		return helper.Message400(err.Error())
	}

	unreadCount, err := ctrl.notificationService.GetUnreadCount(userID)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, fiber.Map{
		"notifications":      notifications,
		"unread_count":       unreadCount,
		"notification_count": len(notifications),
	}, "Unread notifications retrieved successfully")
}

// GetUnreadCount retrieves the count of unread notifications
func (ctrl *NotificationController) GetUnreadCount(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	count, err := ctrl.notificationService.GetUnreadCount(userID)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, fiber.Map{
		"unread_count": count,
	}, "Unread count retrieved successfully")
}

// MarkAsRead marks a specific notification as read
func (ctrl *NotificationController) MarkAsRead(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	notificationID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return helper.Message400("Invalid notification ID")
	}

	err = ctrl.notificationService.MarkAsRead(uint(notificationID), userID)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, nil, "Notification marked as read")
}

// MarkAllAsRead marks all notifications as read for the authenticated user
func (ctrl *NotificationController) MarkAllAsRead(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	err := ctrl.notificationService.MarkAllAsRead(userID)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, nil, "All notifications marked as read")
}

// DeleteNotification deletes a specific notification
func (ctrl *NotificationController) DeleteNotification(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	notificationID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return helper.Message400("Invalid notification ID")
	}

	err = ctrl.notificationService.DeleteNotification(uint(notificationID), userID)
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, nil, "Notification deleted successfully")
}

// TestDeadlineNotifications manually triggers deadline notifications (for testing)
func (ctrl *NotificationController) TestDeadlineNotifications(c *fiber.Ctx) error {
	err := ctrl.notificationService.CheckAndNotifyApproachingDeadlines()
	if err != nil {
		return helper.Message400(err.Error())
	}

	return helper.Message200(c, nil, "Deadline notifications check completed")
}
