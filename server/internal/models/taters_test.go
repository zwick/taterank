package models

import (
	"testing"

	"taterank.com/internal/assert"
	"taterank.com/internal/database"
)

func TestTaterModelGet(t *testing.T) {
	db, err := database.GetTestDynamoDBClient(t, database.TestConfigOptions{})

	if err != nil {
		t.Errorf("Error getting DynamoDB client: %v", err)
	}

	taterModel := TaterModel{DB: db}

	// Get the Taters
	taters, err := taterModel.Get()

	if err != nil {
		t.Errorf("Error getting taters: %v", err)
	}

	assert.Equals(t, len(taters), 10)
}

func TestTaterModelGetByID(t *testing.T) {
	db, err := database.GetTestDynamoDBClient(t, database.TestConfigOptions{})

	if err != nil {
		t.Errorf("Error getting DynamoDB client: %v", err)
	}

	taterModel := TaterModel{DB: db}

	tests := []struct {
		id          string
		expected    bool
		Name        string
		Description string
	}{
		{"46db56c79761", true, "Curly Fries", "Curly fries are a type of French fry characterized by their helical shape, which is formed by cutting the potato in a spiral shape before frying."},
		{"abc1234", false, "", ""},
	}

	for _, test := range tests {
		tater, err := taterModel.GetByID(test.id)

		if err != nil {
			t.Errorf("Error getting tater: %v", err)
		}

		if test.expected {
			assert.Equals(t, tater.ID, test.id)
			assert.Equals(t, tater.Name, test.Name)
			assert.Equals(t, tater.Description, test.Description)
		} else {
			if tater != nil {
				t.Errorf("Expected tater to be nil, got: %v", tater)
			}
		}

	}
}
