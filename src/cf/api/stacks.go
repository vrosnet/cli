package api

import (
	"cf/configuration"
	"cf/errors"
	"cf/models"
	"cf/net"
	"fmt"
	"net/url"
)

type PaginatedStackResources struct {
	Resources []StackResource
}

type StackResource struct {
	Resource
	Entity StackEntity
}

func (resource StackResource) ToFields() (fields models.Stack) {
	fields.Guid = resource.Metadata.Guid
	fields.Name = resource.Entity.Name
	fields.Description = resource.Entity.Description
	return
}

type StackEntity struct {
	Name        string
	Description string
}

type StackRepository interface {
	FindByName(name string) (stack models.Stack, apiErr errors.Error)
	FindAll() (stacks []models.Stack, apiErr errors.Error)
}

type CloudControllerStackRepository struct {
	config  configuration.Reader
	gateway net.Gateway
}

func NewCloudControllerStackRepository(config configuration.Reader, gateway net.Gateway) (repo CloudControllerStackRepository) {
	repo.config = config
	repo.gateway = gateway
	return
}

func (repo CloudControllerStackRepository) FindByName(name string) (stack models.Stack, apiErr errors.Error) {
	path := fmt.Sprintf("%s/v2/stacks?q=%s", repo.config.ApiEndpoint(), url.QueryEscape("name:"+name))
	stacks, apiErr := repo.findAllWithPath(path)
	if apiErr != nil {
		return
	}

	if len(stacks) == 0 {
		apiErr = errors.NewErrorWithMessage("Stack '%s' not found", name)
		return
	}

	stack = stacks[0]
	return
}

func (repo CloudControllerStackRepository) FindAll() (stacks []models.Stack, apiErr errors.Error) {
	path := fmt.Sprintf("%s/v2/stacks", repo.config.ApiEndpoint())
	return repo.findAllWithPath(path)
}

func (repo CloudControllerStackRepository) findAllWithPath(path string) (stacks []models.Stack, apiErr errors.Error) {
	resources := new(PaginatedStackResources)
	apiErr = repo.gateway.GetResource(path, repo.config.AccessToken(), resources)
	if apiErr != nil {
		return
	}

	for _, r := range resources.Resources {
		stacks = append(stacks, r.ToFields())
	}
	return
}
