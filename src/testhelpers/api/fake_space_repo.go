package api

import (
	"cf"
	"cf/errors"
	"cf/models"
)

type FakeSpaceRepository struct {
	CurrentSpace models.Space

	Spaces []models.Space

	FindByNameName     string
	FindByNameSpace    models.Space
	FindByNameErr      bool
	FindByNameNotFound bool

	FindByNameInOrgName    string
	FindByNameInOrgOrgGuid string
	FindByNameInOrgSpace   models.Space

	SummarySpace models.Space

	CreateSpaceName    string
	CreateSpaceOrgGuid string
	CreateSpaceExists  bool
	CreateSpaceSpace   models.Space

	RenameSpaceGuid string
	RenameNewName   string

	DeletedSpaceGuid string
}

func (repo FakeSpaceRepository) GetCurrentSpace() (space models.Space) {
	return repo.CurrentSpace
}

func (repo FakeSpaceRepository) ListSpaces(callback func(models.Space) bool) errors.Error {
	for _, space := range repo.Spaces {
		if !callback(space) {
			break
		}
	}
	return nil
}

func (repo *FakeSpaceRepository) FindByName(name string) (space models.Space, apiErr errors.Error) {
	repo.FindByNameName = name

	var foundSpace bool = false
	for _, someSpace := range repo.Spaces {
		if name == someSpace.Name {
			foundSpace = true
			space = someSpace
			break
		}
	}

	if repo.FindByNameErr || !foundSpace {
		apiErr = errors.NewErrorWithMessage("Error finding space by name.")
	}

	if repo.FindByNameNotFound {
		apiErr = errors.NewModelNotFoundError("Space", name)
	}

	return
}

func (repo *FakeSpaceRepository) FindByNameInOrg(name, orgGuid string) (space models.Space, apiErr errors.Error) {
	repo.FindByNameInOrgName = name
	repo.FindByNameInOrgOrgGuid = orgGuid
	space = repo.FindByNameInOrgSpace
	return
}

func (repo *FakeSpaceRepository) GetSummary() (space models.Space, apiErr errors.Error) {
	space = repo.SummarySpace
	return
}

func (repo *FakeSpaceRepository) Create(name string, orgGuid string) (space models.Space, apiErr errors.Error) {
	if repo.CreateSpaceExists {
		apiErr = errors.NewError("Space already exists", cf.SPACE_EXISTS)
		return
	}
	repo.CreateSpaceName = name
	repo.CreateSpaceOrgGuid = orgGuid
	space = repo.CreateSpaceSpace
	return
}

func (repo *FakeSpaceRepository) Rename(spaceGuid, newName string) (apiErr errors.Error) {
	repo.RenameSpaceGuid = spaceGuid
	repo.RenameNewName = newName
	return
}

func (repo *FakeSpaceRepository) Delete(spaceGuid string) (apiErr errors.Error) {
	repo.DeletedSpaceGuid = spaceGuid
	return
}
