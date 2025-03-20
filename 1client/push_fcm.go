package client

import (
    "fmt"
    "github.com/go-toast/toast"
)

func SendNotification(title, message string) {
	actions := []toast.Action{
		{"protocol", "Да", "http://localhost:8081/yes"},
		{"protocol", "Нет", "http://localhost:8081/no"},
	}

    notification := toast.Notification{
        AppID:   "MyCLIApp",
        Title:   title,
        Message: message,
        Actions: actions,
    }

    err := notification.Push()
    if err != nil {
        fmt.Printf("Ошибка при отправке уведомления: %v\n", err)
        return
    }
}