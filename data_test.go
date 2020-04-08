package main

import (
	"github.com/stretchr/testify/require"
	"github.com/uber/h3-go"
	"strings"
	"testing"
)

func TestSqlDb_AddUser(t *testing.T) {
	id, err := Db.AddUser(NewUser("Test_Mario Benedetti", -32.8181, -56.5064, ServiceProvider))
	require.NotEqual(t, id, 0)
	require.NoError(t, err)

	user, err := Db.GetUser(id)
	require.NoError(t, err)
	require.Equal(t, user.Id, id)
	require.Equal(t, user.Category, Category(ServiceProvider))
	require.Equal(t, user.GeoCord.Latitude, -32.8181)
	require.Equal(t, user.GeoCord.Longitude, -56.5064)
	require.Equal(t, len(user.H3Positions), 16)

	user1, err := Db.DeleteUser(id)
	require.NoError(t, err)
	require.Equal(t, user1.Id, id)
	require.Equal(t, user1.Category, Category(ServiceProvider))
	require.Equal(t, user1.GeoCord.Latitude, -32.8181)
	require.Equal(t, user1.GeoCord.Longitude, -56.5064)
	require.Equal(t, len(user1.H3Positions), 16)

}

func TestSqlDb_ListUsers(t *testing.T) {
	var err error

	user1 := NewUser("Test_Ruben Dario", 12.732229, -86.123326, Client)
	user1.Id, err = Db.AddUser(user1)
	require.NoError(t, err)

	user2 := NewUser("Test_Mario Benedetti", -32.8181, -56.5064, ServiceProvider)
	user2.Id, err = Db.AddUser(user2)
	require.NoError(t, err)
	user3 := NewUser("Test_Pablo Neruda", -36.143747, -71.827252, Client)
	user3.Id, err = Db.AddUser(user3)
	require.NoError(t, err)

	t.Log("List with category")
	userList, err := Db.ListUsers(Generic)
	require.NoError(t, err)
	require.Equal(t, len(userList), 3)

	for i := range userList {
		_, err = Db.DeleteUser(userList[i].Id)
		require.NoError(t, err)
	}

	require.Equal(t, user1.Name, userList[0].Name)
	require.Equal(t, user1.Category, userList[0].Category)
	require.Equal(t, user1.GeoCord.Longitude, userList[0].GeoCord.Longitude)
	require.Equal(t, user1.GeoCord.Latitude, userList[0].GeoCord.Latitude)

	require.Equal(t, user2.Name, userList[1].Name)
	require.Equal(t, user2.Category, userList[1].Category)
	require.Equal(t, user2.GeoCord.Longitude, userList[1].GeoCord.Longitude)
	require.Equal(t, user2.GeoCord.Latitude, userList[1].GeoCord.Latitude)

	require.Equal(t, user3.Name, userList[2].Name)
	require.Equal(t, user3.Category, userList[2].Category)
	require.Equal(t, user3.GeoCord.Longitude, userList[2].GeoCord.Longitude)
	require.Equal(t, user3.GeoCord.Latitude, userList[2].GeoCord.Latitude)

	require.Equal(t, h3.H3Index(userList[0].H3Positions[0]), h3.H3Index(0x806dfffffffffff))
	require.Equal(t, h3.H3Index(userList[1].H3Positions[0]), h3.H3Index(0x80c3fffffffffff))
	require.Equal(t, h3.H3Index(userList[2].H3Positions[0]), h3.H3Index(0x80b3fffffffffff))
	require.False(t, h3.AreNeighbors(h3.H3Index(userList[0].H3Positions[0]), h3.H3Index(userList[1].H3Positions[0])))
	require.False(t, h3.AreNeighbors(h3.H3Index(userList[0].H3Positions[0]), h3.H3Index(userList[2].H3Positions[0])))
	require.True(t, h3.AreNeighbors(h3.H3Index(userList[1].H3Positions[0]), h3.H3Index(userList[2].H3Positions[0])))

}

func TestSqlDb_UpdateUser(t *testing.T) {
	var err error
	user1 := NewUser("Test_Ruben Dario", 12.732229, -86.123326, Client)
	user1.Id, err = Db.AddUser(user1)
	require.NoError(t, err)
	userToUpdate := NewUser(user1.Name, user1.GeoCord.Latitude+1, user1.GeoCord.Longitude-1, user1.Category)

	newUser, err := Db.UpdateUser(user1.Id, userToUpdate.GeoCord.Latitude, userToUpdate.GeoCord.Longitude, userToUpdate.H3Positions)
	require.NoError(t, err)
	require.Equal(t, user1.Id, newUser.Id)
	require.Equal(t, user1.GeoCord.Longitude, newUser.GeoCord.Longitude+1)
	require.Equal(t, user1.GeoCord.Latitude, newUser.GeoCord.Latitude-1)

	_, err = Db.DeleteUser(user1.Id)
	require.NoError(t, err)
}

func TestSqlDb_GetCloseUsers(t *testing.T) {
	var err error
	user1 := NewUser("Test_Ruben Dario", 12.732229, -86.123326, Client)
	user1.Id, err = Db.AddUser(user1)
	require.NoError(t, err)

	user2 := NewUser("Test_Mario Benedetti", -32.8181, -56.5064, ServiceProvider)
	user2.Id, err = Db.AddUser(user2)
	require.NoError(t, err)

	user3 := NewUser("Test_Pablo Neruda", -36.143747, -71.827252, Client)
	user3.Id, err = Db.AddUser(user3)
	require.NoError(t, err)

	user4 := NewUser("Test_Victor Jara", -36.534865, -72.434204, ServiceProvider)
	user4.Id, err = Db.AddUser(user4)
	require.NoError(t, err)

	resolution := 1
	userList, err := Db.GetCloseUsers(resolution, user3.H3Positions[resolution], Generic)
	require.NoError(t, err)
	require.Len(t, userList, 2)

	userList1, err := Db.GetCloseUsers(resolution, user1.H3Positions[resolution], Client)
	require.NoError(t, err)
	require.Len(t, userList1, 1)
}

func TestSqlDb_Close(t *testing.T) {
	t.Cleanup(func() {
		users, err := Db.ListUsers(Generic)
		require.NoError(t, err)
		for i := range users {
			if strings.Split(users[i].Name, "_")[0] == "Test" {
				_, err := Db.DeleteUser(users[i].Id)
				require.NoError(t, err)
			}
		}
		defer Db.Close()
	})
}
