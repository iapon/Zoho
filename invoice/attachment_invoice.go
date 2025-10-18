package invoice

import (
	"fmt"

	zoho "github.com/iapon/zoho"
)

func (c *API) AttachInvoiceFile(request interface{}, invoiceId string, file []byte, filename string) (data EmailInvoiceResponse, err error) {
	endpoint := zoho.Endpoint{
		URL:          fmt.Sprintf("%s%s/%s/attachment", InvoiceAPIEndpoint, InvoicesModule, invoiceId),
		Method:       zoho.HTTPPost,
		ResponseData: &EmailInvoiceResponse{},
		URLParameters: map[string]zoho.Parameter{
			"can_send_in_mail": zoho.Parameter("true"),
			"filter_by":        "",
			"send_attachment":  zoho.Parameter("true"),
		},
		RequestBody:    &request,
		BodyFormat:     zoho.FILE_BYTE,
		AttachmentByte: file,
		Attachment:     filename,
		Headers: map[string]string{
			InvoiceAPIEndpointHeader: c.OrganizationID,
		},
	}
	err = c.Zoho.HTTPRequest(&endpoint)
	if err != nil {
		return EmailInvoiceResponse{}, fmt.Errorf("failed to update invoice: %s", err)
	}

	if v, ok := endpoint.ResponseData.(*EmailInvoiceResponse); ok {
		// Check if the request succeeded
		if v.Code != 0 {
			return *v, fmt.Errorf("failed to update invoice: %s", v.Message)
		}
		return *v, nil
	}
	return EmailInvoiceResponse{}, fmt.Errorf("data retrieved was not 'EmailInvoiceResponse'")
}
func (c *API) DeleteInvoiceFile(request interface{}, invoiceId string) (data DeleteAttachmentResponse, err error) {
	endpoint := zoho.Endpoint{
		URL:    fmt.Sprintf("%s%s/%s/attachment", InvoiceAPIEndpoint, InvoicesModule, invoiceId),
		Method: zoho.HTTPDelete,
		URLParameters: map[string]zoho.Parameter{
			"filter_by": "",
		},
		RequestBody: &request,
		BodyFormat:  zoho.JSON_STRING,
		Headers: map[string]string{
			InvoiceAPIEndpointHeader: c.OrganizationID,
		},
		ResponseData: &DeleteAttachmentResponse{},
	}

	err = c.Zoho.HTTPRequest(&endpoint)
	if err != nil {
		return DeleteAttachmentResponse{}, fmt.Errorf("Failed to delete file: %s", err)
	}

	if v, ok := endpoint.ResponseData.(*DeleteAttachmentResponse); ok {
		// Check if the request succeeded
		if v.Code != 0 {
			return *v, fmt.Errorf("Failed to delete file: %s", v.Message)
		}
		return *v, nil
	}
	return DeleteAttachmentResponse{}, fmt.Errorf("Data retrieved was not 'DeleteAttachmentResponse'")
}
func (c *API) GetInvoicePDF(invoiceId string) ([]byte, error) {
	err := c.CheckForSavedTokens()
	if err == zoho.ErrTokenExpired {
		err := c.RefreshTokenRequest()
		if err != nil {
			return nil, fmt.Errorf("Failed to refresh the access token: %s: %s", InvoicesModule, err)
		}
	}
	client := &http.Client{}
	endpointURL := fmt.Sprintf("%s%s/%s/attachment", InvoiceAPIEndpoint, InvoicesModule, invoiceId),
	q := url.Values{}
	q.Set("organization_id", c.OrganizationID)
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", endpointURL, q.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create a request for %s: %s", InvoicesModule, err)
	}

	// Add global authorization header
	req.Header.Add("Authorization", "Zoho-oauthtoken "+c.GetOauthToken())
	req.Header.Add(InvoiceAPIEndpointHeader, c.OrganizationID)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to perform request for %s: %s", InvoicesModule, err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read body of response for %s: got status %s: %s", InvoicesModule, zoho.ResolveStatus(resp), err)
	}
	return body, nil
}


type DeleteAttachmentResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
