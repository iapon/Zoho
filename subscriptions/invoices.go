package subscriptions

import (
	"fmt"
	"strconv"

	zoho "github.com/iapon/zoho"
)

type InvoiceStatus string

// Proper names for Invoice statuses
const (
	InvoiceStatusAll           InvoiceStatus = "Status.All"
	InvoiceStatusSent          InvoiceStatus = "Status.Sent"
	InvoiceStatusDraft         InvoiceStatus = "Status.Draft"
	InvoiceStatusOverDue       InvoiceStatus = "Status.OverDue"
	InvoiceStatusPaid          InvoiceStatus = "Status.Paid"
	InvoiceStatusPartiallyPaid InvoiceStatus = "Status.PartiallyPaid"
	InvoiceStatusVoid          InvoiceStatus = "Status.Void"
	InvoiceStatusUnpaid        InvoiceStatus = "Status.Unpaid"
	InvoiceStatusPending       InvoiceStatus = "Status.Pending" // Pending status is not present in zoho documentation, but works
)

// listInvoicesWithParams will return the list of invoices that match the given invoice status
// and additional filter defined by parameter name and value (allows to filter by `customer_id` and `subscription_id`)
// https://www.zoho.com/subscriptions/api/v1/#Invoices_List_all_invoices
func (s *API) listInvoicesWithParams(status InvoiceStatus, paramName, paramValue string) (data InvoicesResponse, err error) {
	if status == "" {
		status = InvoiceStatusAll
	}
	endpoint := zoho.Endpoint{
		Name:         "invoices",
		URL:          fmt.Sprintf("https://subscriptions.zoho.%s/api/v1/invoices", s.ZohoTLD),
		Method:       zoho.HTTPGet,
		ResponseData: &InvoicesResponse{},
		URLParameters: map[string]zoho.Parameter{
			"filter_by": zoho.Parameter(status),
		},
		Headers: map[string]string{
			ZohoSubscriptionsEndpointHeader: s.OrganizationID,
		},
	}

	if paramName != "" {
		endpoint.URLParameters[paramName] = zoho.Parameter(paramValue)
	}

	err = s.Zoho.HTTPRequest(&endpoint)
	if err != nil {
		return InvoicesResponse{}, fmt.Errorf("Failed to retrieve invoices: %s", err)
	}

	if v, ok := endpoint.ResponseData.(*InvoicesResponse); ok {
		return *v, nil
	}

	return InvoicesResponse{}, fmt.Errorf("Data retrieved was not 'InvoicesResponse'")
}

// ListAllInvoices will return the list of invoices that match the given invoice status
// https://www.zoho.com/subscriptions/api/v1/#Invoices_List_all_invoices
func (s *API) ListAllInvoices(status InvoiceStatus) (data InvoicesResponse, err error) {
	return s.listInvoicesWithParams(status, "", "")
}

// ListInvoicesForSubscription will return the list of invoices that match the given invoice status and subscription ID
// https://www.zoho.com/subscriptions/api/v1/#Invoices_List_all_invoices
func (s *API) ListInvoicesForSubscription(status InvoiceStatus, subscriptionID string) (data InvoicesResponse, err error) {
	return s.listInvoicesWithParams(status, "subscription_id", subscriptionID)
}

// ListInvoicesForSubscription will return the list of invoices that match the given invoice status and customer ID
// https://www.zoho.com/subscriptions/api/v1/#Invoices_List_all_invoices
func (s *API) ListInvoicesForCustomer(status InvoiceStatus, customerID string) (data InvoicesResponse, err error) {
	return s.listInvoicesWithParams(status, "customer_id", customerID)
}

// GetInvoice will return the subscription specified by id
// https://www.zoho.com/subscriptions/api/v1/#Invoices_Retrieve_a_subscription
func (s *API) GetInvoice(id string) (data InvoiceResponse, err error) {
	endpoint := zoho.Endpoint{
		Name:         "invoices",
		URL:          fmt.Sprintf("https://subscriptions.zoho.%s/api/v1/invoices/%s", s.ZohoTLD, id),
		Method:       zoho.HTTPGet,
		ResponseData: &InvoiceResponse{},
		Headers: map[string]string{
			ZohoSubscriptionsEndpointHeader: s.OrganizationID,
		},
	}

	err = s.Zoho.HTTPRequest(&endpoint)
	if err != nil {
		return InvoiceResponse{}, fmt.Errorf("Failed to retrieve invoice (%s): %s", id, err)
	}

	if v, ok := endpoint.ResponseData.(*InvoiceResponse); ok {
		return *v, nil
	}

	return InvoiceResponse{}, fmt.Errorf("Data retrieved was not 'InvoiceResponse'")
}

