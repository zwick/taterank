package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"taterank.com/internal/database"
)

func TestTaterModelList(t *testing.T) {
	t.Run("lists all taters", func(t *testing.T) {
		testBag := setup(t)

		taters, err := testBag.model.List()

		assert.NoError(t, err)
		assert.Greater(t, len(taters), 0)
	})

	t.Run("returns error with bad config", func(t *testing.T) {
		db, err := database.GetTestDynamoDBClient(t, database.TestConfigOptions{Endpoint: "bad"})
		assert.NoError(t, err)

		taterModel := TaterModel{DB: db}

		_, err = taterModel.List()
		assert.Error(t, err)
	})
}

func TestTaterModelGet(t *testing.T) {

	t.Run("gets tater by ID", func(t *testing.T) {
		testBag := setup(t)

		name := strPtr("Test Tater")
		description := strPtr("This is a test tater")

		input := TaterFields{
			Name:        name,
			Description: description,
		}

		id, err := testBag.model.Create(input)

		assert.NoError(t, err)

		tater, err := testBag.model.Get(*id)

		assert.NoError(t, err)
		assert.Equal(t, tater.ID, *id)
		assert.Equal(t, tater.Name, name)
		assert.Equal(t, tater.Description, description)
	})

	t.Run("returns error for non-existent tater", func(t *testing.T) {
		testBag := setup(t)

		result, err := testBag.model.Get("id-that-definitely-does-not-exist")

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("returns error with bad config", func(t *testing.T) {
		db, err := database.GetTestDynamoDBClient(t, database.TestConfigOptions{Endpoint: "bad"})
		assert.NoError(t, err)

		taterModel := TaterModel{DB: db}

		_, err = taterModel.Get("id")
		assert.Error(t, err)
	})

}

func TestTaterModelUpdate(t *testing.T) {
	t.Run("updates all fields", func(t *testing.T) {
		testBag := setup(t)

		originalTater := testBag.createTestTater("Test Tater", "This is a test tater")

		fields := TaterFields{
			Name:        strPtr("Amazing Potatoes"),
			Description: strPtr("Are with you always!"),
		}

		assert.NotEqual(t, originalTater, fields)

		err := testBag.model.Update(originalTater.ID, fields)
		assert.NoError(t, err)

		updatedTater, err := testBag.model.Get(originalTater.ID)
		assert.NoError(t, err)

		assert.Equal(t, originalTater.ID, updatedTater.ID)
		assert.Equal(t, fields, updatedTater.TaterFields)
	})

	t.Run("updates some fields", func(t *testing.T) {
		testBag := setup(t)

		originalTater := testBag.createTestTater("Test Tater", "This is a test tater")

		name := "My New Name"

		fields := TaterFields{
			Name: strPtr(name),
		}

		err := testBag.model.Update(originalTater.ID, fields)
		assert.NoError(t, err)

		updatedTater, err := testBag.model.Get(originalTater.ID)
		assert.NoError(t, err)

		assert.Equal(t, originalTater.ID, updatedTater.ID)
		assert.Equal(t, originalTater.Description, updatedTater.Description)
		assert.NotEqual(t, originalTater.Name, updatedTater.Name)
		assert.Equal(t, fields.Name, updatedTater.Name)
	})

	t.Run("handles empty fields", func(t *testing.T) {

		testBag := setup(t)
		originalTater := testBag.createTestTater("Test Tatert", "This is a test tater")

		fields := TaterFields{
			Name: strPtr(""),
		}

		err := testBag.model.Update(originalTater.ID, fields)

		assert.NoError(t, err)

		updatedTater, _ := testBag.model.Get(originalTater.ID)

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

func TestTaterModelCreate(t *testing.T) {
	t.Run("creates tater with all fields", func(t *testing.T) {
		testBag := setup(t)

		fields := TaterFields{
			Name:        strPtr("Amazing Potatoes"),
			Description: strPtr("Are with you always!"),
		}

		id, err := testBag.model.Create(fields)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		tater, err := testBag.model.Get(*id)
		assert.NoError(t, err)

		assert.Equal(t, &tater.ID, id)
		assert.Equal(t, tater.TaterFields, fields)
	})

	t.Run("creates tater with some fields", func(t *testing.T) {
		testBag := setup(t)

		fieldsWithName := TaterFields{
			Name: strPtr("Amazing Potatoes"),
		}

		id, err := testBag.model.Create(fieldsWithName)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		tater, err := testBag.model.Get(*id)

		assert.NoError(t, err)
		assert.Equal(t, &tater.ID, id)
		assert.Equal(t, tater.TaterFields, fieldsWithName)

		fieldsWithDescription := TaterFields{
			Description: strPtr("Neat Potato Description!"),
		}

		id, err = testBag.model.Create(fieldsWithDescription)

		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		tater, err = testBag.model.Get(*id)

		assert.NoError(t, err)
		assert.Equal(t, &tater.ID, id)
		assert.Equal(t, tater.TaterFields, fieldsWithDescription)
	})

	t.Run("returns error with bad config", func(t *testing.T) {
		db, err := database.GetTestDynamoDBClient(t, database.TestConfigOptions{Endpoint: "bad"})

		if err != nil {
			t.Errorf("Error getting DynamoDB client: %v", err)
		}

		taterModel := TaterModel{DB: db}

		var fields TaterFields

		_, err = taterModel.Create(fields)

		assert.Error(t, err)
	})
}

func TestSanitizer(t *testing.T) {
	id := TaterPreparationsPrefix + "46db56c79761"

	tater := Tater{
		ID: id,
		TaterFields: TaterFields{
			Name:        strPtr("Test Name"),
			Description: strPtr("Test Description"),
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
				Name:        strPtr(name),
				Description: strPtr(description),
			},
		},
		{
			ID: TaterPreparationsPrefix + "52kd01kdl2ds",
			TaterFields: TaterFields{
				Name:        strPtr(name),
				Description: strPtr(description),
			},
		},
	}

	collectionSanitizer(taters)

	assert.Equal(t, taters[0].ID, "46db56c79761")
	assert.Equal(t, taters[1].ID, "52kd01kdl2ds")
}

type TestBag struct {
	model TaterModel
	t     *testing.T
}

func setup(t *testing.T) TestBag {
	db, err := database.GetTestDynamoDBClient(t, database.TestConfigOptions{})
	assert.NoError(t, err)

	taterModel := TaterModel{DB: db}

	bag := TestBag{
		model: taterModel,
		t:     t,
	}

	return bag
}

func strPtr(s string) *string {
	return &s
}

func (b *TestBag) createTestTater(name string, description string) *Tater {
	id, err := b.model.Create(TaterFields{
		Name:        strPtr(name),
		Description: strPtr(description),
	})

	assert.NoError(b.t, err)

	tater, err := b.model.Get(*id)

	assert.NoError(b.t, err)

	return tater
}
