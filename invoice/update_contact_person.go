package invoice

import (
	"fmt"

	zoho "github.com/iapon/zoho"
)

// https://www.zoho.com/invoice/api/v3/contact-persons/#update-a-contact-person
func (c *API) UpdateContactPerson(request any, contactPersonID string) (data UpdateContactPersonResponse, err error) {
	endpoint := zoho.Endpoint{
		Name:         InvoicesModule,
		URL:          fmt.Sprintf("%s%s/%s/%s", InvoiceAPIEndpoint, ContactsModule, ContactsPersonSubModule, contactPersonID),
		Method:       zoho.HTTPPut,
		ResponseData: &UpdateContactPersonResponse{},
		URLParameters: map[string]zoho.Parameter{
			"filter_by": "",
		},
		RequestBody: &request,
		BodyFormat:  zoho.JSON_STRING,
		Headers: map[string]string{
			InvoiceAPIEndpointHeader: c.OrganizationID,
		},
	}

	/*for k, v := range params {
		endpoint.URLParameters[k] = v
	}*/

	err = c.Zoho.HTTPRequest(&endpoint)
	if err != nil {
		return UpdateContactPersonResponse{}, fmt.Errorf("Failed to create contact: %s", err)
	}

	if v, ok := endpoint.ResponseData.(*UpdateContactPersonResponse); ok {
		// Check if the request succeeded
		if v.Code != 0 {
			return *v, fmt.Errorf("Failed to update contact: %s", v.Message)
		}
		return *v, nil
	}
	return UpdateContactPersonResponse{}, fmt.Errorf("Data retrieved was not 'UpdateContactResponse'")
}

type UpdateContactPersonRequest struct {
	ContactPerson
	EnablePortal *bool `json:"enable_portal"`
}

type UpdateContactPersonResponse struct {
	Code          int             `json:"code"`
	Message       string          `json:"message"`
	ContactPerson []ContactPerson `json:"contact_person"`
}
