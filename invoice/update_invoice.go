package invoice

import (
	"encoding/json"
	"fmt"

	zoho "github.com/iapon/zoho"
)

//https://www.zoho.com/invoice/api/v3/#Invoices_Update_an_invoice
//func (c *API) UpdateRecurringInvoice(request interface{}, OrganizationID string, params map[string]zoho.Parameter) (data UpdateInvoiceResponse, err error) {
func (c *API) UpdateInvoice(request interface{}, invoiceId string) (data UpdateInvoiceResponse, err error) {
	endpoint := zoho.Endpoint{
		Name:         ContactsModule,
		URL:          fmt.Sprintf("%s%s/%s", InvoiceAPIEndpoint, InvoicesModule, invoiceId),
		Method:       zoho.HTTPPut,
		ResponseData: &UpdateInvoiceResponse{},
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
		return UpdateInvoiceResponse{}, fmt.Errorf("Failed to update invoice: %s", err)
	}

	if v, ok := endpoint.ResponseData.(*UpdateInvoiceResponse); ok {
		// Check if the request succeeded
		if v.Code != 0 {
			return *v, fmt.Errorf("Failed to update invoice: %s", v.Message)
		}
		return *v, nil
	}
	return UpdateInvoiceResponse{}, fmt.Errorf("Data retrieved was not 'UpdateInvoiceResponse'")
}

type UpdateInvoiceRequest struct {
	CustomerId         string               `json:"customer_id"`
	InvoicedEstimateId string               `json:"invoiced_estimate_id,omitempty"`
	CustomFields       []CustomFieldRequest `json:"custom_fields,omitempty"`
	ProjectId          string               `json:"project_id,omitempty"`
	CustomBody         string               `json:"custom_body,omitempty"`
	CustomSubject      string               `json:"custom_subject,omitempty"`
	Reason             string               `json:"reason,omitempty"`
	TaxAuthorityId     string               `json:"tax_authority_id,omitempty"`
	TaxExemptionId     string               `json:"tax_exemption_id,omitempty"`
	Invoice            Invoice              `json:"invoice"`
}

type UpdateInvoiceResponse struct {
	Code    int64   `json:"code"`
	Message string  `json:"message"`
	Invoice Invoice `json:"invoice"`
}

type EmailInvoiceResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func (u *UpdateInvoiceRequest) MarshalJSON() ([]byte, error) {
	type updateInvoiceRequest struct {
		CustomerId string `json:"customer_id"`
		//ContactName          string   `json:"contact_name,omitempty"`
		ContactPersons        []string             `json:"contact_persons,omitempty"`
		InvoiceNumber         string               `json:"invoice_number,omitempty"`
		ReferenceNumber       string               `json:"reference_number,omitempty"`
		PlaceOfSupply         string               `json:"place_of_supply,omitempty"`
		GstTreatment          string               `json:"gst_treatment,omitempty"`
		GstNo                 string               `json:"gst_no,omitempty"`
		TemplateId            string               `json:"template_id,omitempty"`
		Date                  string               `json:"date,omitempty"`
		PaymentTerms          int64                `json:"payment_terms,omitempty"`
		PaymentTermsLabel     string               `json:"payment_terms_label,omitempty"`
		DueDate               string               `json:"due_date,omitempty"`
		Discount              float64              `json:"discount,omitempty"`
		IsDiscountBeforeTax   bool                 `json:"is_discount_before_tax,omitempty"`
		DiscountType          string               `json:"discount_type,omitempty"`
		IsInclusiveTax        bool                 `json:"is_inclusive_tax,omitempty"`
		ExchangeRate          float64              `json:"exchange_rate,omitempty"`
		RecurringInvoiceId    string               `json:"recurring_invoice_id,omitempty"`
		InvoicedEstimateId    string               `json:"invoiced_estimate_id,omitempty"`
		SalespersonId         string               `json:"salesperson_id,omitempty"`
		CustomFields          []CustomFieldRequest `json:"custom_fields,omitempty"`
		ProjectId             string               `json:"project_id,omitempty"`
		LineItems             []InvoiceLineItem    `json:"line_items"`
		PaymentOptions        PaymentOptions       `json:"payment_options"`
		AllowPartialPayments  bool                 `json:"allow_partial_payments"`
		CustomBody            string               `json:"custom_body,omitempty"`
		CustomSubject         string               `json:"custom_subject,omitempty"`
		Notes                 string               `json:"notes,omitempty"`
		Terms                 string               `json:"terms,omitempty"`
		ShippingCharge        float64              `json:"shipping_charge,omitempty"`
		Adjustment            float64              `json:"adjustment,omitempty"`
		AdjustmentDescription string               `json:"adjustment_description"`
		Reason                string               `json:"reason,omitempty"`
		TaxAuthorityId        string               `json:"tax_authority_id,omitempty"`
		TaxExemptionId        string               `json:"tax_exemption_id,omitempty"`
	}
	uir := updateInvoiceRequest{
		ContactPersons:        u.Invoice.ContactPersons,
		InvoiceNumber:         u.Invoice.InvoiceNumber,
		ReferenceNumber:       u.Invoice.ReferenceNumber,
		PlaceOfSupply:         u.Invoice.PlaceOfSupply,
		GstTreatment:          u.Invoice.GstTreatment,
		GstNo:                 u.Invoice.GstNo,
		TemplateId:            u.Invoice.TemplateId,
		Date:                  u.Invoice.Date,
		PaymentTerms:          u.Invoice.PaymentTerms,
		PaymentTermsLabel:     u.Invoice.PaymentTermsLabel,
		DueDate:               u.Invoice.DueDate,
		Discount:              u.Invoice.Discount,
		IsDiscountBeforeTax:   u.Invoice.IsDiscountBeforeTax,
		DiscountType:          u.Invoice.DiscountType,
		IsInclusiveTax:        u.Invoice.IsInclusiveTax,
		ExchangeRate:          u.Invoice.ExchangeRate,
		RecurringInvoiceId:    u.Invoice.RecurringInvoiceId,
		InvoicedEstimateId:    u.InvoicedEstimateId,
		SalespersonId:         u.Invoice.SalespersonId,
		ProjectId:             u.ProjectId,
		LineItems:             u.Invoice.LineItems,
		PaymentOptions:        u.Invoice.PaymentOptions,
		AllowPartialPayments:  u.Invoice.AllowPartialPayments,
		CustomBody:            u.CustomBody,
		CustomSubject:         u.CustomSubject,
		Notes:                 u.Invoice.Notes,
		Terms:                 u.Invoice.Terms,
		ShippingCharge:        u.Invoice.ShippingCharge,
		Adjustment:            u.Invoice.Adjustment,
		AdjustmentDescription: u.Invoice.AdjustmentDescription,
		Reason:                u.Reason,
		TaxAuthorityId:        u.TaxAuthorityId,
		TaxExemptionId:        u.TaxExemptionId,
	}
	return json.Marshal(uir)
}
func (c *API) EmailInvoice(request interface{}, invoiceId string) (data EmailInvoiceResponse, err error) {
	endpoint := zoho.Endpoint{
		URL:          fmt.Sprintf("%s%s/%s/email", InvoiceAPIEndpoint, InvoicesModule, invoiceId),
		Method:       zoho.HTTPPost,
		ResponseData: &EmailInvoiceResponse{},
		URLParameters: map[string]zoho.Parameter{
			"filter_by": "",
		},
		RequestBody: &request,
		BodyFormat:  zoho.JSON,
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
func (c *API) EmailInvoiceWithFile(request interface{}, invoiceId string, file []byte, filename string) (data EmailInvoiceResponse, err error) {
	endpoint := zoho.Endpoint{
		URL:          fmt.Sprintf("%s%s/%s/email", InvoiceAPIEndpoint, InvoicesModule, invoiceId),
		Method:       zoho.HTTPPost,
		ResponseData: &EmailInvoiceResponse{},
		URLParameters: map[string]zoho.Parameter{
			"filter_by":       "",
			"send_attachment": zoho.Parameter("true"),
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
