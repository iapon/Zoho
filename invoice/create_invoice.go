package invoice

import (
	"fmt"
	"log"

	zoho "github.com/iapon/zoho"
	"github.com/kr/pretty"
)

//https://www.zoho.com/invoice/api/v3/#Invoices_Create_an_invoice
//func (c *API) CreateInvoice(request interface{}, OrganizationID string, params map[string]zoho.Parameter) (data ListContactsResponse, err error) {
func (c *API) CreateInvoice(request interface{}, pars map[string]zoho.Parameter, mark bool) (data CreateInvoiceResponse, err error) {

	endpoint := zoho.Endpoint{
		Name:         InvoicesModule,
		URL:          fmt.Sprintf(InvoiceAPIEndpoint+"%s", InvoicesModule),
		Method:       zoho.HTTPPost,
		ResponseData: &CreateInvoiceResponse{},
		URLParameters: map[string]zoho.Parameter{
			"filter_by": "",
		},
		RequestBody: request,
		BodyFormat:  zoho.JSON_STRING,
		Headers: map[string]string{
			InvoiceAPIEndpointHeader: c.OrganizationID,
		},
	}
	if pars != nil && len(pars) > 0 {
		for k, v := range pars {
			endpoint.URLParameters[k] = v
		}
	}

	/*for k, v := range params {
		endpoint.URLParameters[k] = v
	}*/

	log.Printf("%# v", pretty.Formatter(request))
	err = c.Zoho.HTTPRequest(&endpoint)
	if err != nil {
		return CreateInvoiceResponse{}, fmt.Errorf("Failed to create invoice: %s", err)
	}

	// Mark the invoice as sent before returning details
	if v, ok := endpoint.ResponseData.(*CreateInvoiceResponse); ok {
		// Check if the creation succeeded
		if v.Code != 0 {
			return *v, fmt.Errorf("Failed to create invoice: %s", v.Message)
		}
		if mark {
			endpointSent := zoho.Endpoint{
				Name:         InvoicesModule,
				URL:          fmt.Sprintf(InvoiceAPIEndpoint+"%s/%s/status/sent", InvoicesModule, v.Invoice.InvoiceId),
				Method:       zoho.HTTPPost,
				ResponseData: &InvoiceSent{},
				BodyFormat:   zoho.JSON_STRING,
				Headers: map[string]string{
					InvoiceAPIEndpointHeader: c.OrganizationID,
				},
			}
			err = c.Zoho.HTTPRequest(&endpointSent)
			if err != nil {
				return *v, fmt.Errorf("Failed to mark invoice as sent: %s", err)
			}
		}
		return *v, nil
	}

	return CreateInvoiceResponse{}, fmt.Errorf("Data retrieved was not 'CreateInvoiceResponse'")
}
func (c *API) SetSent(invoiceId string) error {
	endpointSent := zoho.Endpoint{
		Name:         InvoicesModule,
		URL:          fmt.Sprintf(InvoiceAPIEndpoint+"%s/%s/status/sent", InvoicesModule, invoiceId),
		Method:       zoho.HTTPPost,
		ResponseData: &InvoiceSent{},
		BodyFormat:   zoho.JSON_STRING,
		Headers: map[string]string{
			InvoiceAPIEndpointHeader: c.OrganizationID,
		},
	}
	err := c.Zoho.HTTPRequest(&endpointSent)
	if err != nil {
		return fmt.Errorf("Failed to mark invoice as sent: %s", err)
	}
	return nil
}

