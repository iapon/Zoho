package invoice

import (
	"fmt"

	zoho "github.com/iapon/zoho"
)

// https://www.zohoapis.com/invoice/v3/contacts/460000000026049/contactpersons/460000000026051 \
func (c *API) GetContactPerson(contactID, contactPersonID string) (data GetContactPersonResponse, err error) {
	endpoint := zoho.Endpoint{
		Name:         ContactsModule,
		URL:          fmt.Sprintf("https://www.zohoapis.%s/invoice/v3/%s/%s/%s", c.ZohoTLD, ContactsModule, ContactsPersonSubModule, contactPersonID),
		Method:       zoho.HTTPGet,
		ResponseData: &GetContactPersonResponse{},
		URLParameters: map[string]zoho.Parameter{
			"filter_by": "",
		},
		BodyFormat: zoho.JSON_STRING,
		Headers: map[string]string{
			InvoiceAPIEndpointHeader: c.OrganizationID,
		},
	}

	err = c.Zoho.HTTPRequest(&endpoint)
	if err != nil {
		return GetContactPersonResponse{}, fmt.Errorf("Failed to retrieve contact: %s", err)
	}

	if v, ok := endpoint.ResponseData.(*GetContactPersonResponse); ok {
		// Check if the request succeeded
		if v.Code != 0 {
			return *v, fmt.Errorf("Failed to retrieve contact: %s", v.Message)
		}
		return *v, nil
	}
	return GetContactPersonResponse{}, fmt.Errorf("Data retrieved was not 'GetContactPersonResponse'")
}

type GetContactPersonResponse struct {
	Code          int           `json:"code"`
	Message       string        `json:"message"`
	ContactPerson ContactPerson `json:"contact_person"`
}