// AddAttachment attaches a file to an invoice
// https://www.zoho.com/subscriptions/api/v1/#Invoices_Add_attachment_to_an_invoice
func (s *API) AddAttachment(id, file string, canSendInEmail bool) (data AttachementResponse, err error) {
	endpoint := zoho.Endpoint{
		Name:         "invoices",
		URL:          fmt.Sprintf("https://subscriptions.zoho.%s/api/v1/invoices/%s/attachment", s.ZohoTLD, id),
		Method:       zoho.HTTPPost,
		ResponseData: &AttachementResponse{},
		Attachment:   file,
		BodyFormat:   zoho.FILE,
		URLParameters: map[string]zoho.Parameter{
			"can_send_in_mail": zoho.Parameter(strconv.FormatBool(canSendInEmail)),
		},
		Headers: map[string]string{
			ZohoSubscriptionsEndpointHeader: s.OrganizationID,
		},
	}

	err = s.Zoho.HTTPRequest(&endpoint)
	if err != nil {
		return AttachementResponse{}, fmt.Errorf("Failed to attach file to invoice (%s): %s", id, err)
	}

	if v, ok := endpoint.ResponseData.(*AttachementResponse); ok {
		return *v, nil
	}

	return AttachementResponse{}, fmt.Errorf("Data retrieved was not 'AttachementResponse'")
}

// EmailInvoice sends an invoice in email
// https://www.zoho.com/subscriptions/api/v1/#Invoices_Email_an_invoice
func (s *API) EmailInvoice(id string, request EmailInvoiceRequest) (data EmailInvoiceResponse, err error) {
	endpoint := zoho.Endpoint{
		Name:         "invoices",
		URL:          fmt.Sprintf("https://subscriptions.zoho.%s/api/v1/invoices/%s/email", s.ZohoTLD, id),
		Method:       zoho.HTTPPost,
		ResponseData: &EmailInvoiceResponse{},
		RequestBody:  request,
		Headers: map[string]string{
			ZohoSubscriptionsEndpointHeader: s.OrganizationID,
		},
	}

	err = s.Zoho.HTTPRequest(&endpoint)
	if err != nil {
		return EmailInvoiceResponse{}, fmt.Errorf("Failed to email invoice (%s): %s", id, err)
	}

	if v, ok := endpoint.ResponseData.(*EmailInvoiceResponse); ok {
		return *v, nil
	}

	return EmailInvoiceResponse{}, fmt.Errorf("Data retrieved was not 'EmailInvoiceResponse'")
}

// AddItems adds items to pending invoice
// https://www.zoho.com/subscriptions/api/v1/#Invoices_Add_items_to_a_pending_invoice
func (s *API) AddItems(id string, request AddItemsRequest) (data AddItemsResponse, err error) {
	endpoint := zoho.Endpoint{
		Name:         "invoices",
		URL:          fmt.Sprintf("https://subscriptions.zoho.%s/api/v1/invoices/%s/lineitems", s.ZohoTLD, id),
		Method:       zoho.HTTPPost,
		ResponseData: &AddItemsResponse{},
		RequestBody:  request,
		Headers: map[string]string{
			ZohoSubscriptionsEndpointHeader: s.OrganizationID,
		},
	}

	err = s.Zoho.HTTPRequest(&endpoint)
	if err != nil {
		return AddItemsResponse{}, fmt.Errorf("Failed to add items to invoice (%s): %s", id, err)
	}

	if v, ok := endpoint.ResponseData.(*AddItemsResponse); ok {
		return *v, nil
	}

	return AddItemsResponse{}, fmt.Errorf("Data retrieved was not 'AddItemsResponse'")
}

