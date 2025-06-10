package events

import (
	"encoding/json"
	"lateslip/models"
)

func NotifyStudent(studentID string, message string) {
	clientManager.mu.Lock()
	defer clientManager.mu.Unlock()

	if client, exists := clientManager.studentClients[studentID]; exists {
		// If message is already JSON, use it directly
		client.Send <- []byte(message)
	}
}

func NotifyAdmins(message string, lateSlip models.LateSlip) {
	msg := map[string]interface{}{
		"type":    "NEW_LATE_SLIP_REQUEST",
		"message": message,
		"data": map[string]interface{}{
			"id":        lateSlip.ID.Hex(),
			"studentId": lateSlip.StudentID.Hex(),
			"reason":    lateSlip.Reason,
			"status":    lateSlip.Status,
			"createdAt": lateSlip.CreatedAt,
		},
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return
	}

	clientManager.mu.Lock()
	defer clientManager.mu.Unlock()

	for client := range clientManager.adminClients {
		client.Send <- jsonMsg
	}
}