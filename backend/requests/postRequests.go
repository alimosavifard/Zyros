package requests

import (
)

type PostRequest struct {
    Title    string `json:"title" validate:"required,min=3"`
    Content  string `json:"content" validate:"required,min=10"`
    Type     string `json:"type" validate:"required,oneof=post article"`
    Lang     string `json:"lang" validate:"required,oneof=fa en"`
    ImageUrl string `json:"imageUrl" validate:"omitempty,url"` // اختیاری
}