package operations

import (
	"encoding/json"
	"fmt"
	"path/filepath"
)

type profileError struct {
	error string
}

func (err profileError) Error() string {
	return err.error
}

type ProfileFileOps interface {
	FileExists(path string) bool
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, contents []byte) error

	GetConfigDir() string
	GetPrivateKeyPath() string
	GetPublicKeyPath() string
	GetCertPath() string
	WriteLogin(auth string, token string, endpoint string) error

	GetAuth() string
	GetToken() string
	GetEndpoint() string
}

type ProfileOperation struct {
	opname               string
	name                 string
	file_ops             ProfileFileOps
	current_profile_file string
}


func NewProfileOperation(opname string, name string, file_ops ProfileFileOps) *ProfileOperation {
	return &ProfileOperation{
		opname:               opname,
		name:                 name,
		file_ops:             file_ops,
		current_profile_file: filepath.Join(file_ops.GetConfigDir(), "current_profile"),
	}
}

func (oper ProfileOperation) run() *runResult {
	switch oper.opname {
	case "create":
		return oper.createProfile()
	case "delete":
		return oper.deleteProfile()
	case "use":
		return oper.useProfile()
	case "show":
		return oper.showProfile()
	}
	return error_result("Unknown command")
}

func (oper ProfileOperation) createProfile() *runResult {
	if oper.name == "" {
		return error_result("Please enter a name")
	}

	name, _ := oper.currentProfileName()
	profile := oper.getProfile(name)
	oper.saveProfile(name, &profile)

	newProfile := oper.getProfile(oper.name)
	oper.setProfile(newProfile)

	return ok_result()
}

func (oper ProfileOperation) deleteProfile() *runResult {
	if oper.name == "" {
		return error_result("Please enter a name")
	}
	return ok_result()
}

func (oper ProfileOperation) useProfile() *runResult {
	if oper.name == "" {
		return error_result("Please enter a name")
	}

	profile, err := oper.loadProfile(oper.name)
	if err != nil {
		return error_result(err.Error())
	}
	oper.setProfile(*profile)

	return ok_result()
}

func (oper ProfileOperation) showProfile() *runResult {
	var profile_name, _ = oper.currentProfileName()
	var url = oper.file_ops.GetEndpoint()

	if oper.name != "" {
		pfile, err := oper.loadProfile(oper.name)
		if err != nil {
			return error_result(err.Error())
		}

		profile_name = pfile.name
		url = pfile.login.Endpoint
	}

	fmt.Println("Profile:", profile_name)
	fmt.Println("URL:", url)

	return ok_result()
}

func (oper ProfileOperation) initializeProfiles() {
	profile := oper.getProfile("default")
	oper.setProfile(profile)
}

func (oper ProfileOperation) getProfile(name string) ProfileFile {
	var pfile ProfileFile

	pfile.name = name
	pfile.login.Auth = oper.file_ops.GetAuth()
	pfile.login.Token = oper.file_ops.GetToken()
	pfile.login.Endpoint = oper.file_ops.GetEndpoint()

	privateKeyBytes, _ := oper.file_ops.ReadFile(oper.file_ops.GetPrivateKeyPath())
	publicKeyBytes, _ := oper.file_ops.ReadFile(oper.file_ops.GetPublicKeyPath())
	certBytes, _ := oper.file_ops.ReadFile(oper.file_ops.GetCertPath())

	pfile.privateKey = string(privateKeyBytes)
	pfile.publicKey = string(publicKeyBytes)
	pfile.cert = string(certBytes)

	return pfile
}

func (oper ProfileOperation) setProfile(profile ProfileFile) {
	oper.file_ops.WriteFile(oper.current_profile_file, []byte(profile.name))

	oper.file_ops.WriteLogin(profile.login.Auth, profile.login.Token, profile.login.Endpoint)

	oper.file_ops.WriteFile(oper.file_ops.GetPrivateKeyPath(), []byte(profile.privateKey))
	oper.file_ops.WriteFile(oper.file_ops.GetPublicKeyPath(), []byte(profile.publicKey))
	oper.file_ops.WriteFile(oper.file_ops.GetCertPath(), []byte(profile.cert))
}

func (oper ProfileOperation) currentProfileName() (string, error) {
	if !oper.file_ops.FileExists(oper.current_profile_file) {
		oper.initializeProfiles()
	}

	var profile_name = "default"
	if oper.file_ops.FileExists(oper.current_profile_file) {
		contents, err := oper.file_ops.ReadFile(oper.current_profile_file)
		if err != nil {
			return "", err
		}

		profile_name = string(contents)
	}
	return profile_name, nil
}

func (oper ProfileOperation) profilePath(name string) string {
	return filepath.Join(oper.file_ops.GetConfigDir(), "profile_" + name)
}

func (oper ProfileOperation) profileExists(name string) bool {
	return oper.file_ops.FileExists(oper.profilePath(name))
}

func (oper ProfileOperation) loadProfile(name string) (*ProfileFile, error) {
	if !oper.profileExists(name) {
		return nil, &profileError{error: "profile " + name + " does not exist"}
	}

	profilePath := oper.profilePath(name)
	var pfile ProfileFile
	profileJson, err := oper.file_ops.ReadFile(profilePath)
	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal(profileJson, &pfile)
		if err != nil {
			return nil, err
		}
	}

	return &pfile, nil
}

func (oper ProfileOperation) saveProfile(name string, profile *ProfileFile) (error) {
	b, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	err = oper.file_ops.WriteFile(oper.profilePath(name), b)
	if err != nil {
		return err
	}
	return nil
}
