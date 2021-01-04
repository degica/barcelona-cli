package operations

import (
	"github.com/degica/barcelona-cli/config"
	"testing"
)

type MockProfileFileOps struct {
}

func (op MockProfileFileOps) FileExists(path string) bool {
	return true
}

func (op MockProfileFileOps) ReadFile(path string) ([]byte, error) {
	return []byte(""), nil
}

func (op MockProfileFileOps) WriteFile(path string, contents []byte) error {
	return nil
}

func (op MockProfileFileOps) GetConfigDir() string {
	return ""
}

func (op MockProfileFileOps) GetPrivateKeyPath() string {
	return ""
}

func (op MockProfileFileOps) GetPublicKeyPath() string {
	return ""
}

func (op MockProfileFileOps) GetCertPath() string {
	return ""
}

func (op MockProfileFileOps) WriteLogin(auth string, token string, endpoint string) error {
	return nil
}

func (op MockProfileFileOps) GetAuth() string {
	return ""
}

func (op MockProfileFileOps) GetToken() string {
	return ""
}

func (op MockProfileFileOps) GetEndpoint() string {
	return ""
}

type MockProfileManipulation struct {
}

func (_ MockProfileManipulation) getProfile() (*profileFile, error) {
	return nil, nil
}

func (_ MockProfileManipulation) setProfile(profileFile) error {
	return nil
}

func (_ MockProfileManipulation) FileExists(string) bool {
	return false
}

func (_ MockProfileManipulation) currentProfileFile() string {
	return ""
}

func (_ MockProfileManipulation) currentProfileName() (string, error) {
	return "", nil
}

func (_ MockProfileManipulation) saveProfile(string, *profileFile) error {
	return nil
}

func (_ MockProfileManipulation) loadProfile(string) (*profileFile, error) {
	return nil, nil
}

func (_ MockProfileManipulation) GetEndpoint() string {
	return ""
}

// ========================================================
// TestCreateProfile
// ========================================================
//

type MockProfileManipulationForCreateProfile struct {
	MockProfileManipulation
	profiles   map[string]*profileFile
	profileName string
	currentProfile *profileFile
	getProfileError error
}

func (m MockProfileManipulationForCreateProfile) currentProfileName() (string, error) {
	return m.profileName, nil
}

func (m MockProfileManipulationForCreateProfile) getProfile() (*profileFile, error) {
	return m.currentProfile, m.getProfileError
}

func (m *MockProfileManipulationForCreateProfile) saveProfile(name string, pfile *profileFile) error {
	m.profiles[name] = pfile
	return nil
}

func (m *MockProfileManipulationForCreateProfile) setProfile(pfile profileFile) error {
	m.profileName = pfile.Name
	return nil
}

func TestCreateProfileWithNothing(t *testing.T) {
	// This simulates not having a .bcn directory
	oper := &MockProfileManipulationForCreateProfile{
		profiles:   map[string]*profileFile{},
		profileName: "adefault",
		getProfileError: profileError{},
	}

	result := createProfile(oper, "newprofile")

	if result.is_error != true {
		t.Errorf("Expected an error but did not get one")
	}
}

func TestCreateProfileWithExistingProfile(t *testing.T) {
	oper := &MockProfileManipulationForCreateProfile{
		profiles:   map[string]*profileFile{},
		profileName: "adefault",
		currentProfile: &profileFile{},
	}

	result := createProfile(oper, "newprofile")

	if result.is_error == true {
		t.Errorf("Expected no error")
	}

	if oper.profileName != "newprofile" {
		t.Errorf("Expected current profile to be 'newprofile' but was not")
	}

	if oper.profiles["adefault"] == nil {
		t.Errorf("Expected adefault to be saved but was not")
	}
}

// ========================================================
// TestUseProfile
// ========================================================
// use a profile

type MockProfileManipulationForUseProfile struct {
	MockProfileManipulation
	curprofile       string
	savedProfile     string
	loadedProfile    string
	currentProfile   *profileFile
	loadProfileError error
}