// CollectChargeViaCreditCard collects charge via credit card
// https://www.zoho.com/subscriptions/api/v1/#Invoices_Collect_charge_via_credit_card
// Note: Real life reply for this request differs from Zoho documentation,
// so CollectChargeViaCreditCardResponse was updated to include both top level objects:
// - 'payment' per documentation
// - 'invoice' per real life reply
func (s *API) CollectChargeViaCreditCard(id string, request CollectChargeViaCreditCardRequest) (data CollectChargeViaCreditCardResponse, err error) {
	endpoint := zoho.Endpoint{
		Name:         "invoices",
		URL:          fmt.Sprintf("https://subscriptions.zoho.%s/api/v1/invoices/%s/collect", s.ZohoTLD, id),
		Method:       zoho.HTTPPost,
		ResponseData: &CollectChargeViaCreditCardResponse{},
		RequestBody:  request,
		Headers: map[string]string{
			ZohoSubscriptionsEndpointHeader: s.OrganizationID,
		},
	}

	err = s.Zoho.HTTPRequest(&endpoint)
	if err != nil {
		return CollectChargeViaCreditCardResponse{}, fmt.Errorf("Failed to collect charge via credit card (%s): %s", id, err)
	}

	if v, ok := endpoint.ResponseData.(*CollectChargeViaCreditCardResponse); ok {
		return *v, nil
	}

	return CollectChargeViaCreditCardResponse{}, fmt.Errorf("Data retrieved was not 'CollectChargeViaCreditCardResponse'")
}

// CollectChargeViaBankAccount collects charge via bank account
// https://www.zoho.com/subscriptions/api/v1/#Invoices_Collect_charge_via_bank_account
func (s *API) CollectChargeViaBankAccount(id string, request CollectChargeViaBankAccountRequest) (data CollectChargeViaBankAccountResponse, err error) {
	endpoint := zoho.Endpoint{
		Name:         "invoices",
		URL:          fmt.Sprintf("https://subscriptions.zoho.%s/api/v1/invoices/%s/collect", s.ZohoTLD, id),
		Method:       zoho.HTTPPost,
		ResponseData: &CollectChargeViaBankAccountResponse{},
		RequestBody:  request,
		Headers: map[string]string{
			ZohoSubscriptionsEndpointHeader: s.OrganizationID,
		},
	}

	err = s.Zoho.HTTPRequest(&endpoint)
	if err != nil {
		return CollectChargeViaBankAccountResponse{}, fmt.Errorf("Failed to collect charge via bank account (%s): %s", id, err)
	}

	if v, ok := endpoint.ResponseData.(*CollectChargeViaBankAccountResponse); ok {
		return *v, nil
	}

	return CollectChargeViaBankAccountResponse{}, fmt.Errorf("Data retrieved was not 'CollectChargeViaBankAccountResponse'")
}

type CollectChargeViaBankAccountRequest struct {
	AccountID string `json:"account_id"`
}

type CollectChargeViaBankAccountResponse struct {
	Code    int64                        `json:"code"`
	Message string                       `json:"message"`
	Invoice CollectChargeInvoiceResponse `json:"invoice,omitempty"`
	Payment struct {
		PaymentID       string  `json:"payment_id"`
		PaymentMode     string  `json:"payment_mode"`
		Amount          float64 `json:"amount"`
		AmountRefunded  float64 `json:"amount_refunded"`
		BankCharges     float64 `json:"bank_charges"`
		Date            string  `json:"date"`
		Status          string  `json:"status"`
		ReferenceNumber string  `json:"reference_number"`
		DueDate         string  `json:"due_date"`
		AmountDue       float64 `json:"amount_due"`
		Description     string  `json:"description"`
		CustomerID      string  `json:"customer_id"`
		CustomerName    string  `json:"customer_name"`
		Email           string  `json:"email"`
		Autotransaction struct {
			AutotransactionID    string `json:"autotransaction_id"`
			PaymentGateway       string `json:"payment_gateway"`
			GatewayTransactionID string `json:"gateway_transaction_id"`
			GatewayErrorMessage  string `json:"gateway_error_message"`
			AccountID            string `json:"account_id"`
		} `json:"autotransaction"`
		Invoices []struct {
			InvoiceID     string  `json:"invoice_id"`
			InvoiceNumber string  `json:"invoice_number"`
			Date          string  `json:"date"`
			InvoiceAmount float64 `json:"invoice_amount"`
			AmountApplied float64 `json:"amount_applied"`
			BalanceAmount float64 `json:"balance_amount"`
		} `json:"invoices"`
		CurrencyCode   string        `json:"currency_code"`
		CurrencySymbol string        `json:"currency_symbol"`
		CustomFields   []CustomField `json:"custom_fields"`
		CreatedTime    string        `json:"created_time"`
		UpdatedTime    string        `json:"updated_time"`
	} `json:"payment,omitempty"`
}

type CollectChargeViaCreditCardRequest struct {
	CardID string `json:"card_id"`
}

