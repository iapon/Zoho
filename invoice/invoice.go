package invoice

import (
	"math/rand"

	zoho "github.com/iapon/zoho"
)

const (
	ContactsModule           string = "contacts"
	ContactsPersonSubModule  string = "contactpersons"
	InvoicesModule           string = "invoices"
	InvoiceAPIEndpointHeader string = "X-com-zoho-invoice-organizationid"
	ItemsModule              string = "items"
	RecurringInvoicesModule  string = "recurringinvoices"
	CustomerPaymentsModule   string = "customerpayments"
)

var InvoiceAPIEndpoint string = "https://www.zohoapis.com/invoice/v3/"

type CustomFieldRequest struct {
	CustomfieldID string      `json:"customfield_id,omitempty"`
	Label         string      `json:"label"`
	Value         interface{} `json:"value,omitempty"`
}

// API is used for interacting with the Zoho expense API
// the exposed methods are primarily access to expense modules which provide access to expense Methods
type API struct {
	*zoho.Zoho
	id string
}

func (c *API) SetBooking() {
	InvoiceAPIEndpoint = "https://www.zohoapis.com/books/v3/"
}

// New returns a *invoice.API with the provided zoho.Zoho as an embedded field
func New(z *zoho.Zoho) *API {
	id := func() string {
		var id []byte
		keyspace := "abcdefghijklmnopqrutuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		for i := 0; i < 25; i++ {
			id = append(id, keyspace[rand.Intn(len(keyspace))])
		}
		return string(id)
	}()

	API := API{
		Zoho: z,
		id:   id,
	}
	return &API
}
