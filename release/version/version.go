package version

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/alexejk/go-release-tools/config"
	"github.com/coreos/go-semver/semver"
	log "github.com/sirupsen/logrus"
)

const DevelopmentVersionPreRelease = "dev"

type Handler struct {
	versionFile     string
	versionProperty string

	versionStringCache string
}

func NewVersionHandler(workDir string) *Handler {

	v := &Handler{
		versionFile:     getVersionFile(workDir),
		versionProperty: config.GetString(config.ProjectVersionProperty),
	}

	_, err := os.Stat(v.versionFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	return v
}

func (v *Handler) GetVersion() (*semver.Version, error) {

	verStr, err := v.readVersionStringFromFile()
	if err != nil {
		return nil, err
	}

	ver, err := semver.NewVersion(verStr)
	if err != nil {
		return nil, err
	}

	return ver, nil
}

func (v *Handler) SetVersion(version *semver.Version) error {

	// Write to file
	return v.writeVersionStringToFile(version.String())
}

func (v *Handler) ReleaseVersion(version *semver.Version) *semver.Version {

	newVersion := &semver.Version{}
	*newVersion = *version

	newVersion.PreRelease = ""
	newVersion.Metadata = ""

	return newVersion
}

func (v *Handler) NextDevelopmentVersion(version *semver.Version) *semver.Version {

	incrementType := config.GetString(config.ProjectVersionIncrementType)
	if incrementType == "" {
		incrementType = "patch"
	}

	newVersion := &semver.Version{}
	*newVersion = *version

	switch strings.ToLower(incrementType) {
	case "major":
		newVersion.BumpMajor()
	case "minor":
		newVersion.BumpMinor()
	case "patch":
		newVersion.BumpPatch()
	default:
		log.Errorf("Unknown version increment type '%s', Doing patch-increment", incrementType)
		newVersion.BumpPatch()
	}

	newVersion.PreRelease = DevelopmentVersionPreRelease
	newVersion.Metadata = ""

	return newVersion
}

func (v *Handler) versionRegexp() *regexp.Regexp {
	return regexp.MustCompile(`(?m:^(\s*` + v.versionProperty + `\s*=\s*)([a-zA-Z0-9-.]*)(\s*)$)`)
}

func (v *Handler) readVersionFile() (string, error) {

	fileBytes, err := ioutil.ReadFile(v.versionFile)
	if err != nil {
		return "", fmt.Errorf("cannot read version file: %s", err.Error())
	}

	return string(fileBytes), nil
}

func (v *Handler) readVersionStringFromFile() (string, error) {

	if v.versionStringCache != "" {
		return v.versionStringCache, nil
	}

	log.Debug("Parsing version file...")

	versionFile, err := v.readVersionFile()
	if err != nil {
		return "", err
	}

	captureGroups := v.versionRegexp().FindStringSubmatch(versionFile)

	if captureGroups != nil {
		return captureGroups[2], nil
	}

	return "", errors.New("unable to find matching version property in the version file")
}

func (v *Handler) writeVersionStringToFile(newVersion string) error {

	// Wiping cache
	v.versionStringCache = ""

	fileInfo, _ := os.Stat(v.versionFile)

	versionFile, err := v.readVersionFile()
	if err != nil {
		return err
	}

	newVersionFile := v.versionRegexp().ReplaceAllString(versionFile, "${1}"+newVersion+"${3}")

	if err = ioutil.WriteFile(v.versionFile, []byte(newVersionFile), fileInfo.Mode().Perm()); err != nil {
		return err
	}

	// Re-write the cache
	v.versionStringCache = newVersion
	return nil
}

func getVersionFile(workDir string) string {

	return path.Join(workDir, config.GetString(config.ProjectVersionFile))
}

func (v *Handler) InterpolateVersionInString(input string) string {

	currentVersion, _ := v.GetVersion()

	res := strings.Replace(input, "${version}", currentVersion.String(), -1)

	return res
}
