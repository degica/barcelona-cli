package operations

import (
	"encoding/json"
	"fmt"
	"path/filepath"
)

type ProfileFileOps interface {
	FileExists(path string) bool
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, contents []byte) error
	GetConfigDir() string
	GetLoginEndpoint() string
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
	return ok_result()
}

func (oper ProfileOperation) showProfile() *runResult {
	var profile_name = "default"
	var url = oper.file_ops.GetLoginEndpoint()

	if oper.name == "" {
		if oper.file_ops.FileExists(oper.current_profile_file) {
			contents, err := oper.file_ops.ReadFile(oper.current_profile_file)
			if err != nil {
				return error_result(err.Error())
			}

			profile_name = string(contents)
		}
	} else {
		profile_file := filepath.Join(oper.file_ops.GetConfigDir(), "profile_"+oper.name)
		if oper.file_ops.FileExists(profile_file) {

			var profile ProfileFile
			profileJson, err := oper.file_ops.ReadFile(profile_file)
			if err != nil {
				return error_result(err.Error())
			}
		
			err = json.Unmarshal(profileJson, &profile)
			if err != nil {
				return error_result(err.Error())
			}

			profile_name = profile.name
			url = profile.login.Endpoint

		} else {
			return error_result("Profile does not exist")
		}
	}

	fmt.Println("Profile:", profile_name)
	fmt.Println("URL:", url)

	return ok_result()
}

type profileError struct {
	error string
}

func (err profileError) Error() string {
	return err.error
}

func (oper ProfileOperation) profilePath(name string) string {
	return filepath.Join(oper.file_ops.GetConfigDir(), "profile_" + name)
}

func (oper ProfileOperation) profileExists(name string) bool {
	return oper.file_ops.FileExists(oper.profilePath(name))
}

func (oper ProfileOperation) loadProfile(name string) (*ProfileFile, error) {
	if (!oper.profileExists(name)) {
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
