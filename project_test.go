package dtrack

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProjectService_Clone(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		APIPermissions: []string{
			PermissionPortfolioManagement,
		},
	})

	project, err := client.Project.Create(context.Background(), Project{
		Name:    "acme-app",
		Version: "1.0.0",
	})
	require.NoError(t, err)

	token, err := client.Project.Clone(context.Background(), ProjectCloneRequest{
		ProjectUUID: project.UUID,
		Version:     "2.0.0",
	})
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestProjectService_Clone_v4_10(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		Version: "4.10.1",
		APIPermissions: []string{
			PermissionPortfolioManagement,
		},
	})

	project, err := client.Project.Create(context.Background(), Project{
		Name:    "acme-app",
		Version: "1.0.0",
	})
	require.NoError(t, err)

	token, err := client.Project.Clone(context.Background(), ProjectCloneRequest{
		ProjectUUID: project.UUID,
		Version:     "2.0.0",
	})
	require.NoError(t, err)
	require.Empty(t, token)
}
func TestProjectService_Latest(t *testing.T) {
	client := setUpContainer(t, testContainerOptions{
		Version: "4.12.7",
		APIPermissions: []string{
			PermissionPortfolioManagement,
			PermissionViewPortfolio,
		},
	})
	var name = "acme-app"
	project, err := client.Project.Create(context.Background(), Project{
		Name:     name,
		Version:  "1.0.0",
		IsLatest: OptionalBoolOf(true),
	})
	require.NoError(t, err)
	latest, err := client.Project.Latest(context.Background(), name)

	require.NoError(t, err)
	require.Equal(t, project.Version, latest.Version)

	token, err := client.Project.Clone(context.Background(), ProjectCloneRequest{
		ProjectUUID:     project.UUID,
		Version:         "2.0.0",
		MakeCloneLatest: OptionalBoolOf(true),
	})
	require.NoError(t, err)
	require.NotEmpty(t, token)

	for {
		processing, err := client.Event.IsBeingProcessed(context.Background(), token)
		require.NoError(t, err)
		if !processing {
			break
		}
	}

	latest, err = client.Project.Latest(context.Background(), name)

	require.NoError(t, err)
	require.Equal(t, "2.0.0", latest.Version)

}