type CreateInvoiceRequest struct {
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
	PaymentOptions        PaymentOptions       `json:"payment_options,omitempty"`
	AllowPartialPayments  bool                 `json:"allow_partial_payments,omitempty"`
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

type InvoiceSent struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

type CreateInvoiceResponse struct {
	Code    int64   `json:"code"`
	Message string  `json:"message"`
	Invoice Invoice `json:"invoice"`
}
type Invoice struct {
	InvoiceId             string            `json:"invoice_id,omitempty"`
	AchPaymentInitiated   bool              `json:"ach_payment_initiated,omitempty"`
	InvoiceNumber         string            `json:"invoice_number,omitempty"`
	IsPreGst              bool              `json:"is_pre_gst,omitempty"`
	PlaceOfSupply         string            `json:"place_of_supply,omitempty"`
	GstNo                 string            `json:"gst_no,omitempty"`
	GstTreatment          string            `json:"gst_treatment,omitempty"`
	Date                  string            `json:"date,omitempty"`
	Status                string            `json:"status,omitempty"`
	PaymentTerms          int64             `json:"payment_terms,omitempty"`
	PaymentTermsLabel     string            `json:"payment_terms_label,omitempty"`
	DueDate               string            `json:"due_date,omitempty"`
	PaymentExpectedDate   string            `json:"payment_expected_date,omitempty"`
	LastPaymentDate       string            `json:"last_payment_date,omitempty"`
	ReferenceNumber       string            `json:"reference_number,omitempty"`
	CustomerId            string            `json:"customer_id,omitempty"`
	CustomerName          string            `json:"customer_name,omitempty"`
	ContactPersons        []string          `json:"contact_persons,omitempty"`
	CurrencyId            string            `json:"currency_id,omitempty"`
	CurrencyCode          string            `json:"currency_code,omitempty"`
	ExchangeRate          float64           `json:"exchange_rate,omitempty"`
	Discount              float64           `json:"discount,omitempty"`
	IsDiscountBeforeTax   bool              `json:"is_discount_before_tax,omitempty"`
	DiscountType          string            `json:"discount_type,omitempty"`
	IsInclusiveTax        bool              `json:"is_inclusive_tax,omitempty"`
	RecurringInvoiceId    string            `json:"recurring_invoice_id,omitempty"`
	IsViewedByClient      bool              `json:"is_viewed_by_client,omitempty"`
	HasAttachment         bool              `json:"has_attachment,omitempty"`
	ClientViewedTime      string            `json:"client_viewed_time,omitempty"`
	LineItems             []InvoiceLineItem `json:"line_items,omitempty"`
	ShippingCharge        float64           `json:"shipping_charge,omitempty"`
	Adjustment            float64           `json:"adjustment,omitempty"`
	AdjustmentDescription string            `json:"adjustment_description,omitempty"`
	SubTotal              float64           `json:"sub_total,omitempty"`
	TaxTotal              float64           `json:"tax_total,omitempty"`
	Total                 float64           `json:"total,omitempty"`
	Taxes                 []struct {
		TaxName   string  `json:"tax_name,omitempty"`
		TaxAmount float64 `json:"tax_amount,omitempty"`
	} `json:"taxes,omitempty"`
	PaymentReminderEnabled bool           `json:"payment_reminder_enabled,omitempty"`
	PaymentMade            float64        `json:"payment_made,omitempty"`
	CreditsApplied         float64        `json:"credits_applied,omitempty"`
	TaxAmountWithheld      float64        `json:"tax_amount_withheld,omitempty"`
	Balance                float64        `json:"balance,omitempty"`
	WriteOffAmount         float64        `json:"write_off_amount,omitempty"`
	AllowPartialPayments   bool           `json:"allow_partial_payments,omitempty"`
	PricePrecision         int64          `json:"price_precision,omitempty"`
	PaymentOptions         PaymentOptions `json:"payment_options,omitempty"`
	IsEmailed              bool           `json:"is_emailed,omitempty"`
	RemindersSent          int64          `json:"reminders_sent,omitempty"`
	LastReminderSentDate   string         `json:"last_reminder_sent_date,omitempty"`
	BillingAddress         ContactAddress `json:"billing_address,omitempty"`
	ShippingAddress        ContactAddress `json:"shipping_address,omitempty"`
	Notes                  string         `json:"notes,omitempty"`
	Terms                  string         `json:"terms,omitempty"`
	CustomFields           []struct {
		CustomfieldId string      `json:"customfield_id,omitempty"`
		DataType      string      `json:"data_type,omitempty"`
		Index         int64       `json:"index,omitempty"`
		Label         string      `json:"label,omitempty"`
		ShowOnPdf     bool        `json:"show_on_pdf,omitempty"`
		ShowInAllPdf  bool        `json:"show_in_all_pdf,omitempty"`
		Value         interface{} `json:"value,omitempty"`
	} `json:"custom_fields,omitempty"`
	TemplateId       string `json:"template_id,omitempty"`
	TemplateName     string `json:"template_name,omitempty"`
	CreatedTime      string `json:"created_time,omitempty"`
	LastModifiedTime string `json:"last_modified_time,omitempty"`
	AttachmentName   string `json:"attachment_name,omitempty"`
	CanSendInMail    bool   `json:"can_send_in_mail,omitempty"`
	SalespersonId    string `json:"salesperson_id,omitempty"`
	SalespersonName  string `json:"salesperson_name,omitempty"`
	InvoiceUrl       string `json:"invoice_url,omitempty"`
}
type InvoiceLineItem struct {
	LineItemId       string               `json:"line_item_id,omitempty"`
	ItemId           string               `json:"item_id,omitempty"`
	Description      string               `json:"description,omitempty"`
	ProjectId        string               `json:"project_id,omitempty"`
	ProjectName      string               `json:"project_name,omitempty"`
	TimeEntryIds     []string             `json:"time_entry_ids,omitempty"`
	ItemType         string               `json:"item_type,omitempty"`
	ProductType      string               `json:"product_type,omitempty"`
	ExpenseId        string               `json:"expense_id,omitempty"`
	Name             string               `json:"name,omitempty"`
	ItemOrder        float64              `json:"item_order,omitempty"`
	BcyRate          float64              `json:"bcy_rate,omitempty"`
	Rate             float64              `json:"rate,omitempty"`
	Quantity         float64              `json:"quantity,omitempty"`
	Unit             string               `json:"unit,omitempty"`
	DiscountAmount   float64              `json:"discount_amount,omitempty"`
	Discount         float64              `json:"discount,omitempty"`
	TaxId            string               `json:"tax_id,omitempty"`
	TaxExemptionId   string               `json:"tax_exemption_id,omitempty"`
	TaxName          string               `json:"tax_name,omitempty"`
	TaxType          string               `json:"tax_type,omitempty"`
	TaxPercentage    float64              `json:"tax_percentage,omitempty"`
	ItemTotal        float64              `json:"item_total,omitempty"`
	HsnOrSac         int64                `json:"hsn_or_sac,omitempty"`
	ItemCustomFields []CustomFieldRequest `json:"item_custom_fields,omitempty"`
}
