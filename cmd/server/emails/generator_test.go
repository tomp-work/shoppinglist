package emails_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tomp-work/shoppinglist/cmd/server/emails"
	"github.com/tomp-work/shoppinglist/cmd/server/models"
)

const expected = `Dear Clare,

Please can you pick up the following shopping:

- Wine (£15)
- Beer (£10)
- Gin (£25)

The total price of the shop should be £50.

Thanks,

Tom
`

func TestGenerateShoppingListEmail(t *testing.T) {
	items := []*models.Item{
		{SeqNum: 0, Name: "Wine", Price: 15},
		{SeqNum: 1, Name: "Beer", Price: 10},
		{SeqNum: 2, Name: "Gin", Price: 25},
	}
	details := models.ListDetails{
		TotalPrice: 50,
	}
	content, err := emails.GenerateShoppingListEmail("Clare", "Tom", details, items)
	require.NoError(t, err)
	require.Equal(t, expected, content)
}