func (m MockProfileManipulationForUseProfile) currentProfileName() (string, error) {
	return m.curprofile, nil
}

func (_ MockProfileManipulationForUseProfile) getProfile() (*profileFile, error) {
	return &profileFile{}, nil
}

func (m *MockProfileManipulationForUseProfile) saveProfile(name string, pfile *profileFile) error {
	m.savedProfile = name
	return nil
}

func (m *MockProfileManipulationForUseProfile) setProfile(pfile profileFile) error {
	m.currentProfile = &pfile
	return nil
}

func (m *MockProfileManipulationForUseProfile) loadProfile(name string) (*profileFile, error) {
	m.loadedProfile = name
	pfile := &profileFile{
		Name: name,
	}
	return pfile, m.loadProfileError
}

func TestUseProfile(t *testing.T) {
	oper := &MockProfileManipulationForUseProfile{
		curprofile: "p1",
	}

	result := useProfile(oper, "p2")

	if result.is_error == true {
		t.Errorf("Expected no error")
	}

	if oper.savedProfile != "p1" {
		t.Errorf("Expected p1 profile to be saved but was not.")
	}

	if oper.loadedProfile != "p2" {
		t.Errorf("Expected loaded profile to be p2 but got " + oper.loadedProfile)
	}

	if oper.currentProfile.Name != "p2" {
		t.Errorf("Expected loaded profile to be p2 but got " + oper.currentProfile.Name)
	}
}

func TestUseNonexistentProfile(t *testing.T) {
	oper := &MockProfileManipulationForUseProfile{
		curprofile:       "p1",
		loadProfileError: profileError{},
	}

	result := useProfile(oper, "p2")

	if result.is_error != true {
		t.Errorf("Expected an error")
	}
}

// ========================================================
// ExampleShowProfile
// ========================================================
// show the current profile

type MockProfileManipulationForShowProfile struct {
	MockProfileManipulation

	loadProfileResult *profileFile
	loadProfileError  error
}

func (_ MockProfileManipulationForShowProfile) GetEndpoint() string {
	return "https://default.endpoint"
}

func (_ MockProfileManipulationForShowProfile) currentProfileName() (string, error) {
	return "thedefaultprofile", nil
}

func (m MockProfileManipulationForShowProfile) loadProfile(string) (*profileFile, error) {
	return m.loadProfileResult, m.loadProfileError
}

func Example_showProfile_with_no_param() {

	oper := MockProfileManipulationForShowProfile{}

	showProfile(oper, "")

	// Output:
	// Profile: thedefaultprofile
	// URL: https://default.endpoint
}

func Example_showProfile_with_specified_profile() {

	oper := MockProfileManipulationForShowProfile{
		loadProfileResult: &profileFile{
			Name: "aspecificprofile",
			Login: config.Login{
				Endpoint: "https://specific.endpoint",
			},
		},
	}

	showProfile(oper, "blabla")

	// Output:
	// Profile: aspecificprofile
	// URL: https://specific.endpoint
}

func Example_showProfile_with_specified_nonexistent_profile() {

	oper := MockProfileManipulationForShowProfile{
		loadProfileError: &profileError{
			error: "some error",
		},
	}

	showProfile(oper, "blabla")

	// Output:
}

func TestShowProfileWithNonexistentProfile(t *testing.T) {

	oper := MockProfileManipulationForShowProfile{
		loadProfileError: &profileError{
			error: "some error",
		},
	}

	result := showProfile(oper, "blabla")

	if result.is_error != true {
		t.Errorf("Expected an error returned.")
	}

	if result.message != "some error" {
		t.Errorf("Expected 'some error' to returned, but got: " + result.message)
	}
}

// ========================================================
// TestInitializeProfile
// ========================================================
// initialize profiles if its not initialized

type MockProfileOperationForTestInitializeProfile struct {
	MockProfileManipulation

	getProfileOutput profileFile
	getProfileError  error

	profileThatWasSet *profileFile
	setProfileError   error

	currentProfileNameString string
	currentProfileNameError  error

	fileExistsBool bool
}

