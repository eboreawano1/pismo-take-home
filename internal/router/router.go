package router

import "pismo-take-home/internal/event"

func RouteEvent(e event.Event) string {
	switch e.EventType {
		case "payment_authorized":
			return "analytics"
		case "user_sign_up":
			return "notifications"
		default:
			return ""
	}
}