type CollectChargeViaCreditCardResponse struct {
	Code    int64                        `json:"code"`
	Message string                       `json:"message"`
	Invoice CollectChargeInvoiceResponse `json:"invoice,omitempty"`
	Payment struct {
		PaymentID       string  `json:"payment_id"`
		PaymentMode     string  `json:"payment_mode"`
		Amount          float64 `json:"amount"`
		AmountRefunded  float64 `json:"amount_refunded"`
		BankCharges     float64 `json:"bank_charges"`
		Date            string  `json:"date"`
		Status          string  `json:"status"`
		ReferenceNumber string  `json:"reference_number"`
		Description     string  `json:"description"`
		CustomerID      string  `json:"customer_id"`
		CustomerName    string  `json:"customer_name"`
		Email           string  `json:"email"`
		Autotransaction struct {
			AutotransactionID    string `json:"autotransaction_id"`
			PaymentGateway       string `json:"payment_gateway"`
			GatewayTransactionID string `json:"gateway_transaction_id"`
			GatewayErrorMessage  string `json:"gateway_error_message"`
			CardID               string `json:"card_id"`
			LastFourDigits       string `json:"last_four_digits"`
			ExpiryMonth          int64  `json:"expiry_month"`
			ExpiryYear           int64  `json:"expiry_year"`
		} `json:"autotransaction"`
		Invoices []struct {
			InvoiceID     string  `json:"invoice_id"`
			InvoiceNumber string  `json:"invoice_number"`
			Date          string  `json:"date"`
			InvoiceAmount float64 `json:"invoice_amount"`
			AmountApplied float64 `json:"amount_applied"`
			BalanceAmount float64 `json:"balance_amount"`
		} `json:"invoices"`
		CurrencyCode   string        `json:"currency_code"`
		CurrencySymbol string        `json:"currency_symbol"`
		CustomFields   []CustomField `json:"custom_fields"`
		CreatedTime    string        `json:"created_time"`
		UpdatedTime    string        `json:"updated_time"`
	} `json:"payment,omitempty"`
}