func (op MockProfileOperationForTestInitializeProfile) FileExists(path string) bool {
	return op.fileExistsBool
}

func (m MockProfileOperationForTestInitializeProfile) getProfile() (*profileFile, error) {
	return &m.getProfileOutput, m.getProfileError
}

func (m *MockProfileOperationForTestInitializeProfile) setProfile(profile profileFile) error {
	m.profileThatWasSet = &profile
	return m.setProfileError
}

func (m MockProfileOperationForTestInitializeProfile) currentProfileName() (string, error) {
	return m.currentProfileNameString, m.currentProfileNameError
}

func TestInitializeProfileWhenNoFilesExist(t *testing.T) {
	oper := &MockProfileOperationForTestInitializeProfile{
		getProfileError: profileError{
			error: "some error due to file load",
		},
		fileExistsBool: false,
	}

	err := initializeProfiles(oper)

	if err == nil {
		t.Errorf("Expected error but got none")
	}
}

func TestInitializeProfileWhenNoProfilesExist(t *testing.T) {
	oper := &MockProfileOperationForTestInitializeProfile{
		getProfileError: nil,
		fileExistsBool:  false,
	}

	err := initializeProfiles(oper)

	if err != nil {
		t.Errorf("Expected no error but got: " + err.Error())
	}

	if oper.profileThatWasSet == nil {
		t.Errorf("Expected a profile to be set")
	}
}

func TestInitializeProfileWhenProfilesExist(t *testing.T) {
	oper := &MockProfileOperationForTestInitializeProfile{
		getProfileError: profileError{
			error: "should not error",
		},
		fileExistsBool: true,
	}

	err := initializeProfiles(oper)

	if err != nil {
		t.Errorf("Expected no error but got: " + err.Error())
	}

	if oper.profileThatWasSet != nil {
		t.Errorf("Expected no profile to be set")
	}
}

// ========================================================
// TestGetProfile
// ========================================================
// gets the current settings and assigns a name to it

type MockProfileFileOpsForTestGetProfile struct {
	MockProfileFileOps
}

func (op MockProfileFileOpsForTestGetProfile) GetPrivateKeyPath() string {
	return "/private.key"
}

func (op MockProfileFileOpsForTestGetProfile) GetPublicKeyPath() string {
	return "/public.key"
}

func (op MockProfileFileOpsForTestGetProfile) GetCertPath() string {
	return "/cert.cert"
}

func (op MockProfileFileOpsForTestGetProfile) ReadFile(path string) ([]byte, error) {
	if path == "/private.key" {
		return []byte("aprivatekey"), nil
	}

	if path == "/public.key" {
		return []byte("apublickey"), nil
	}

	if path == "/cert.cert" {
		return []byte("acert"), nil
	}

	if path == "/profilename" {
		return []byte("thename"), nil
	}

	return []byte(""), nil
}

func (op MockProfileFileOpsForTestGetProfile) GetAuth() string {
	return "anauth"
}

func (op MockProfileFileOpsForTestGetProfile) GetToken() string {
	return "atoken"
}

func (op MockProfileFileOpsForTestGetProfile) GetEndpoint() string {
	return "anendpoint"
}

func TestGetProfile(t *testing.T) {
	ops := &MockProfileFileOpsForTestGetProfile{}

	oper := &ProfileOperation{
		current_profile_file: "/profilename",
		file_ops:             ops,
	}

	pfile, err := oper.getProfile()

	if err != nil {
		t.Errorf("Did not expect error")
	}

	if pfile.Name != "thename" {
		t.Errorf("Expected 'thename' but got: " + pfile.Name)
	}

	if pfile.PrivateKey != "aprivatekey" {
		t.Errorf("Expected 'aprivatekey' but got: " + pfile.PrivateKey)
	}

	if pfile.PublicKey != "apublickey" {
		t.Errorf("Expected 'apublickey' but got: " + pfile.PublicKey)
	}

	if pfile.Cert != "acert" {
		t.Errorf("Expected 'acert' but got: " + pfile.Cert)
	}

	if pfile.Login.Auth != "anauth" {
		t.Errorf("Expected 'anauth' but got: " + pfile.Login.Auth)
	}

	if pfile.Login.Endpoint != "anendpoint" {
		t.Errorf("Expected 'anendpoint' but got: " + pfile.Login.Endpoint)
	}

	if pfile.Login.Token != "atoken" {
		t.Errorf("Expected 'atoken' but got: " + pfile.Login.Token)
	}
}

