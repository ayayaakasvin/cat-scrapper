package domain

import catphotofetch "github.com/ayayaakasvin/cat-photo-fetch"

type Engine interface {
	SaveImage(*Job, *catphotofetch.Image) (string, error)
}