type CollectChargeInvoiceResponse struct {
	AchPaymentInitiated     bool    `json:"ach_payment_initiated"`
	Adjustment              float64 `json:"adjustment"`
	AdjustmentDescription   string  `json:"adjustment_description"`
	AllowPartialPayments    bool    `json:"allow_partial_payments"`
	ApproverID              string  `json:"approver_id"`
	AutoRemindersConfigured bool    `json:"auto_reminders_configured"`
	Balance                 float64 `json:"balance"`
	BcyAdjustment           float64 `json:"bcy_adjustment"`
	BcyDiscountTotal        float64 `json:"bcy_discount_total"`
	BcyShippingCharge       float64 `json:"bcy_shipping_charge"`
	BcySubTotal             float64 `json:"bcy_sub_total"`
	BcyTaxTotal             float64 `json:"bcy_tax_total"`
	BcyTotal                float64 `json:"bcy_total"`
	BillingAddress          struct {
		Address   string `json:"address"`
		Attention string `json:"attention"`
		City      string `json:"city"`
		Country   string `json:"country"`
		Fax       string `json:"fax"`
		Phone     string `json:"phone"`
		State     string `json:"state"`
		Street    string `json:"street"`
		Street2   string `json:"street2"`
		Zip       string `json:"zip"`
	} `json:"billing_address"`
	CanEditItems       bool   `json:"can_edit_items"`
	CanSendInMail      bool   `json:"can_send_in_mail"`
	CanSendInvoiceSms  bool   `json:"can_send_invoice_sms"`
	CanSkipPaymentInfo bool   `json:"can_skip_payment_info"`
	ClientViewedTime   string `json:"client_viewed_time"`
	Contactpersons     []struct {
		ContactpersonID string `json:"contactperson_id"`
		Email           string `json:"email"`
		Mobile          string `json:"mobile"`
		Phone           string `json:"phone"`
		ZcrmContactID   string `json:"zcrm_contact_id"`
	} `json:"contactpersons"`
	Coupons         []Coupon `json:"coupons"`
	CreatedByID     string   `json:"created_by_id"`
	CreatedDate     string   `json:"created_date"`
	CreatedTime     string   `json:"created_time"`
	Credits         []Credit `json:"credits"`
	CreditsApplied  float64  `json:"credits_applied"`
	CurrencyCode    string   `json:"currency_code"`
	CurrencyID      string   `json:"currency_id"`
	CurrencySymbol  string   `json:"currency_symbol"`
	CustomFieldHash struct {
	} `json:"custom_field_hash"`
	CustomFields            []CustomField `json:"custom_fields"`
	CustomerCustomFieldHash struct {
	} `json:"customer_custom_field_hash"`
	CustomerCustomFields        []CustomField `json:"customer_custom_fields"`
	CustomerID                  string        `json:"customer_id"`
	CustomerName                string        `json:"customer_name"`
	Date                        string        `json:"date"`
	DiscountPercent             float64       `json:"discount_percent"`
	DiscountTotal               float64       `json:"discount_total"`
	Documents                   []interface{} `json:"documents"`
	DueDate                     string        `json:"due_date"`
	Email                       string        `json:"email"`
	ExchangeRate                float64       `json:"exchange_rate"`
	InprocessTransactionPresent bool          `json:"inprocess_transaction_present"`
	InvoiceDate                 string        `json:"invoice_date"`
	InvoiceID                   string        `json:"invoice_id"`
	InvoiceItems                []struct {
		AccountID        string        `json:"account_id"`
		AccountName      string        `json:"account_name"`
		Code             string        `json:"code"`
		Description      string        `json:"description"`
		DiscountAmount   float64       `json:"discount_amount"`
		ItemCustomFields []CustomField `json:"item_custom_fields"`
		ItemID           string        `json:"item_id"`
		ItemTotal        float64       `json:"item_total"`
		Name             string        `json:"name"`
		Price            float64       `json:"price"`
		ProductID        string        `json:"product_id"`
		ProductType      string        `json:"product_type"`
		Quantity         int64         `json:"quantity"`
		Tags             []Tag         `json:"tags"`
		TaxID            string        `json:"tax_id"`
		TaxName          string        `json:"tax_name"`
		TaxPercentage    float64       `json:"tax_percentage"`
		TaxType          string        `json:"tax_type"`
		Unit             string        `json:"unit"`
	} `json:"invoice_items"`
	InvoiceNumber          string `json:"invoice_number"`
	InvoiceURL             string `json:"invoice_url"`
	IsInclusiveTax         bool   `json:"is_inclusive_tax"`
	IsReverseChargeApplied bool   `json:"is_reverse_charge_applied"`
	IsViewedByClient       bool   `json:"is_viewed_by_client"`
	IsViewedInMail         bool   `json:"is_viewed_in_mail"`
	LastModifiedByID       string `json:"last_modified_by_id"`
	MailFirstViewedTime    string `json:"mail_first_viewed_time"`
	MailLastViewedTime     string `json:"mail_last_viewed_time"`
	Notes                  string `json:"notes"`
	Number                 string `json:"number"`
	PageWidth              string `json:"page_width"`
	PaymentExpectedDate    string `json:"payment_expected_date"`
	PaymentGateways        []struct {
		PaymentGateway string `json:"payment_gateway"`
	} `json:"payment_gateways"`
	PaymentMade            int64  `json:"payment_made"`
	PaymentReminderEnabled bool   `json:"payment_reminder_enabled"`
	PaymentTerms           int64  `json:"payment_terms"`
	PaymentTermsLabel      string `json:"payment_terms_label"`
	Payments               []struct {
		Amount               float64 `json:"amount"`
		AmountRefunded       float64 `json:"amount_refunded"`
		BankCharges          float64 `json:"bank_charges"`
		CardType             string  `json:"card_type"`
		Date                 string  `json:"date"`
		Description          string  `json:"description"`
		ExchangeRate         float64 `json:"exchange_rate"`
		GatewayTransactionID string  `json:"gateway_transaction_id"`
		InvoicePaymentID     string  `json:"invoice_payment_id"`
		LastFourDigits       string  `json:"last_four_digits"`
		PaymentID            string  `json:"payment_id"`
		PaymentMode          string  `json:"payment_mode"`
		ReferenceNumber      string  `json:"reference_number"`
		SettlementStatus     string  `json:"settlement_status"`
	} `json:"payments"`
	PricePrecision  int64  `json:"price_precision"`
	PricebookID     string `json:"pricebook_id"`
	ReferenceID     string `json:"reference_id"`
	ReferenceNumber string `json:"reference_number"`
	SalespersonID   string `json:"salesperson_id"`
	SalespersonName string `json:"salesperson_name"`
	ShippingAddress struct {
		Address   string `json:"address"`
		Attention string `json:"attention"`
		City      string `json:"city"`
		Country   string `json:"country"`
		Fax       string `json:"fax"`
		Phone     string `json:"phone"`
		State     string `json:"state"`
		Street    string `json:"street"`
		Street2   string `json:"street2"`
		Zip       string `json:"zip"`
	} `json:"shipping_address"`
	ShippingCharge                        float64       `json:"shipping_charge"`
	ShippingChargeExclusiveOfTax          float64       `json:"shipping_charge_exclusive_of_tax"`
	ShippingChargeExclusiveOfTaxFormatted string        `json:"shipping_charge_exclusive_of_tax_formatted"`
	ShippingChargeInclusiveOfTax          float64       `json:"shipping_charge_inclusive_of_tax"`
	ShippingChargeInclusiveOfTaxFormatted string        `json:"shipping_charge_inclusive_of_tax_formatted"`
	ShippingChargeTax                     string        `json:"shipping_charge_tax"`
	ShippingChargeTaxFormatted            string        `json:"shipping_charge_tax_formatted"`
	ShippingChargeTaxID                   string        `json:"shipping_charge_tax_id"`
	ShippingChargeTaxName                 string        `json:"shipping_charge_tax_name"`
	ShippingChargeTaxPercentage           string        `json:"shipping_charge_tax_percentage"`
	ShippingChargeTaxType                 string        `json:"shipping_charge_tax_type"`
	Status                                string        `json:"status"`
	StopReminderUntilPaymentExpectedDate  bool          `json:"stop_reminder_until_payment_expected_date"`
	SubTotal                              float64       `json:"sub_total"`
	SubmitterID                           string        `json:"submitter_id"`
	Subscriptions                         []interface{} `json:"subscriptions"`
	TaxRounding                           string        `json:"tax_rounding"`
	TaxTotal                              float64       `json:"tax_total"`
	Taxes                                 []interface{} `json:"taxes"`
	TemplateID                            string        `json:"template_id"`
	TemplateName                          string        `json:"template_name"`
	TemplateType                          string        `json:"template_type"`
	Terms                                 string        `json:"terms"`
	Total                                 float64       `json:"total"`
	TransactionType                       string        `json:"transaction_type"`
	UnbilledChargesID                     string        `json:"unbilled_charges_id"`
	UnusedCreditsReceivableAmount         float64       `json:"unused_credits_receivable_amount"`
	UpdatedTime                           string        `json:"updated_time"`
	VatTreatment                          string        `json:"vat_treatment"`
	WriteOffAmount                        float64       `json:"write_off_amount"`
	ZcrmPotentialID                       string        `json:"zcrm_potential_id"`
}