// ========================================================
// TestSetProfile
// ========================================================
// sets the current settings given a profile

type MockProfileFileOpsForTestSetProfile struct {
	MockProfileFileOps
	filesContents   map[string]string
	writtenAuth     string
	writtenToken    string
	writtenEndpoint string
}

func (op MockProfileFileOpsForTestSetProfile) GetPrivateKeyPath() string {
	return "/privatekeyfile"
}

func (op MockProfileFileOpsForTestSetProfile) GetPublicKeyPath() string {
	return "/publickeyfile"
}

func (op MockProfileFileOpsForTestSetProfile) GetCertPath() string {
	return "/certfile"
}

func (op *MockProfileFileOpsForTestSetProfile) WriteFile(path string, contents []byte) error {
	op.filesContents[path] = string(contents)
	return nil
}

func (op *MockProfileFileOpsForTestSetProfile) WriteLogin(auth string, token string, endpoint string) error {
	op.writtenToken = token
	op.writtenAuth = auth
	op.writtenEndpoint = endpoint
	return nil
}

func TestSetProfile(t *testing.T) {
	ops := &MockProfileFileOpsForTestSetProfile{
		filesContents: map[string]string{},
	}

	oper := &ProfileOperation{
		current_profile_file: "/profilenamefile",
		file_ops:             ops,
	}

	pfile := profileFile{
		Name: "testo",
		Login: config.Login{
			Auth:     "theauth",
			Token:    "thetoken",
			Endpoint: "https://theendpoint",
		},
		PrivateKey: "theprivatekey",
		PublicKey:  "thepublickey",
		Cert:       "thecert",
	}

	err := oper.setProfile(pfile)

	if err != nil {
		t.Errorf("Did not expect error")
	}

	if ops.filesContents["/privatekeyfile"] != "theprivatekey" {
		t.Errorf("Expected privatekeyfile to contain 'theprivatekey' instead got: " + ops.filesContents["/privatekeyfile"])
	}

	if ops.filesContents["/publickeyfile"] != "thepublickey" {
		t.Errorf("Expected publickeyfile to contain 'thepublickey' instead got: " + ops.filesContents["/publickeyfile"])
	}

	if ops.filesContents["/certfile"] != "thecert" {
		t.Errorf("Expected certfile to contain 'thecert' instead got: " + ops.filesContents["/certfile"])
	}

	if ops.filesContents["/profilenamefile"] != "testo" {
		t.Errorf("Expected profilenamefile to contain 'testo' instead got: " + ops.filesContents["/profilenamefile"])
	}

	if ops.writtenAuth != "theauth" {
		t.Errorf("Expected writtenAuth to contain 'theauth' instead got: " + ops.writtenAuth)
	}

	if ops.writtenToken != "thetoken" {
		t.Errorf("Expected writtenToken to contain 'thetoken' instead got: " + ops.writtenToken)
	}

	if ops.writtenEndpoint != "https://theendpoint" {
		t.Errorf("Expected writtenEndpoint to contain 'https://theendpoint' instead got: " + ops.writtenEndpoint)
	}
}

// ========================================================
// TestCurrentProfileName
// ========================================================

type MockProfileFileOpsForTestCurrentProfileName struct {
	MockProfileFileOps
	profileExists bool
}

func (op *MockProfileFileOpsForTestCurrentProfileName) FileExists(name string) bool {
	return op.profileExists
}

