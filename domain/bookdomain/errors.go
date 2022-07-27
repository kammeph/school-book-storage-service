package bookdomain

import "fmt"

func ErrApplyEventBookAlreadyExists(eventType, id string) error {
	return fmt.Errorf("can not apply %s: Book with ID %s already exists", eventType, id)
}

func ErrApplyEventBookWithIDNotFound(eventType, id string) error {
	return fmt.Errorf("can not apply %s: Book with ID %s not found", eventType, id)
}
