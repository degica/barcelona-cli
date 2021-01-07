package operations

import (
	"io"
)

type MockSshcmdOperationApiClient struct {
}

func (m MockSshcmdOperationApiClient) Post(path string, body io.Reader) ([]byte, error) {
	return nil, nil
}

type MockSshcmdOperationConfig struct {
}

func (m MockSshcmdOperationConfig) GetCertPath() string {
	return ""
}

func (m MockSshcmdOperationConfig) GetPrivateKeyPath() string {
	return ""
}

func (m MockSshcmdOperationConfig) IsDebug() bool {
	return false
}

type MockSshcmdOperationCommandRunner struct {
}

func (m MockSshcmdOperationCommandRunner) RunCommand(name string, arg ...string) error {
	return nil
}

func ExampleSshcmdOperation_run_output() {
	client := &MockSshcmdOperationApiClient{}
	mockConfig := &MockSshcmdOperationConfig{}
	mockCmdRunner := &MockSshcmdOperationCommandRunner{}
	oper := NewSshcmdOperation(client, "asd", "123.123.123.123", mockConfig, mockCmdRunner)

	oper.run()
	// Output:
}