type AddItemsRequest struct {
	InvoiceItems []InvoiceItemRequest `json:"invoice_items,omitempty"`
}

type InvoiceItemRequest struct {
	Code           string  `json:"code,omitempty"`
	ProductID      string  `json:"product_id,omitempty"`
	Name           string  `json:"name,omitempty"`
	Description    string  `json:"description,omitempty"`
	Price          float64 `json:"price,omitempty"`
	Quantity       float64 `json:"quantity,omitempty"`
	TaxID          string  `json:"tax_id,omitempty"`
	TaxExemptionID string  `json:"tax_exemption_id,omitempty"`
}

type AddItemsResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Invoice struct {
		InvoiceID           string `json:"invoice_id"`
		Number              string `json:"number"`
		Status              string `json:"status"`
		InvoiceDate         string `json:"invoice_date"`
		DueDate             string `json:"due_date"`
		PaymentExpectedDate string `json:"payment_expected_date"`
		AchPaymentInitiated bool   `json:"ach_payment_initiated"`
		TransactionType     string `json:"transaction_type"`
		CustomerID          string `json:"customer_id"`
		CustomerName        string `json:"customer_name"`
		Email               string `json:"email"`
		InvoiceItems        []struct {
			ItemID           string        `json:"item_id"`
			Name             string        `json:"name"`
			Description      string        `json:"description"`
			Tags             []Tag         `json:"tags"`
			ItemCustomFields []CustomField `json:"item_custom_fields"`
			Code             string        `json:"code"`
			Price            float64       `json:"price"`
			Quantity         float64       `json:"quantity"`
			DiscountAmount   float64       `json:"discount_amount"`
			ItemTotal        float64       `json:"item_total"`
			TaxID            string        `json:"tax_id"`
			ProductType      string        `json:"product_type"`
			HsnOrSac         string        `json:"hsn_or_sac"`
			TaxExemptionID   string        `json:"tax_exemption_id"`
			TaxExemptionCode string        `json:"tax_exemption_code"`
		} `json:"invoice_items"`
		Coupons []struct {
			CouponCode     string  `json:"coupon_code"`
			CouponName     string  `json:"coupon_name"`
			DiscountAmount float64 `json:"discount_amount"`
		} `json:"coupons"`
		Credits []struct {
			CreditnoteID      string  `json:"creditnote_id"`
			CreditnotesNumber string  `json:"creditnotes_number"`
			CreditedDate      string  `json:"credited_date"`
			CreditedAmount    float64 `json:"credited_amount"`
		} `json:"credits"`
		Total          float64 `json:"total"`
		PaymentMade    float64 `json:"payment_made"`
		Balance        float64 `json:"balance"`
		CreditsApplied float64 `json:"credits_applied"`
		WriteOffAmount float64 `json:"write_off_amount"`
		Payments       []struct {
			PaymentID            string  `json:"payment_id"`
			PaymentMode          string  `json:"payment_mode"`
			InvoicePaymentID     string  `json:"invoice_payment_id"`
			GatewayTransactionID string  `json:"gateway_transaction_id"`
			Description          string  `json:"description"`
			Date                 string  `json:"date"`
			ReferenceNumber      string  `json:"reference_number"`
			Amount               float64 `json:"amount"`
			BankCharges          float64 `json:"bank_charges"`
			ExchangeRate         float64 `json:"exchange_rate"`
		} `json:"payments"`
		CurrencyCode    string  `json:"currency_code"`
		CurrencySymbol  string  `json:"currency_symbol"`
		CreatedTime     string  `json:"created_time"`
		UpdatedTime     string  `json:"updated_time"`
		SalespersonID   string  `json:"salesperson_id"`
		SalespersonName string  `json:"salesperson_name"`
		InvoiceURL      string  `json:"invoice_url"`
		BillingAddress  Address `json:"billing_address"`
		ShippingAddress Address `json:"shipping_address"`
		Comments        []struct {
			CommentID       string `json:"comment_id"`
			Description     string `json:"description"`
			CommentedByID   string `json:"commented_by_id"`
			CommentedBy     string `json:"commented_by"`
			CommentType     string `json:"comment_type"`
			Time            string `json:"time"`
			OperationType   string `json:"operation_type"`
			TransactionID   string `json:"transaction_id"`
			TransactionType string `json:"transaction_type"`
		} `json:"comments"`
		CustomFields []CustomField `json:"custom_fields"`
	} `json:"invoice"`
}