func (op MockProfileFileOpsForTestCurrentProfileName) ReadFile(path string) ([]byte, error) {
	return []byte("profilename"), nil
}

func TestCurrentProfileNameWhenNoneExists(t *testing.T) {
	ops := &MockProfileFileOpsForTestCurrentProfileName{}
	ops.profileExists = false

	oper := &ProfileOperation{
		file_ops: ops,
	}

	name, err := oper.currentProfileName()

	if err != nil {
		t.Errorf("Expected no error but got: " + err.Error())
	}

	if name != "default" {
		t.Errorf("Expected 'default' but got: " + name)
	}
}

func TestCurrentProfileNameWhenOneExists(t *testing.T) {
	ops := &MockProfileFileOpsForTestCurrentProfileName{}
	ops.profileExists = true

	oper := &ProfileOperation{
		file_ops: ops,
	}

	name, err := oper.currentProfileName()

	if err != nil {
		t.Errorf("Expected no error but got: " + err.Error())
	}

	if name != "profilename" {
		t.Errorf("Expected 'profilename' but got: " + name)
	}
}

// ========================================================
// TestProfilePath
// ========================================================

type MockProfileFileOpsForTestProfilePath struct {
	MockProfileFileOps
}

func (op MockProfileFileOpsForTestProfilePath) GetConfigDir() string {
	return "/something"
}

func TestProfilePath(t *testing.T) {
	ops := &MockProfileFileOpsForTestProfilePath{}

	oper := &ProfileOperation{
		file_ops: ops,
	}

	profilePath := oper.profilePath("abc")

	if profilePath != "/something/profile_abc" {
		t.Errorf("Expected '/something/profile_abc' but got: " + profilePath)
	}
}

// ========================================================
// TestProfileExists
// ========================================================

type MockProfileFileOpsForTestProfileExists struct {
	MockProfileFileOps
	checkedFile       string
	assertedExistence bool
}

func (op *MockProfileFileOpsForTestProfileExists) FileExists(name string) bool {
	op.checkedFile = name
	return op.assertedExistence
}

func TestProfileExists(t *testing.T) {
	ops := &MockProfileFileOpsForTestProfileExists{}

	oper := &ProfileOperation{
		file_ops: ops,
	}

	oper.profileExists("foobar")

	if ops.checkedFile != "profile_foobar" {
		t.Errorf("Expected 'profile_foobar' but got: " + ops.checkedFile)
	}
}

func TestProfileExistsReturnsResultOfFileExists(t *testing.T) {
	ops := &MockProfileFileOpsForTestProfileExists{}

	oper := &ProfileOperation{
		file_ops: ops,
	}

	ops.assertedExistence = true
	trueResult := oper.profileExists("foobar")

	if trueResult != true {
		t.Errorf("FileExists returned true but profileExists returned false")
	}

	ops.assertedExistence = false
	falseResult := oper.profileExists("foobar")

	if falseResult != false {
		t.Errorf("FileExists returned false but profileExists returned true")
	}
}

// ========================================================
// TestLoadNonexistentProfile
// ========================================================

type MockProfileFileOpsForTestLoadNonexistentProfile struct {
	MockProfileFileOps
}

func (op MockProfileFileOpsForTestLoadNonexistentProfile) FileExists(path string) bool {
	return false
}

func TestLoadNonexistentProfile(t *testing.T) {
	loader := &MockProfileFileOpsForTestLoadNonexistentProfile{}

	oper := &ProfileOperation{
		file_ops: loader,
	}

	_, err := oper.loadProfile("idontexist")

	if err == nil {
		t.Errorf("Should throw an error")
		return
	}

	if err.Error() != "profile idontexist does not exist" {
		t.Errorf("Shuold throw an error 'profile idontexist does not exist' but got: " + err.Error())
	}
}

// ========================================================
// TestLoadProfile
// ========================================================

type MockProfileFileOpsForTestLoadProfile struct {
	MockProfileFileOps
	readFile string
}

