package transfer

// CurrencyCreateRequest
type CurrencyCreateRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=255"`
	Code    string `json:"code" validate:"required,len=3"`
	NumCode int16  `json:"num_code" validate:"required,min=1,max=999"`
	Symbol  string `json:"symbol" validate:"required,min=1,max=8"`
}

// CurrencyUpdateRequest
type CurrencyUpdateRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=255"`
	Code    string `json:"code" validate:"required,len=3"`
	NumCode int16  `json:"num_code" validate:"required,min=1,max=999"`
	Symbol  string `json:"symbol" validate:"required,min=1,max=8"`
}

// CurrencyResponse
type CurrencyResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	NumCode   int16  `json:"num_code"`
	Symbol    string `json:"symbol"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
