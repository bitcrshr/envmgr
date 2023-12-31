// Code generated by ent, DO NOT EDIT.

package ent

import (
	"github.com/bitcrshr/envmgr/api/ent/environment"
	"github.com/bitcrshr/envmgr/api/ent/project"
	"github.com/bitcrshr/envmgr/api/ent/schema"
	"github.com/bitcrshr/envmgr/api/ent/user"
	"github.com/bitcrshr/envmgr/api/ent/variable"
	"github.com/google/uuid"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	environmentFields := schema.Environment{}.Fields()
	_ = environmentFields
	// environmentDescName is the schema descriptor for name field.
	environmentDescName := environmentFields[1].Descriptor()
	// environment.NameValidator is a validator for the "name" field. It is called by the builders before save.
	environment.NameValidator = environmentDescName.Validators[0].(func(string) error)
	projectFields := schema.Project{}.Fields()
	_ = projectFields
	// projectDescDisplayName is the schema descriptor for display_name field.
	projectDescDisplayName := projectFields[1].Descriptor()
	// project.DisplayNameValidator is a validator for the "display_name" field. It is called by the builders before save.
	project.DisplayNameValidator = projectDescDisplayName.Validators[0].(func(string) error)
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescAuth0ID is the schema descriptor for auth0_id field.
	userDescAuth0ID := userFields[1].Descriptor()
	// user.Auth0IDValidator is a validator for the "auth0_id" field. It is called by the builders before save.
	user.Auth0IDValidator = userDescAuth0ID.Validators[0].(func(string) error)
	// userDescAppMetadata is the schema descriptor for app_metadata field.
	userDescAppMetadata := userFields[10].Descriptor()
	// user.DefaultAppMetadata holds the default value on creation for the app_metadata field.
	user.DefaultAppMetadata = userDescAppMetadata.Default.(map[string]string)
	// userDescUserMetadata is the schema descriptor for user_metadata field.
	userDescUserMetadata := userFields[11].Descriptor()
	// user.DefaultUserMetadata holds the default value on creation for the user_metadata field.
	user.DefaultUserMetadata = userDescUserMetadata.Default.(map[string]string)
	// userDescID is the schema descriptor for id field.
	userDescID := userFields[0].Descriptor()
	// user.DefaultID holds the default value on creation for the id field.
	user.DefaultID = userDescID.Default.(func() uuid.UUID)
	variableFields := schema.Variable{}.Fields()
	_ = variableFields
	// variableDescKey is the schema descriptor for key field.
	variableDescKey := variableFields[1].Descriptor()
	// variable.KeyValidator is a validator for the "key" field. It is called by the builders before save.
	variable.KeyValidator = variableDescKey.Validators[0].(func(string) error)
	// variableDescValue is the schema descriptor for value field.
	variableDescValue := variableFields[2].Descriptor()
	// variable.DefaultValue holds the default value on creation for the value field.
	variable.DefaultValue = variableDescValue.Default.(string)
}
