package requirements

import (
	"cf"
	"cf/configuration"
	"cf/terminal"
	"fmt"
)

type TargetedOrgRequirement interface {
	Requirement
	GetOrganization() cf.OrganizationFields
}

type targetedOrgApiRequirement struct {
	ui     terminal.UI
	config *configuration.Configuration
}

func newTargetedOrgRequirement(ui terminal.UI, config *configuration.Configuration) TargetedOrgRequirement {
	return targetedOrgApiRequirement{ui, config}
}

func (req targetedOrgApiRequirement) Execute() (success bool) {
	if !req.config.HasOrganization() {
		message := fmt.Sprintf("No org targeted, use '%s' to target an org.",
			terminal.CommandColor(cf.Name()+" target -o ORG"))
		req.ui.Failed(message)
		return false
	}

	return true
}

func (req targetedOrgApiRequirement) GetOrganization() (org cf.OrganizationFields) {
	return req.config.Organization
}
