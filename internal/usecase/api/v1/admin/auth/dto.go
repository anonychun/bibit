package auth

type SignInRequest struct {
	IpAddress    string `json:"-"`
	UserAgent    string `json:"-"`
	EmailAddress string `json:"emailAddress" validate:"required|email" field:"emailAddress" label:"Email address"`
	Password     string `json:"password" validate:"required|minLen:8" field:"password" label:"Password"`
}

type SignInResponse struct {
	Token string
}

type SignOutRequest struct {
	Token string
}

type MeResponse struct {
	Admin struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"admin"`
}
