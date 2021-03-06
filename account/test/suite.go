package test

import (
	"testing"

	"github.com/apihub/apihub/account"
	"github.com/apihub/apihub/errors"
	. "gopkg.in/check.v1"
)

var app account.App
var user account.User
var plugin account.Plugin
var team account.Team
var service account.Service
var token account.Token
var hook account.Hook

//Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type StorableSuite struct {
	Storage account.Storable
}

func (s *StorableSuite) SetUpTest(c *C) {
	user = account.User{Name: "Alice", Email: "alice@example.org", Password: "123456"}
	token = account.Token{AccessToken: "secret-token", Expires: 10, Type: "Token", User: &user}
	team = account.Team{Name: "ApiHub Team", Alias: "apihub", Users: []string{user.Email}, Owner: user.Email, Apps: []account.App{}, Services: []account.Service{}}
	service = account.Service{Endpoint: "http://example.org/api", Subdomain: "apihub", Team: team.Alias, Owner: user.Email, Transformers: []string{}}
	app = account.App{ClientId: "ios", ClientSecret: "secret", Name: "Ios App", Team: team.Alias, Owner: user.Email, RedirectUris: []string{"http://www.example.org/auth"}}
	plugin = account.Plugin{Name: "cors", Service: service.Subdomain, Config: map[string]interface{}{"version": 1}}
	hook = account.Hook{Name: "service.update", Events: []string{"service.update"}, Config: account.HookConfig{Address: "http://www.example.org"}}
}

