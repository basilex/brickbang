package transfer

// CountryCreateRequest
type CountryCreateRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=255"`
	Iso2    string `json:"iso2" validate:"required,len=2"`
	Iso3    string `json:"iso3" validate:"required,len=3"`
	NumCode int16  `json:"num_code" validate:"required,min=1,max=999"`
}

// CountryUpdateRequest
type CountryUpdateRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=255"`
	Iso2    string `json:"iso2" validate:"required,len=2"`
	Iso3    string `json:"iso3" validate:"required,len=3"`
	NumCode int16  `json:"num_code" validate:"required,min=1,max=999"`
}

// CountryResponse
type CountryResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Iso2      string `json:"iso2"`
	Iso3      string `json:"iso3"`
	NumCode   int16  `json:"num_code"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
