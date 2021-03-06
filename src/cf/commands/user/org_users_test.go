package user_test

import (
	. "cf/commands/user"
	"cf/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	testapi "testhelpers/api"
	testassert "testhelpers/assert"
	testcmd "testhelpers/commands"
	testconfig "testhelpers/configuration"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
)

func callOrgUsers(args []string, reqFactory *testreq.FakeReqFactory, userRepo *testapi.FakeUserRepository) (ui *testterm.FakeUI) {
	ui = &testterm.FakeUI{}

	config := testconfig.NewRepositoryWithDefaults()

	cmd := NewOrgUsers(ui, config, userRepo)
	ctxt := testcmd.NewContext("org-users", args)

	testcmd.RunCommand(cmd, ctxt, reqFactory)
	return
}

var _ = Describe("Listing users in an org", func() {
	It("TestOrgUsersFailsWithUsage", func() {
		reqFactory := &testreq.FakeReqFactory{}
		userRepo := &testapi.FakeUserRepository{}
		ui := callOrgUsers([]string{}, reqFactory, userRepo)
		Expect(ui.FailedWithUsage).To(BeTrue())

		ui = callOrgUsers([]string{"Org1"}, reqFactory, userRepo)
		Expect(ui.FailedWithUsage).To(BeFalse())
	})

	It("TestOrgUsersRequirements", func() {
		reqFactory := &testreq.FakeReqFactory{}
		userRepo := &testapi.FakeUserRepository{}
		args := []string{"Org1"}

		reqFactory.LoginSuccess = false
		callOrgUsers(args, reqFactory, userRepo)
		Expect(testcmd.CommandDidPassRequirements).To(BeFalse())

		reqFactory.LoginSuccess = true
		callOrgUsers(args, reqFactory, userRepo)
		Expect(testcmd.CommandDidPassRequirements).To(BeTrue())

		Expect("Org1").To(Equal(reqFactory.OrganizationName))
	})

	It("TestOrgUsers", func() {
		org := models.Organization{}
		org.Name = "Found Org"
		org.Guid = "found-org-guid"

		userRepo := &testapi.FakeUserRepository{}
		user := models.UserFields{}
		user.Username = "user1"
		user2 := models.UserFields{}
		user2.Username = "user2"
		user3 := models.UserFields{}
		user3.Username = "user3"
		user4 := models.UserFields{}
		user4.Username = "user4"
		userRepo.ListUsersByRole = map[string][]models.UserFields{
			models.ORG_MANAGER:     []models.UserFields{user, user2},
			models.BILLING_MANAGER: []models.UserFields{user4},
			models.ORG_AUDITOR:     []models.UserFields{user3},
		}

		reqFactory := &testreq.FakeReqFactory{
			LoginSuccess: true,
			Organization: org,
		}

		ui := callOrgUsers([]string{"Org1"}, reqFactory, userRepo)

		Expect(userRepo.ListUsersOrganizationGuid).To(Equal("found-org-guid"))
		testassert.SliceContains(ui.Outputs, testassert.Lines{
			{"Getting users in org", "Found Org", "my-user"},
			{"ORG MANAGER"},
			{"user1"},
			{"user2"},
			{"BILLING MANAGER"},
			{"user4"},
			{"ORG AUDITOR"},
			{"user3"},
		})
	})

	It("lists all org users", func() {
		org := models.Organization{}
		org.Name = "Found Org"
		org.Guid = "found-org-guid"

		userRepo := &testapi.FakeUserRepository{}
		user := models.UserFields{}
		user.Username = "user1"
		user2 := models.UserFields{}
		user2.Username = "user2"
		userRepo.ListUsersByRole = map[string][]models.UserFields{
			models.ORG_USER: []models.UserFields{user, user2},
		}

		reqFactory := &testreq.FakeReqFactory{
			LoginSuccess: true,
			Organization: org,
		}

		ui := callOrgUsers([]string{"-a", "Org1"}, reqFactory, userRepo)

		Expect(userRepo.ListUsersOrganizationGuid).To(Equal("found-org-guid"))
		testassert.SliceContains(ui.Outputs, testassert.Lines{
			{"Getting users in org", "Found Org", "my-user"},
			{"USERS"},
			{"user1"},
			{"user2"},
		})
	})
})
