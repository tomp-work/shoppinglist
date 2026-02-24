package emails

import (
	"fmt"

	"github.com/tomp-work/shoppinglist/cmd/server/models"
)

func GenerateShoppingListEmail(toName, fromName string, details models.ListDetails, items []*models.Item) (string, error) {
	if toName == "" || fromName == "" {
		return "", fmt.Errorf("GenerateShoppingListEmail: missing name")
	}
	if len(items) == 0 {
		return "", fmt.Errorf("GenerateShoppingListEmail: missing items")
	}
	content := fmt.Sprintf("Dear %s,\n", toName)
	content += "\n"
	content += "Please can you pick up the following shopping:\n"
	content += "\n"
	for _, item := range items {
		content += fmt.Sprintf("- %s (£%d)\n", item.Name, item.Price)
	}
	content += "\n"
	content += fmt.Sprintf("The total price of the shop should be £%d.\n", details.TotalPrice)
	content += "\n"
	content += "Thanks,\n"
	content += "\n"
	content += fmt.Sprintf("%s\n", fromName)
	return content, nil
}
