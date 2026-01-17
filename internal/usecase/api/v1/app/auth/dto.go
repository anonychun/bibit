package auth

type SignUpRequest struct {
	IpAddress    string `json:"-"`
	UserAgent    string `json:"-"`
	Name         string `json:"name" validate:"required" field:"name" label:"Name"`
	EmailAddress string `json:"emailAddress" validate:"required|email" field:"emailAddress" label:"Email address"`
	Password     string `json:"password" validate:"required|minLen:8" field:"password" label:"Password"`
}

type SignUpResponse struct {
	Token string
}

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
	User struct {
		Id           string `json:"id"`
		Name         string `json:"name"`
		EmailAddress string `json:"emailAddress"`
	} `json:"user"`
}
