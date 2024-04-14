package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"taterank.com/internal/database"
)

func TestTaterModelList(t *testing.T) {
	t.Run("lists all taters", func(t *testing.T) {
		db, err := database.GetTestDynamoDBClient(t, database.TestConfigOptions{})

		if err != nil {
			t.Errorf("Error getting DynamoDB client: %v", err)
		}

		taterModel := TaterModel{DB: db}

		// Get the Taters
		taters, err := taterModel.List()

		if err != nil {
			t.Errorf("Error getting taters: %v", err)
		}

		assert.Greater(t, len(taters), 0)
	})

	t.Run("returns error with bad config", func(t *testing.T) {
		db, err := database.GetTestDynamoDBClient(t, database.TestConfigOptions{Endpoint: "bad"})

		if err != nil {
			t.Errorf("Error getting DynamoDB client: %v", err)
		}

		taterModel := TaterModel{DB: db}

		// Get the Taters
		_, err = taterModel.List()

		assert.Error(t, err)
	})
}

func TestTaterModelGet(t *testing.T) {

	t.Run("gets taters by ID", func(t *testing.T) {
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
			tater, err := taterModel.Get(test.id)

			if err != nil {
				t.Errorf("Error getting tater: %v", err)
			}

			if test.expected {
				assert.Equal(t, tater.ID, test.id)
				assert.Equal(t, *tater.Name, test.Name)
				assert.Equal(t, *tater.Description, test.Description)
			} else {
				if tater != nil {
					t.Errorf("Expected tater to be nil, got: %v", tater)
				}
			}

		}
	})

	t.Run("returns error with bad config", func(t *testing.T) {
		db, err := database.GetTestDynamoDBClient(t, database.TestConfigOptions{Endpoint: "bad"})

		if err != nil {
			t.Errorf("Error getting DynamoDB client: %v", err)
		}

		taterModel := TaterModel{DB: db}

		// Get the Taters
		_, err = taterModel.Get("id")

		assert.Error(t, err)
	})

}

func TestTaterModelUpdate(t *testing.T) {
	db, err := database.GetTestDynamoDBClient(t, database.TestConfigOptions{})

	if err != nil {
		t.Errorf("Error getting DynamoDB client: %v", err)
	}

	taterModel := TaterModel{DB: db}

	t.Run("updates all fields", func(t *testing.T) {
		taters, err := taterModel.List()

		if err != nil {
			t.Errorf("Error getting taters: %v", err)
		}

		originalTater := taters[1]

		updatedName := "Amazing Potatoes"
		updatedDescription := "Are with you always!"

		fields := TaterFields{
			Name:        &updatedName,
			Description: &updatedDescription,
		}

		assert.NotEqual(t, originalTater.TaterFields, fields)

		err = taterModel.Update(originalTater.ID, fields)

		if err != nil {
			t.Errorf("Error updating taters: %v", err)
		}

		updatedTater, _ := taterModel.Get(originalTater.ID)

		assert.Equal(t, originalTater.ID, updatedTater.ID)
		assert.Equal(t, fields, updatedTater.TaterFields)
	})

	t.Run("updates some fields", func(t *testing.T) {
		taters, err := taterModel.List()

		if err != nil {
			t.Errorf("Error getting taters: %v", err)
		}

		originalTater := taters[2]

		name := "My New Name"

		fields := TaterFields{
			Name: &name,
		}

		err = taterModel.Update(originalTater.ID, fields)

		if err != nil {
			t.Errorf("Error updating taters: %v", err)
		}

		updatedTater, _ := taterModel.Get(originalTater.ID)

		assert.Equal(t, originalTater.ID, updatedTater.ID)
		assert.Equal(t, originalTater.Description, updatedTater.Description)
		assert.NotEqual(t, originalTater.Name, updatedTater.Name)
		assert.Equal(t, fields.Name, updatedTater.Name)
	})

	t.Run("handles empty fields", func(t *testing.T) {
		taters, err := taterModel.List()

		if err != nil {
			t.Errorf("Error getting taters: %v", err)
		}

		originalTater := taters[3]

		name := ""

		fields := TaterFields{
			Name: &name,
		}

		err = taterModel.Update(originalTater.ID, fields)

		if err != nil {
			t.Errorf("Error updating taters: %v", err)
		}

		updatedTater, _ := taterModel.Get(originalTater.ID)

		assert.Equal(t, originalTater.ID, updatedTater.ID)
		assert.Equal(t, originalTater.Description, updatedTater.Description)
		assert.NotEqual(t, originalTater.Name, updatedTater.Name)
		assert.Equal(t, fields.Name, updatedTater.Name)
	})

	t.Run("returns error with bad config", func(t *testing.T) {
		db, err := database.GetTestDynamoDBClient(t, database.TestConfigOptions{Endpoint: "bad"})

		if err != nil {
			t.Errorf("Error getting DynamoDB client: %v", err)
		}

		taterModel := TaterModel{DB: db}

		var fields TaterFields

		err = taterModel.Update("id", fields)

		assert.Error(t, err)
	})
}

func TestSanitizer(t *testing.T) {
	id := TaterPreparationsPrefix + "46db56c79761"
	name := "Test Name"
	description := "Test Description"

	tater := Tater{
		ID: id,
		TaterFields: TaterFields{
			Name:        &name,
			Description: &description,
		},
	}

	sanitizer(&tater)

	assert.Equal(t, tater.ID, "46db56c79761")
}

func TestCollectionSanitizer(t *testing.T) {
	name := "Test Name"
	description := "Test Description"

	taters := []*Tater{
		{
			ID: TaterPreparationsPrefix + "46db56c79761",
			TaterFields: TaterFields{
				Name:        &name,
				Description: &description,
			},
		},
		{
			ID: TaterPreparationsPrefix + "52kd01kdl2ds",
			TaterFields: TaterFields{
				Name:        &name,
				Description: &description,
			},
		},
	}

	collectionSanitizer(taters)

	assert.Equal(t, taters[0].ID, "46db56c79761")
	assert.Equal(t, taters[1].ID, "52kd01kdl2ds")
}
