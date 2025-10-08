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
func (c *API) DeleteInvoiceFile(request interface{}, invoiceId, attachID string) (data DeleteAttachmentResponse, err error) {
	endpoint := zoho.Endpoint{
		URL:    fmt.Sprintf("%s%s/%s/documents/%s", InvoiceAPIEndpoint, InvoicesModule, invoiceId, attachID),
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

type DeleteAttachmentResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
