package types

import "image/color"

type ImageSet struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	AverageColors []*color.RGBA `json:"average_colors"`
}
