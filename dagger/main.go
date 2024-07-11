
package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
)

type Go01 struct{}

func (m *Go01) Publish(ctx context.Context, source *Directory) (string, error) {
	builder := dag.Container().
		From("golang:1.22.1").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithEnvVariable("CGO_ENABLED", "0").
		WithEnvVariable("GOOS", "linux").
		WithExec([]string{"go", "build", "-o", "../bin/main"})

	prodImage := dag.Container().
		From("golang:1.22.1-alpine3.19").
		WithFile("/go/bin/main", builder.File("/bin/main")).
		WithFile("/go/bin/test.json", builder.File("/src/src/test.json")).
		WithWorkdir("/go/bin").
		WithExec([]string{"adduser", "--disabled-password", "--gecos", "--quiet", "--shell", "/bin/bash", "--u", "1000", "noonroot"}).
		WithExec([]string{"chown", "-R", "1000:1000", "/go"}).
		WithEntrypoint([]string{"main"})

	address, err := prodImage.Publish(ctx, fmt.Sprintf("filstimatias/dagger-test:%.0f", math.Floor(rand.Float64()*100)))
	if err != nil {
		return "", err
	}
	return address, nil
}


func (m *Go01) TestAll(ctx context.Context, source *Directory) (string, error) {
	result, err := m.Lint(ctx, source)
	if err != nil {
		return "", err
	}

	return result, nil
}


func (m *Go01) Lint(ctx context.Context, source *Directory) (string, error) {
	return m.Test(ctx, source).
		WithExec([]string{"go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1"}).
		WithExec([]string{"pwd"}).
		WithExec([]string{"golangci-lint", "run", "./src", "./modules/...", "--issues-exit-code=1"}).
		Stdout(ctx)
}
// Returns a container that echoes whatever string argument is provided
func (m *Go01) Test(ctx context.Context, source *Directory) *Container {
	result := m.BuildEnv(source).
		WithExec([]string{"go", "test", "./...", "-v"}).
		WithExec([]string{"go", "mod", "verify"}).
		WithExec([]string{"go", "mod", "download"}).
		WithExec([]string{"go", "build", "-v", "./..."})
	return result
}




func (m *Go01) ConnectTest() *Service {
	return dag.Container().
		From("filstimatias/dagger-test-aux:37").
		WithExposedPort(8080).
		AsService()

}

func (m *Go01) Get(ctx context.Context, source *Directory) (string, error) {

	return m.BuildEnv(source).
		WithServiceBinding("default", m.ConnectTest()).
		WithExec([]string{"go", "run", "/src/src/main.go", "http://default:8080/response/1"}).
		Stdout(ctx)
}


// Build a ready-to-use development environment
func (m *Go01) BuildEnv(source *Directory) *Container {
	return dag.Container().
		From("golang:1.22.1").
		WithDirectory("/src", source).
		WithWorkdir("/src")

}