type EmailInvoiceRequest struct {
	ToMailIds []string `json:"to_mail_ids"`
	CcMailIds []string `json:"cc_mail_ids"`
	Subject   string   `json:"subject"`
	Body      string   `json:"body"`
}

type EmailInvoiceResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

type AttachmentRequest struct {
	CanSendInEmail bool `json:"can_send_in_mail"`
}

type AttachementResponse struct {
	Code      int64  `json:"code"`
	Message   string `json:"message"`
	Documents []struct {
		FileName          string `json:"file_name"`
		FileType          string `json:"file_type"`
		FileSize          int64  `json:"file_size"`
		FileSizeFormatted string `json:"file_size_formatted"`
		DocumentID        string `json:"document_id"`
		AttachmentOrder   int64  `json:"attachment_order"`
	} `json:"documents"`
}

type InvoicesResponse struct {
	Invoices []Invoice `json:"invoices"`
	Code     int64     `json:"code"`
	Message  string    `json:"message"`
}

type InvoiceResponse struct {
	Invoice Invoice `json:"invoice"`
	Code    int64   `json:"code"`
	Message string  `json:"message"`
}

type Invoice struct {
	InvoiceID            string        `json:"invoice_id,omitempty"`
	Number               string        `json:"number,omitempty"`
	Status               string        `json:"status,omitempty"`
	InvoiceDate          string        `json:"invoice_date,omitempty"`
	DueDate              string        `json:"due_date,omitempty"`
	CustomerID           string        `json:"customer_id,omitempty"`
	CustomerName         string        `json:"customer_name,omitempty"`
	Email                string        `json:"email,omitempty"`
	Balance              float64       `json:"balance,omitempty"`
	Total                float64       `json:"total,omitempty"`
	PaymentMade          float64       `json:"payment_made,omitempty"`
	CreditsApplied       float64       `json:"credits_applied,omitempty"`
	WriteOffAmount       float64       `json:"write_off_amount,omitempty"`
	CurrencyCode         string        `json:"currency_code,omitempty"`
	CurrencySymbol       string        `json:"currency_symbol,omitempty"`
	HasAttachment        bool          `json:"has_attachment,omitempty"`
	CreatedTime          string        `json:"created_time,omitempty"`
	UpdatedTime          string        `json:"updated_time,omitempty"`
	SalespersonID        string        `json:"salesperson_id,omitempty"`
	SalespersonName      string        `json:"salesperson_name,omitempty"`
	InvoiceUrl           string        `json:"invoice_url,omitempty"`
	PaymentExpectedDate  string        `json:"payment_expected_date,omitempty"`
	ArchPaymentInitiated interface{}   `json:"ach_payment_initiated,omitempty"` // per documentation this field should be bool, but received empty string
	TransactionType      string        `json:"transaction_type,omitempty"`
	InvoiceItems         []InvoiceItem `json:"invoice_items,omitempty"`
	Coupons              []Coupon      `json:"coupons,omitempty"`
	Credits              []Credit      `json:"credits,omitempty"`
	Payments             []Payment     `json:"payments,omitempty"`
	BillingAddress       Address       `json:"billing_address,omitempty"`
	ShippingAddress      Address       `json:"shipping_address,omitempty"`
	Comments             []Comment     `json:"comments,omitempty"`
	CustomFields         []CustomField `json:"custom_fields,omitempty"`
	CanSendInEmail       bool          `json:"can_send_in_mail,omitempty"`
	Documents            []Document    `json:"documents,omitempty"`
}

