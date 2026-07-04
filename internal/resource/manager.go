//go:generate mockgen -source=manager.go -destination=manager_mock_test.go -package=resource_test
package resource

import (
	"errors"
)

var ErrResourceNotFound = errors.New("resource not found")

type Type string

type Resource interface {
	Update() error
	Type() Type
	Get() any
}

type Manager struct {
	resources map[Type]Resource
}

func NewManager(resources ...Resource) *Manager {
	m := &Manager{
		resources: make(map[Type]Resource),
	}

	for _, resource := range resources {
		if _, ok := m.resources[resource.Type()]; ok {
			panic("resource was registered already")
		}

		m.resources[resource.Type()] = resource
	}

	return m
}

func (u *Manager) Update() error {
	var errs []error
	for _, resource := range u.resources {
		if err := resource.Update(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...) // nil if len(errs) == nil
}

func (u *Manager) Get(resourceType Type) (Resource, error) {
	if resource, ok := u.resources[resourceType]; ok {
		return resource, nil
	}

	return nil, ErrResourceNotFound
}
