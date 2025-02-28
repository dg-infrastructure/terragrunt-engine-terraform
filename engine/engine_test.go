package engine_test

import (
	"context"
	"testing"

	"github.com/gruntwork-io/terragrunt-engine-terraform/engine"

	"github.com/stretchr/testify/require"

	"os"
	"os/exec"

	tgengine "github.com/gruntwork-io/terragrunt-engine-go/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
)

// MockInitServer is a mock implementation of the InitServer interface
type MockInitServer struct {
	mock.Mock
	Responses []*tgengine.InitResponse
}

func (m *MockInitServer) Send(resp *tgengine.InitResponse) error {
	m.Responses = append(m.Responses, resp)
	return nil
}

func (m *MockInitServer) SetHeader(md metadata.MD) error {
	return nil
}

func (m *MockInitServer) SendHeader(md metadata.MD) error {
	return nil
}

func (m *MockInitServer) SetTrailer(md metadata.MD) {
}

func (m *MockInitServer) Context() context.Context {
	return context.TODO()
}

func (m *MockInitServer) SendMsg(msg interface{}) error {
	return nil
}

func (m *MockInitServer) RecvMsg(msg interface{}) error {
	return nil
}

// MockRunServer is a mock implementation of the RunServer interface
type MockRunServer struct {
	mock.Mock
	Responses []*tgengine.RunResponse
}

func (m *MockRunServer) Send(resp *tgengine.RunResponse) error {
	m.Responses = append(m.Responses, resp)
	return nil
}

func (m *MockRunServer) SetHeader(md metadata.MD) error {
	return nil
}

func (m *MockRunServer) SendHeader(md metadata.MD) error {
	return nil
}

func (m *MockRunServer) SetTrailer(md metadata.MD) {
}

func (m *MockRunServer) Context() context.Context {
	return context.TODO()
}

func (m *MockRunServer) SendMsg(msg interface{}) error {
	return nil
}

func (m *MockRunServer) RecvMsg(msg interface{}) error {
	return nil
}

// MockShutdownServer is a mock implementation of the ShutdownServer interface
type MockShutdownServer struct {
	mock.Mock
	Responses []*tgengine.ShutdownResponse
}

func (m *MockShutdownServer) Send(resp *tgengine.ShutdownResponse) error {
	m.Responses = append(m.Responses, resp)
	return nil
}

func (m *MockShutdownServer) SetHeader(md metadata.MD) error {
	return nil
}

func (m *MockShutdownServer) SendHeader(md metadata.MD) error {
	return nil
}

func (m *MockShutdownServer) SetTrailer(md metadata.MD) {
}

func (m *MockShutdownServer) Context() context.Context {
	return context.TODO()
}

func (m *MockShutdownServer) SendMsg(msg interface{}) error {
	return nil
}

func (m *MockShutdownServer) RecvMsg(msg interface{}) error {
	return nil
}

func TestTerraformEngine_Init(t *testing.T) {
	t.Parallel()
	engine := &engine.TerraformEngine{}
	mockStream := &MockInitServer{}

	err := engine.Init(&tgengine.InitRequest{}, mockStream)
	require.NoError(t, err)
	assert.NotEmpty(t, mockStream.Responses)
	assert.Equal(t, "Terraform Initialization completed\n", mockStream.Responses[0].GetStdout())
}

func TestTerraformEngine_Run(t *testing.T) {
	t.Parallel()
	engine := &engine.TerraformEngine{}
	mockStream := &MockRunServer{}

	cmd := "terraform"
	args := []string{"--help"}
	req := &tgengine.RunRequest{
		Command: cmd,
		Args:    args,
		EnvVars: map[string]string{"FOO": "bar"},
	}
	err := engine.Run(req, mockStream)
	require.NoError(t, err)
	assert.NotEmpty(t, mockStream.Responses)
	// merge stdout from all responses to a string
	var output string
	for _, response := range mockStream.Responses {
		if response.GetStdout() != "" {
			output += response.GetStdout()
		}
	}
	assert.Contains(t, output, "Usage: terraform [global options] <subcommand> [args]")
}

func TestTerraformEngineError(t *testing.T) {
	t.Parallel()
	engine := &engine.TerraformEngine{}
	mockStream := &MockRunServer{}

	cmd := "terraform"
	args := []string{"not-a-valid-command"}
	req := &tgengine.RunRequest{
		Command: cmd,
		Args:    args,
	}
	err := engine.Run(req, mockStream)
	require.NoError(t, err)
	assert.NotEmpty(t, mockStream.Responses)
	// merge stdout from all responses to a string
	var output string

	for _, response := range mockStream.Responses {
		if response.GetStderr() != "" {
			output += response.GetStderr()
		}
	}
	// get status code from last response
	code := mockStream.Responses[len(mockStream.Responses)-1].GetResultCode()
	assert.Contains(t, output, "Terraform has no command named \"not-a-valid-command\"")
	assert.NotEqual(t, 0, code)
}

func TestTerraformEngine_Shutdown(t *testing.T) {
	t.Parallel()
	engine := &engine.TerraformEngine{}
	mockStream := &MockShutdownServer{}

	err := engine.Shutdown(&tgengine.ShutdownRequest{}, mockStream)
	require.NoError(t, err)
	assert.Len(t, mockStream.Responses, 1)
	assert.Equal(t, "Terraform Shutdown completed\n", mockStream.Responses[0].GetStdout())
}

func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	cmd := exec.Command(os.Args[3], os.Args[4:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
