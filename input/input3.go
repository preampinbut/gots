package test

type PetOwnership struct {
	Name string `json:"name"`
	Pet  []Pet  `json:"pet"`
}