func (s *StorableSuite) TestUpsertUser(c *C) {
	defer s.Storage.DeleteUser(user)
	err := s.Storage.UpsertUser(user)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestUpdateUser(c *C) {
	s.Storage.UpsertUser(user)
	user.Name = "Bob"
	defer s.Storage.DeleteUser(user)
	err := s.Storage.UpsertUser(user)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestTeams(c *C) {
	defer s.Storage.DeleteTeam(team)
	s.Storage.UpsertTeam(team)
	teams, err := s.Storage.UserTeams(account.User{Email: team.Owner})
	c.Check(err, IsNil)
	c.Assert(teams, DeepEquals, []account.Team{team})
}

func (s *StorableSuite) TestTeamsNotFound(c *C) {
	teams, err := s.Storage.UserTeams(account.User{})
	c.Check(err, IsNil)
	c.Assert(teams, DeepEquals, []account.Team{})
}

func (s *StorableSuite) TestDeleteUser(c *C) {
	s.Storage.UpsertUser(user)
	err := s.Storage.DeleteUser(user)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeleteUserNotFound(c *C) {
	err := s.Storage.DeleteUser(user)
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestFindUserByEmail(c *C) {
	defer s.Storage.DeleteUser(user)
	s.Storage.UpsertUser(user)
	u, err := s.Storage.FindUserByEmail(user.Email)
	c.Assert(u, Equals, user)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestFindUserByEmailNotFound(c *C) {
	_, err := s.Storage.FindUserByEmail("not-found")
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestUpsertTeam(c *C) {
	defer s.Storage.DeleteTeam(team)
	err := s.Storage.UpsertTeam(team)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeleteTeam(c *C) {
	s.Storage.UpsertTeam(team)
	err := s.Storage.DeleteTeam(team)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeleteTeamNotFound(c *C) {
	err := s.Storage.DeleteTeam(team)
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestDeleteTeamByAlias(c *C) {
	s.Storage.UpsertTeam(team)
	err := s.Storage.DeleteTeamByAlias(team.Alias)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeleteTeamByAliasNotFound(c *C) {
	err := s.Storage.DeleteTeamByAlias(team.Alias)
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestFindTeamByAlias(c *C) {
	defer s.Storage.DeleteTeam(team)
	s.Storage.UpsertTeam(team)
	u, err := s.Storage.FindTeamByAlias(team.Alias)
	c.Assert(u, DeepEquals, team)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestFindTeamByAliasNotFound(c *C) {
	_, err := s.Storage.FindTeamByAlias("not-found")
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestTeamServices(c *C) {
	s.Storage.UpsertService(service)
	defer s.Storage.DeleteService(service)

	services, err := s.Storage.TeamServices(team)
	c.Assert(err, IsNil)
	c.Assert(services, DeepEquals, []account.Service{service})
}

func (s *StorableSuite) TestTeamServiceNotFound(c *C) {
	services, err := s.Storage.TeamServices(team)
	c.Assert(err, IsNil)
	c.Assert(services, DeepEquals, []account.Service{})
}

func (s *StorableSuite) TestCreateToken(c *C) {
	defer s.Storage.DeleteToken(token.AccessToken)
	err := s.Storage.CreateToken(token)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeleteToken(c *C) {
	err := s.Storage.CreateToken(token)
	c.Check(err, IsNil)
	err = s.Storage.DeleteToken(token.AccessToken)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDecodeToken(c *C) {
	s.Storage.CreateToken(token)
	var u account.User
	s.Storage.DecodeToken(token.AccessToken, &u)
	c.Assert(u, DeepEquals, user)
}

func (s *StorableSuite) TestUpsertService(c *C) {
	defer s.Storage.DeleteService(service)
	err := s.Storage.UpsertService(service)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeleteService(c *C) {
	s.Storage.UpsertService(service)
	err := s.Storage.DeleteService(service)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeleteServiceNotFound(c *C) {
	err := s.Storage.DeleteService(service)
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestFindServiceBySubdomain(c *C) {
	defer s.Storage.DeleteService(service)
	s.Storage.UpsertService(service)
	serv, err := s.Storage.FindServiceBySubdomain(service.Subdomain)
	c.Assert(serv, DeepEquals, service)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestFindServiceBySubdomainNotFound(c *C) {
	_, err := s.Storage.FindServiceBySubdomain("not-found")
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestUserServices(c *C) {
	s.Storage.UpsertTeam(team)
	defer s.Storage.DeleteTeam(team)
	s.Storage.UpsertService(service)
	defer s.Storage.DeleteService(service)
	another_service := account.Service{Endpoint: "http://example.org/api", Subdomain: "example", Team: team.Alias, Owner: user.Email, Transformers: []string{}}
	s.Storage.UpsertService(another_service)
	defer s.Storage.DeleteService(another_service)

	services, err := s.Storage.UserServices(account.User{Email: team.Owner})
	c.Check(err, IsNil)
	c.Assert(len(services), Equals, 2)
}

func (s *StorableSuite) TestUserServicesNotFound(c *C) {
	services, err := s.Storage.UserServices(user)
	c.Assert(err, IsNil)
	c.Assert(services, DeepEquals, []account.Service{})
}

func (s *StorableSuite) TestUpsertApp(c *C) {
	defer s.Storage.DeleteApp(app)
	err := s.Storage.UpsertApp(app)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeleteApp(c *C) {
	s.Storage.UpsertApp(app)
	err := s.Storage.DeleteApp(app)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestTeamApps(c *C) {
	s.Storage.UpsertApp(app)
	defer s.Storage.DeleteApp(app)

	apps, err := s.Storage.TeamApps(team)
	c.Assert(err, IsNil)
	c.Assert(apps, DeepEquals, []account.App{app})
}

func (s *StorableSuite) TestTeamAppNotFound(c *C) {
	apps, err := s.Storage.TeamApps(team)
	c.Assert(err, IsNil)
	c.Assert(apps, DeepEquals, []account.App{})
}

func (s *StorableSuite) TestDeleteAppNotFound(c *C) {
	nf := account.App{}
	err := s.Storage.DeleteApp(nf)
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestFindAppByClientId(c *C) {
	defer s.Storage.DeleteApp(app)
	s.Storage.UpsertApp(app)
	a, err := s.Storage.FindAppByClientId(app.ClientId)
	c.Assert(a, DeepEquals, app)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestFindAppByClientIdNotFound(c *C) {
	_, err := s.Storage.FindAppByClientId("not-found")
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestUpsertPlugin(c *C) {
	defer s.Storage.DeletePlugin(plugin)
	err := s.Storage.UpsertPlugin(plugin)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeletePlugin(c *C) {
	s.Storage.UpsertPlugin(plugin)
	err := s.Storage.DeletePlugin(plugin)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeletePluginNotFound(c *C) {
	nf := account.Plugin{}
	err := s.Storage.DeletePlugin(nf)
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestDeletePluginsByServiceNotFound(c *C) {
	nf := account.Service{}
	err := s.Storage.DeletePluginsByService(nf)
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestDeletePluginsByService(c *C) {
	err := s.Storage.UpsertService(service)
	c.Check(err, IsNil)
	plugin.Service = service.Subdomain
	err = s.Storage.UpsertPlugin(plugin)
	c.Check(err, IsNil)

	err = s.Storage.DeletePluginsByService(service)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestFindPluginByNameAndService(c *C) {
	defer s.Storage.DeletePlugin(plugin)
	plugin.Service = service.Subdomain
	err := s.Storage.UpsertPlugin(plugin)
	pl, err := s.Storage.FindPluginByNameAndService(plugin.Name, service)
	c.Assert(pl, DeepEquals, plugin)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestFindPluginByNameAndServiceNotFound(c *C) {
	_, err := s.Storage.FindPluginByNameAndService("not-found", service)
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestUpsertHook(c *C) {
	defer s.Storage.DeleteHook(hook)
	err := s.Storage.UpsertHook(hook)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeleteHook(c *C) {
	s.Storage.UpsertHook(hook)
	err := s.Storage.DeleteHook(hook)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestDeleteHookNotFound(c *C) {
	err := s.Storage.DeleteHook(hook)
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestDeleteHooksByTeamNotFound(c *C) {
	nf := account.Team{}
	err := s.Storage.DeleteHooksByTeam(nf)
	_, ok := err.(errors.NotFoundError)
	c.Assert(ok, Equals, true)
}

func (s *StorableSuite) TestDeleteHooksByTeam(c *C) {
	err := s.Storage.UpsertTeam(team)
	c.Check(err, IsNil)
	hook.Team = team.Alias
	err = s.Storage.UpsertHook(hook)
	c.Check(err, IsNil)

	err = s.Storage.DeleteHooksByTeam(team)
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestFindAllHooksByEventAndTeam(c *C) {
	defer s.Storage.DeleteTeam(team)
	s.Storage.UpsertTeam(team)

	defer s.Storage.DeleteHook(hook)
	hook.Events = []string{"service.create"}
	hook.Team = team.Alias
	s.Storage.UpsertHook(hook)

	whs, err := s.Storage.FindHooksByEventAndTeam("service.create", account.ALL_TEAMS)
	c.Assert(whs, DeepEquals, []account.Hook{hook})
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestFindHooksByEventAndTeam(c *C) {
	defer s.Storage.DeleteTeam(team)
	s.Storage.UpsertTeam(team)

	defer s.Storage.DeleteHook(hook)
	hook.Name = "service.create"
	hook.Events = []string{"service.create"}
	hook.Team = team.Alias
	s.Storage.UpsertHook(hook)

	whk := account.Hook{
		Name:   "service.update",
		Events: []string{"service.update"},
		Config: account.HookConfig{Address: "http://www.example.org"},
	}
	defer s.Storage.DeleteHook(whk)
	whk.Events = []string{"service.update"}
	whk.Team = team.Alias
	s.Storage.UpsertHook(whk)

	whs, err := s.Storage.FindHooksByEventAndTeam("service.create", team.Alias)
	c.Assert(whs, DeepEquals, []account.Hook{hook})
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestFindHooksByEventAndTeamNotFound(c *C) {
	whs, err := s.Storage.FindHooksByEventAndTeam("not-found", "not-found")
	c.Assert(whs, DeepEquals, []account.Hook{})
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestFindHooksByEvent(c *C) {
	defer s.Storage.DeleteTeam(team)
	s.Storage.UpsertTeam(team)

	defer s.Storage.DeleteHook(hook)
	hook.Name = "service.create"
	hook.Events = []string{"service.create"}
	hook.Team = team.Alias
	s.Storage.UpsertHook(hook)

	whk := account.Hook{
		Name:   "service.update",
		Events: []string{"service.update"},
		Config: account.HookConfig{Address: "http://www.example.org"},
	}
	defer s.Storage.DeleteHook(whk)
	whk.Events = []string{"service.update"}
	whk.Team = team.Alias
	s.Storage.UpsertHook(whk)

	whs, err := s.Storage.FindHooksByEvent("service.create")
	c.Assert(whs, DeepEquals, []account.Hook{hook})
	c.Check(err, IsNil)
}

func (s *StorableSuite) TestFindHooksByEventNotFound(c *C) {
	whs, err := s.Storage.FindHooksByEvent("not-found")
	c.Assert(whs, DeepEquals, []account.Hook{})
	c.Check(err, IsNil)
}
