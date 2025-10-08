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
		URL:    fmt.Sprintf("%s%s/%s/attachment/%s", InvoiceAPIEndpoint, InvoicesModule, invoiceId, attachID),
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

func (c *API) ListInvoiceAttachments(invoiceId string) (data ListAttachmentResponse, err error) {

	endpoint := zoho.Endpoint{
		URL:    fmt.Sprintf("%s%s/%s/attachment", InvoiceAPIEndpoint, InvoicesModule, invoiceId),
		Method: zoho.HTTPGet,
		URLParameters: map[string]zoho.Parameter{
			"filter_by": "",
		},
		BodyFormat: zoho.JSON_STRING,
		Headers: map[string]string{
			InvoiceAPIEndpointHeader: c.OrganizationID,
		},
		ResponseData: &ListAttachmentResponse{},
	}

	err = c.Zoho.HTTPRequest(&endpoint)
	if err != nil {
		return ListAttachmentResponse{}, fmt.Errorf("Failed to list invoice attachments: %s", err)
	}

	if v, ok := endpoint.ResponseData.(*ListAttachmentResponse); ok {
		return *v, nil
	}
	return ListAttachmentResponse{}, fmt.Errorf("Data retrieved was not 'ListAttachmentResponse'")
}

type ListAttachmentResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    []AttachmentResponse `json:"data"`
}
type AttachmentResponse struct {
	ID          string          `json:"id"`
	CreatedTime string          `json:"created_time"`
	Owner       AttachmentOwner `json:"owner"`
	ParentID    string          `json:"parent_id"`
}
type AttachmentOwner struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	Email string `json:"email"`
}
