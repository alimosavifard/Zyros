package requests

type ArticleRequest struct {
    Title    string `json:"title" validate:"required,min=3"`
    Content  string `json:"content" validate:"required,min=10"`
    Lang     string `json:"lang" validate:"required,oneof=fa en"`
    ImageUrl string `json:"imageUrl" validate:"omitempty,url"`
}

func (r *ArticleRequest) Validate() error {
	return ValidateStruct(r)
}