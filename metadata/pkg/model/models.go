package model

import "movix/api/gen"

type Metadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Director    string `json:"director"`
}

func MetadataToProto(metadata *Metadata) *gen.Metadata {
	return &gen.Metadata{
		Id:          metadata.ID,
		Title:       metadata.Title,
		Description: metadata.Description,
		Director:    metadata.Director,
	}
}

func MetadataFromProto(metadata *gen.Metadata) *Metadata {
	return &Metadata{
		ID:          metadata.Id,
		Title:       metadata.Title,
		Description: metadata.Description,
		Director:    metadata.Director,
	}
}
