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

type profileManipulationInterface interface {
	getProfile() (*profileFile, error)
	setProfile(profileFile) error
	FileExists(string) bool
	currentProfileFile() string
	currentProfileName() (string, error)
	saveProfile(string, *profileFile) error
	loadProfile(string) (*profileFile, error)
	GetEndpoint() string
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
	ops := struct {
		ProfileOperation
		ProfileFileOps
	}{oper, oper.file_ops}

	initializeProfiles(ops)

	switch oper.opname {
	case "create":
		return createProfile(ops, oper.name)
	case "use":
		return useProfile(ops, oper.name)
	case "show":
		return showProfile(ops, oper.name)
	}
	return error_result("Unknown command")
}

func createProfile(oper profileManipulationInterface, name string) *runResult {
	if name == "" {
		return error_result("Please enter a name")
	}

	curr_name, err := oper.currentProfileName()
	if err != nil {
		return error_result(err.Error())
	}

	profile, err2 := oper.getProfile()
	if err2 != nil {
		return error_result(err2.Error())
	}

	profile.Name = curr_name
	saveError := oper.saveProfile(curr_name, profile)
	if saveError != nil {
		return error_result(saveError.Error())
	}

	newProfile, err3 := oper.getProfile()
	if err3 != nil {
		return error_result(err3.Error())
	}

	newProfile.Name = name
	oper.setProfile(*newProfile)

	return ok_result()
}

func useProfile(oper profileManipulationInterface, name string) *runResult {
	if name == "" {
		return error_result("Please enter a name")
	}

	curr_name, err := oper.currentProfileName()
	if err != nil {
		return error_result(err.Error())
	}

	profile, err := oper.getProfile()
	if err != nil {
		return error_result(err.Error())
	}

	profile.Name = curr_name
	saveError := oper.saveProfile(curr_name, profile)
	if saveError != nil {
		return error_result(saveError.Error())
	}

	loadedProfile, err := oper.loadProfile(name)
	if err != nil {
		return error_result(err.Error())
	}
	oper.setProfile(*loadedProfile)

	return ok_result()
}

func showProfile(oper profileManipulationInterface, name string) *runResult {
	var profile_name, _ = oper.currentProfileName()
	var url = oper.GetEndpoint()

	if name != "" {
		pfile, err := oper.loadProfile(name)
		if err != nil {
			return error_result(err.Error())
		}

		profile_name = pfile.Name
		url = pfile.Login.Endpoint
	}

	fmt.Println("Profile:", profile_name)
	fmt.Println("URL:", url)

	return ok_result()
}

func initializeProfiles(oper profileManipulationInterface) error {
	if oper.FileExists(oper.currentProfileFile()) {
		return nil
	}

	profile, err := oper.getProfile()
	if err != nil {
		return err
	}

	profile.Name = "default"
	err1 := oper.setProfile(*profile)
	if err1 != nil {
		return err1
	}

	return nil
}

func (oper ProfileOperation) currentProfileFile() string {
	return oper.current_profile_file
}

func (oper ProfileOperation) getProfile() (*profileFile, error) {
	var pfile profileFile

	name, err := oper.currentProfileName()
	if err != nil {
		return nil, err
	}

	pfile.Name = name
	pfile.Login.Auth = oper.file_ops.GetAuth()
	pfile.Login.Token = oper.file_ops.GetToken()
	pfile.Login.Endpoint = oper.file_ops.GetEndpoint()

	privateKeyBytes, err1 := oper.file_ops.ReadFile(oper.file_ops.GetPrivateKeyPath())
	if err1 != nil {
		return nil, err1
	}

	publicKeyBytes, err2 := oper.file_ops.ReadFile(oper.file_ops.GetPublicKeyPath())
	if err2 != nil {
		return nil, err2
	}

	certBytes, err3 := oper.file_ops.ReadFile(oper.file_ops.GetCertPath())
	if err3 != nil {
		return nil, err3
	}

	pfile.PrivateKey = string(privateKeyBytes)
	pfile.PublicKey = string(publicKeyBytes)
	pfile.Cert = string(certBytes)

	return &pfile, nil
}

func (oper ProfileOperation) setProfile(profile profileFile) error {
	err := oper.file_ops.WriteFile(oper.current_profile_file, []byte(profile.Name))
	if err != nil {
		return err
	}

	err1 := oper.file_ops.WriteLogin(profile.Login.Auth, profile.Login.Token, profile.Login.Endpoint)
	if err1 != nil {
		return err1
	}

	err2 := oper.file_ops.WriteFile(oper.file_ops.GetPrivateKeyPath(), []byte(profile.PrivateKey))
	if err2 != nil {
		return err2
	}

	err3 := oper.file_ops.WriteFile(oper.file_ops.GetPublicKeyPath(), []byte(profile.PublicKey))
	if err3 != nil {
		return err3
	}

	err4 := oper.file_ops.WriteFile(oper.file_ops.GetCertPath(), []byte(profile.Cert))
	if err4 != nil {
		return err4
	}

	return nil
}

func (oper ProfileOperation) currentProfileName() (string, error) {
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
	return filepath.Join(oper.file_ops.GetConfigDir(), "profile_"+name)
}

func (oper ProfileOperation) profileExists(name string) bool {
	return oper.file_ops.FileExists(oper.profilePath(name))
}

func (oper ProfileOperation) loadProfile(name string) (*profileFile, error) {
	if !oper.profileExists(name) {
		return nil, &profileError{error: "profile " + name + " does not exist"}
	}

	profilePath := oper.profilePath(name)
	var pfile profileFile
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

func (oper ProfileOperation) saveProfile(name string, profile *profileFile) error {
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