type InvoiceItem struct {
	ItemID           string        `json:"item_id,omitempty"`
	Name             string        `json:"name,omitempty"`
	Description      string        `json:"description,omitempty"`
	Code             string        `json:"code,omitempty"`
	Tags             []Tag         `json:"tags,omitempty"`
	ItemCustomFields []CustomField `json:"item_custom_fields,omitempty"`
	Price            float64       `json:"price,omitempty"`
	Quantity         float64       `json:"quantity,omitempty"`
	DiscountAmount   float64       `json:"discount_amount,omitempty"`
	ItemTotal        float64       `json:"item_total,omitempty"`
	TaxID            string        `json:"tax_id,omitempty"`
	TaxExemptionID   string        `json:"tax_exemption_id,omitempty"`
	TaxExemptionCode string        `json:"tax_exemption_code,omitempty"`
}

type Coupon struct {
	CouponCode     string  `json:"coupon_code,omitempty"`
	CouponName     string  `json:"coupon_name,omitempty"`
	DiscountAmount float64 `json:"discount_amount,omitempty"`
}

type Credit struct {
	CreditnoteID      string  `json:"creditnote_id,omitempty"`
	CreditnotesNumber string  `json:"creditnotes_number,omitempty"`
	CreditedDate      string  `json:"credited_date,omitempty"`
	CreditedAmount    float64 `json:"credited_amount,omitempty"`
}

type Payment struct {
	PaymentID            string  `json:"payment_id,omitempty"`
	PaymentMode          string  `json:"payment_mode,omitempty"`
	InvoicePaymentID     string  `json:"invoice_payment_id,omitempty"`
	AmountRefunded       float64 `json:"amount_refunded,omitempty"`
	GatewayTransactionID string  `json:"gateway_transaction_id,omitempty"`
	Description          string  `json:"description,omitempty"`
	Date                 string  `json:"date,omitempty"`
	ReferenceNumber      string  `json:"reference_number,omitempty"`
	Amount               float64 `json:"amount,omitempty"`
	BankCharges          float64 `json:"bank_charges,omitempty"`
	ExchangeRate         float64 `json:"exchange_rate,omitempty"`
}

type Comment struct {
	CommentID       string `json:"comment_id,omitempty"`
	Description     string `json:"description,omitempty"`
	CommentedByID   string `json:"commented_by_id,omitempty"`
	CommentedBy     string `json:"commented_by,omitempty"`
	CommentType     string `json:"comment_type,omitempty"`
	Time            string `json:"time,omitempty"`
	OperationType   string `json:"operation_type,omitempty"`
	TransactionID   string `json:"transaction_id,omitempty"`
	TransactionType string `json:"transaction_type,omitempty"`
}

type Document struct {
	FileName          string `json:"file_name,omitempty"`
	FileType          string `json:"file_type,omitempty"`
	FileSize          int64  `json:"file_size,omitempty"`
	FileSizeFormatted string `json:"file_size_formatted,omitempty"`
	DocumentID        string `json:"document_id,omitempty"`
	AttachmentOrder   int64  `json:"attachment_order,omitempty"`
}

type Tag struct {
	TagID       string `json:"tag_id,omitempty"`
	TagOptionID string `json:"tag_option_id,omitempty"`
}
