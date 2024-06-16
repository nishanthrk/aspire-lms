package dto

import "github.com/guregu/null"

type ApplicationCreateRequest struct {
	User            UserObject            `json:"user"`
	LoanApplication LoanApplicationObject `json:"loan_application"`
}

type ApplicationCreateResponse struct {
	Status int `json:"status"`
	Data   struct {
		ApplicationID string `json:"application_id"`
		UserID        string `json:"user_id"`
	} `json:"data"`
	Message string `json:"message"`
}

type ApplicationApproveRequest struct {
	ApplicationID  string  `json:"-" validate:"required"`
	ApprovedAmount float64 `json:"approved_amount" validate:"required"`
	Override       bool    `json:"override"`
}

type ApplicationApproveResponse struct {
	Data    ApproveObject `json:"data"`
	Message string        `json:"message"`
	Status  int           `json:"status"`
}

type ApproveObject struct {
	ApplicationID string `json:"application_id"`
}

type RepaymentRequest struct {
	ApplicationID string  `json:"application_id" validate:"required"`
	PaymentAmount float64 `json:"payment_amount" validate:"required"`
}

type RepaymentResponse struct {
	Data struct {
		PaymentId string `json:"payment_id"`
	} `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type ApplicationDetailsResponse struct {
	Status int `json:"status"`
	Data   struct {
		LoanApplicationObject
		Repayment []Repayment `json:"repayments"`
	} `json:"data"`
}

type ApplicationListResponse struct {
	Status  int                      `json:"status"`
	Message string                   `json:"message"`
	Data    *[]LoanApplicationObject `json:"data"`
}

type UserObject struct {
	UserName     string    `json:"user_name" validate:"required"`
	UserEmail    string    `json:"user_email" validate:"required"`
	UserPassword string    `json:"user_password"`
	MobileNumber string    `json:"mobile_number" validate:"required"`
	CountryCode  string    `json:"country_code" validate:"validateCountryCode"`
	Kyc          KycObject `json:"kyc"`
}

type KycObject struct {
	KycType   string `json:"kyc_type" validate:"required,oneof=PAN SSN AADHAAR"`
	KycNumber string `json:"kyc_number" validate:"required"`
}

type LoanApplicationObject struct {
	ApplicationId     string  `json:"application_id,omitempty"`
	ApplicationStatus string  `json:"application_status,omitempty"`
	LoanAmount        float64 `json:"loan_amount" validate:"required"`
	CurrencyCode      string  `json:"currency_code" validate:"validateCurrencyCode"`
	InterestRate      float64 `json:"interest_rate"`
	LoanTerm          int     `json:"loan_term" validate:"required"`
	LoanTermUnit      string  `json:"loan_term_unit" validate:"required,oneof=WEEKLY MONTHLY"`
	Income            float64 `json:"income" validate:"required"`
	CreditScore       int     `json:"credit_score" validate:"required"`
	ExistingDebts     float64 `json:"existing_debts" validate:"required"`
	CountryCode       string  `json:"country_code" validate:"validateCountryCode"`
}

type Repayment struct {
	InstallmentNumber  int         `json:"installment_number"`
	PaymentDate        string      `json:"payment_date"`
	CurrencyCode       string      `json:"currency"`
	Principle          string      `json:"principle"`
	Interest           string      `json:"interest"`
	Emi                string      `json:"emi"`
	AmountPaid         null.String `json:"amount_paid"`
	OutstandingBalance string      `json:"outstanding_balance,omitempty"`
	RepaymentStatus    string      `json:"repayment_status"`
}
