package tenant

import (
	"fmt"

	"github.com/mdryaan/resgate/internal/config"
	"github.com/mdryaan/resgate/internal/models"
)

func ValidatePriority(priority int) error {
	if priority < config.MinPriority || priority > config.MaxPriority {
		return fmt.Errorf("priority must be between %d and %d", config.MinPriority, config.MaxPriority)
	}
	return nil
}

func HasHigherPriority(a, b *models.Tenant) bool {
	return a.Priority > b.Priority
}
