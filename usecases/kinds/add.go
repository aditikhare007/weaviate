//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2019 Weaviate. All rights reserved.
//  LICENSE WEAVIATE OPEN SOURCE: https://www.semi.technology/playbook/playbook/contract-weaviate-OSS.html
//  LICENSE WEAVIATE ENTERPRISE: https://www.semi.technology/playbook/contract-weaviate-enterprise.html
//  CONCEPT: Bob van Luijt (@bobvanluijt)
//  CONTACT: hello@semi.technology
//

package kinds

import (
	"context"

	"github.com/go-openapi/strfmt"
	uuid "github.com/satori/go.uuid"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/entities/schema"
	"github.com/semi-technologies/weaviate/entities/schema/kind"
	"github.com/semi-technologies/weaviate/usecases/kinds/validation"
)

type addAndGetRepo interface {
	getRepo
	addRepo
}

type addRepo interface {
	AddAction(ctx context.Context, class *models.Action, id strfmt.UUID) error
	AddThing(ctx context.Context, class *models.Thing, id strfmt.UUID) error
	ClassExists(ctx context.Context, id strfmt.UUID) (bool, error)
}

type schemaManager interface {
	UpdatePropertyAddDataType(context.Context, *models.Principal, kind.Kind, string, string, string) error
	GetSchema(principal *models.Principal) (schema.Schema, error)
}

// AddAction Class Instance to the connected DB. If the class contains a network
// ref, it has a side-effect on the schema: The schema will be updated to
// include this particular network ref class.
func (m *Manager) AddAction(ctx context.Context, principal *models.Principal,
	class *models.Action) (*models.Action, error) {

	err := m.authorizer.Authorize(principal, "create", "actions")
	if err != nil {
		return nil, err
	}

	unlock, err := m.locks.LockSchema()
	if err != nil {
		return nil, NewErrInternal("could not aquire lock: %v", err)
	}
	defer unlock()

	return m.addActionToConnectorAndSchema(ctx, principal, class)
}

func (m *Manager) checkIDOrAssignNew(ctx context.Context, id strfmt.UUID) (strfmt.UUID, error) {
	if id == "" {
		newID, err := generateUUID()
		if err != nil {
			return "", NewErrInternal("could not generate id: %v", err)
		}
		return newID, nil
	}

	// only validate ID uniqueness if explicitly set
	if ok, err := m.repo.ClassExists(ctx, id); ok {
		return "", NewErrInvalidUserInput("id '%s' already exists", id)
	} else if err != nil {
		return "", NewErrInternal(err.Error())
	}
	return id, nil
}

func (m *Manager) addActionToConnectorAndSchema(ctx context.Context, principal *models.Principal,
	class *models.Action) (*models.Action, error) {
	id, err := m.checkIDOrAssignNew(ctx, class.ID)
	if err != nil {
		return nil, err
	}
	class.ID = id

	err = m.validateAction(ctx, principal, class)
	if err != nil {
		return nil, NewErrInvalidUserInput("invalid action: %v", err)
	}

	err = m.addNetworkDataTypesForAction(ctx, principal, class)
	if err != nil {
		return nil, NewErrInternal("could not update schema for network refs: %v", err)
	}

	err = m.repo.AddAction(ctx, class, class.ID)
	if err != nil {
		return nil, NewErrInternal("could not store action: %v", err)
	}

	v, err := m.vectorizer.Action(ctx, class)
	if err != nil {
		return nil, NewErrInternal("could not create vector from action: %v", err)
	}

	err = m.vectorRepo.PutAction(ctx, class, v)
	if err != nil {
		return nil, NewErrInternal("could not store vector for thing: %v", err)
	}

	return class, nil
}

func (m *Manager) validateAction(ctx context.Context, principal *models.Principal, class *models.Action) error {
	// Validate schema given in body with the weaviate schema
	if _, err := uuid.FromString(class.ID.String()); err != nil {
		return err
	}

	s, err := m.schemaManager.GetSchema(principal)
	if err != nil {
		return err
	}

	databaseSchema := schema.HackFromDatabaseSchema(s)
	return validation.ValidateActionBody(
		ctx, class, databaseSchema, m.repo, m.network, m.config)
}

// AddThing Class Instance to the connected DB. If the class contains a network
// ref, it has a side-effect on the schema: The schema will be updated to
// include this particular network ref class.
func (m *Manager) AddThing(ctx context.Context, principal *models.Principal,
	class *models.Thing) (*models.Thing, error) {

	err := m.authorizer.Authorize(principal, "create", "things")
	if err != nil {
		return nil, err
	}

	unlock, err := m.locks.LockSchema()
	if err != nil {
		return nil, NewErrInternal("could not aquire lock: %v", err)
	}
	defer unlock()

	return m.addThingToConnectorAndSchema(ctx, principal, class)
}

func (m *Manager) addThingToConnectorAndSchema(ctx context.Context, principal *models.Principal,
	class *models.Thing) (*models.Thing, error) {
	id, err := m.checkIDOrAssignNew(ctx, class.ID)
	if err != nil {
		return nil, err
	}
	class.ID = id

	err = m.validateThing(ctx, principal, class)
	if err != nil {
		return nil, NewErrInvalidUserInput("invalid thing: %v", err)
	}

	err = m.addNetworkDataTypesForThing(ctx, principal, class)
	if err != nil {
		return nil, NewErrInternal("could not update schema for network refs: %v", err)
	}

	err = m.repo.AddThing(ctx, class, class.ID)
	if err != nil {
		return nil, NewErrInternal("could not store thing: %v", err)
	}

	v, err := m.vectorizer.Thing(ctx, class)
	if err != nil {
		return nil, NewErrInternal("could not create vector from thing: %v", err)
	}

	err = m.vectorRepo.PutThing(ctx, class, v)
	if err != nil {
		return nil, NewErrInternal("could not store vector for thing: %v", err)
	}

	return class, nil
}

func (m *Manager) validateThing(ctx context.Context, principal *models.Principal,
	class *models.Thing) error {
	// Validate schema given in body with the weaviate schema
	if _, err := uuid.FromString(class.ID.String()); err != nil {
		return err
	}

	s, err := m.schemaManager.GetSchema(principal)
	if err != nil {
		return err
	}

	// Validate schema given in body with the weaviate schema
	databaseSchema := schema.HackFromDatabaseSchema(s)
	return validation.ValidateThingBody(
		ctx, class, databaseSchema, m.repo, m.network, m.config)
}

func (m *Manager) addNetworkDataTypesForThing(ctx context.Context, principal *models.Principal, class *models.Thing) error {
	refSchemaUpdater := newReferenceSchemaUpdater(ctx, principal, m.schemaManager, m.network, class.Class, kind.Thing)
	return refSchemaUpdater.addNetworkDataTypes(class.Schema)
}

func (m *Manager) addNetworkDataTypesForAction(ctx context.Context, principal *models.Principal, class *models.Action) error {
	refSchemaUpdater := newReferenceSchemaUpdater(ctx, principal, m.schemaManager, m.network, class.Class, kind.Action)
	return refSchemaUpdater.addNetworkDataTypes(class.Schema)
}