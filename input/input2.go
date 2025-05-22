package test

type PersonalInfo struct {
	Name  string `json:"name"`
	Tel   string `json:"tel,omitempty"`
	Email string `json:"email,omitempty"`
}