func (op *MockProfileFileOpsForTestLoadProfile) ReadFile(path string) ([]byte, error) {
	op.readFile = path
	return []byte(`{"name":"hello2","login":{"auth":"theauth","token":"thetoken","endpoint":"theendpoint"},"privateKey":"theprivatekey","publicKey":"thepublickey","cert":"thecert"}`), nil
}

func TestLoadProfile(t *testing.T) {
	loader := &MockProfileFileOpsForTestLoadProfile{}

	oper := &ProfileOperation{
		file_ops: loader,
	}

	profile, err := oper.loadProfile("theprofile")

	if err != nil {
		t.Errorf("Threw an error: " + err.Error())
	}

	if loader.readFile != "profile_theprofile" {
		t.Errorf("Read the wrong file: " + loader.readFile)
	}

	if profile.Name != "hello2" {
		t.Errorf("Expected Name to be 'hello2' but got " + profile.Name)
	}

	if profile.PrivateKey != "theprivatekey" {
		t.Errorf("Expected PrivateKey to be 'theprivatekey' but got " + profile.PrivateKey)
	}

	if profile.PublicKey != "thepublickey" {
		t.Errorf("Expected PublicKey to be 'thepublickey' but got " + profile.PublicKey)
	}

	if profile.Cert != "thecert" {
		t.Errorf("Expected Cert to be 'thecert' but got " + profile.Cert)
	}

	if profile.Login.Auth != "theauth" {
		t.Errorf("Expected Auth to be 'theauth' but got " + profile.Login.Auth)
	}

	if profile.Login.Endpoint != "theendpoint" {
		t.Errorf("Expected Endpoint to be 'theendpoint' but got " + profile.Login.Endpoint)
	}

	if profile.Login.Token != "thetoken" {
		t.Errorf("Expected Token to be 'thetoken' but got " + profile.Login.Token)
	}
}

// ========================================================
// TestSaveProfile
// ========================================================

type MockProfileFileOpsForTestSaveProfile struct {
	MockProfileFileOps
	writtenFile  string
	writtenBytes []byte
}

func (op *MockProfileFileOpsForTestSaveProfile) WriteFile(path string, contents []byte) error {
	op.writtenFile = path
	op.writtenBytes = contents
	return nil
}

func TestSaveProfile(t *testing.T) {
	saver := &MockProfileFileOpsForTestSaveProfile{}

	oper := &ProfileOperation{
		file_ops: saver,
	}

	profile := &profileFile{}

	err := oper.saveProfile("hello", profile)

	if err != nil {
		t.Errorf("Threw an error: " + err.Error())
	}

	if saver.writtenBytes == nil {
		t.Errorf("writtenBytes nil")
	}
}

// ========================================================
// TestSaveProfileContents
// ========================================================

type MockProfileFileOpsForTestSaveProfileContents struct {
	MockProfileFileOpsForTestSaveProfile
}

func (op MockProfileFileOpsForTestSaveProfileContents) GetConfigDir() string {
	return "/home/test/bcn/configdir"
}

func TestSaveProfileContents(t *testing.T) {
	saver := &MockProfileFileOpsForTestSaveProfileContents{}

	oper := &ProfileOperation{
		file_ops: saver,
	}

	profile := &profileFile{
		Name: "hello2",
		Login: config.Login{
			Auth:     "theauth",
			Token:    "thetoken",
			Endpoint: "https://theendpoint",
		},
		PrivateKey: "theprivatekey",
		PublicKey:  "thepublickey",
		Cert:       "thecert",
	}

	oper.saveProfile("hello2", profile)

	writtenContent := string(saver.writtenBytes)
	if writtenContent != `{"name":"hello2","login":{"auth":"theauth","token":"thetoken","endpoint":"https://theendpoint"},"privateKey":"theprivatekey","publicKey":"thepublickey","cert":"thecert"}` {
		t.Errorf("writtenbytes was " + writtenContent)
	}

	if saver.writtenFile != "/home/test/bcn/configdir/profile_hello2" {
		t.Errorf("Saved to wrong file: " + saver.writtenFile)
	}
}
