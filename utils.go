package gormet

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// isValid checks the validity of the entity based on the validation configuration.
func isValid[T any](entity *T, validable bool) error {
	if validable {
		validate := validator.New()
		return validate.Struct(entity)
	}

	return nil
}

func getPrimaryKeyFieldName(db *gorm.DB, model interface{}) (string, error) {
	// Initialize GORM statement
	stmt := &gorm.Statement{DB: db}

	// Parse the model to get *schema.Schema
	if err := stmt.Parse(model); err != nil {
		return "", err
	}

	// Loop over schema fields to find the primary key
	for _, field := range stmt.Schema.Fields {
		fmt.Printf("%s - %v\n", field.DBName, field.PrimaryKey)
		if field.PrimaryKey {
			return field.DBName, nil
		}
	}

	return "", fmt.Errorf("no primary key found")
}

// getCompositePrimaryKeyFieldNames retrieves the names of fields that are part of the composite primary key
// func getCompositePrimaryKeyFieldNames(model interface{}) ([]string, error) {
// 	var primaryKeyFieldNames []string

// 	// Get the type of the model
// 	modelType := reflect.TypeOf(model)

// 	// Iterate over the fields of the model using reflection
// 	for i := 0; i < modelType.NumField(); i++ {
// 		field := modelType.Field(i)

// 		// Check if the field is marked as primary key by Gorm
// 		if gormTag, exists := field.Tag.Lookup("gorm"); exists {
// 			if gormTagParts := parseGormTag(gormTag); gormTagParts["PRIMARY_KEY"] == "true" {
// 				primaryKeyFieldNames = append(primaryKeyFieldNames, field.Name)
// 			}
// 		}
// 	}

// 	// If no primary key fields found, return an error
// 	if len(primaryKeyFieldNames) == 0 {
// 		return nil, fmt.Errorf("No primary key fields found in the struct")
// 	}

// 	return primaryKeyFieldNames, nil
// }

// func parseGormTag(input string) map[string]string {
// 	// Converter para maiúsculas
// 	input = strings.ToUpper(input)

// 	// Substituir vírgulas por ponto e vírgula
// 	input = strings.Replace(input, ",", ";", -1)

// 	// Dividir a string com base no ponto e vírgula
// 	tags := strings.Split(input, ";")

// 	// Criar um mapa para armazenar as chaves e valores
// 	tagMap := make(map[string]string)

// 	for _, tag := range tags {
// 		var value string // Inicializar o valor como uma string vazia

// 		parts := strings.SplitN(tag, ":", 2)
// 		key := strings.TrimSpace(parts[0])
// 		if len(parts) == 2 {
// 			value = strings.TrimSpace(parts[1])
// 		}
// 		tagMap[key] = value
// 	}
// 	return tagMap
// }
