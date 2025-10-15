package utils

import "github.com/hayden-erickson/ai-evaluation/internal/models"

// UniqueIntSlice removes duplicate integers from a slice
func UniqueIntSlice(input []int) ([]int, error) {
	keys := make(map[int]bool)
	result := []int{}
	for _, val := range input {
		if _, exists := keys[val]; !exists {
			keys[val] = true
			result = append(result, val)
		}
	}
	return result, nil
}

// ConvertToStringSlice converts GateAccessCodes to a slice of access code strings
func ConvertToStringSlice(gacs models.GateAccessCodes) []string {
	result := make([]string, len(gacs))
	for i, gac := range gacs {
		result[i] = gac.AccessCode
	}
	return result
}
