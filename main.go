package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/google/rpmpack"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type InputFile struct {
	Path        string
	Destination string
	Mode        os.FileMode
	Owner       string
	Group       string
}

func ParseInputFile(v string) (InputFile, error) {
	parts := strings.Split(v, ":")
	if len(parts) < 2 || len(parts) > 5 {
		return InputFile{}, errors.Errorf("invalid input file: %s", v)
	}

	path := parts[0]
	destination := parts[1]
	mode := os.FileMode(0644)
	owner := "root"
	group := "root"
	if len(parts) >= 3 {
		m, err := strconv.ParseUint(parts[2], 8, 32)
		if err != nil {
			return InputFile{}, errors.Wrapf(err, "invalid input file mode: %s", v)
		}
		mode = os.FileMode(m)
	}
	if len(parts) >= 4 {
		owner = parts[3]
	}
	if len(parts) >= 5 {
		group = parts[4]
	}

	return InputFile{
		Path:        path,
		Destination: destination,
		Mode:        mode,
		Owner:       owner,
		Group:       group,
	}, nil
}

var (
	FlagName = &cli.StringFlag{
		Name:     "name",
		Usage:    "name of the package",
		Required: true,
		EnvVars:  []string{"RPM_ASSEMBLER_NAME"},
	}

	FlagSummary = &cli.StringFlag{
		Name:     "summary",
		Usage:    "summary of the package",
		Required: false,
		EnvVars:  []string{"RPM_ASSEMBLER_SUMMARY"},
	}

	FlagDescription = &cli.StringFlag{
		Name:     "description",
		Usage:    "description of the package",
		Required: false,
		EnvVars:  []string{"RPM_ASSEMBLER_DESCRIPTION"},
	}

	FlagVersion = &cli.StringFlag{
		Name:     "version",
		Usage:    "version of the package",
		Value:    "0.0.0",
		Required: false,
		EnvVars:  []string{"RPM_ASSEMBLER_VERSION"},
	}

	FlagRelease = &cli.StringFlag{
		Name:     "release",
		Usage:    "release of the package",
		Value:    "0",
		Required: false,
		EnvVars:  []string{"RPM_ASSEMBLER_RELEASE"},
	}

	FlagArch = &cli.StringFlag{
		Name:     "arch",
		Usage:    "architecture of the package. this is usually one of: noarch, x86_64, aarch64, armv7hl, i686, ppc64, ppc64le, s390x",
		Value:    "noarch",
		Required: false,
		EnvVars:  []string{"RPM_ASSEMBLER_ARCH"},
	}

	FlagOS = &cli.StringFlag{
		Name:     "os",
		Usage:    "operating system of the package",
		Required: false,
		EnvVars:  []string{"RPM_ASSEMBLER_OS"},
	}

	FlagVendor = &cli.StringFlag{
		Name:     "vendor",
		Usage:    "vendor of the package",
		Required: false,
		EnvVars:  []string{"RPM_ASSEMBLER_VENDOR"},
	}

	FlagURL = &cli.StringFlag{
		Name:     "url",
		Usage:    "url of the package",
		Required: false,
		EnvVars:  []string{"RPM_ASSEMBLER_URL"},
	}

	FlagPackager = &cli.StringFlag{
		Name:     "packager",
		Usage:    "packager of the package",
		Required: false,
		EnvVars:  []string{"RPM_ASSEMBLER_PACKAGER"},
	}

	FlagGroup = &cli.StringFlag{
		Name:     "group",
		Usage:    "group of the package",
		Required: false,
		EnvVars:  []string{"RPM_ASSEMBLER_GROUP"},
	}

	FlagLicence = &cli.StringFlag{
		Name:     "licence",
		Usage:    "licence of the package",
		Required: false,
		EnvVars:  []string{"RPM_ASSEMBLER_LICENCE"},
	}

	FlagEpoch = &cli.IntFlag{
		Name:     "epoch",
		Usage:    "epoch of the package",
		Value:    0,
		Required: false,
		EnvVars:  []string{"RPM_ASSEMBLER_EPOCH"},
	}

	FlagProvides = &cli.StringSliceFlag{
		Name:     "provides",
		Usage:    "provides of the package",
		Required: false,
		EnvVars:  []string{"RPM_ASSEMBLER_PROVIDES"},
	}

	FlagRequires = &cli.StringSliceFlag{
		Name:     "requires",
		Usage:    "requires of the package",
		Required: false,
		EnvVars:  []string{"RPM_ASSEMBLER_REQUIRES"},
	}

	FlagConflicts = &cli.StringSliceFlag{
		Name:     "conflicts",
		Usage:    "conflicts of the package",
		Required: false,
		EnvVars:  []string{"RPM_ASSEMBLER_CONFLICTS"},
	}

	FlagOutput = &cli.StringFlag{
		Name:     "output",
		Usage:    "output file. if not specified, the package will be written to the current working directory",
		Required: false,
		EnvVars:  []string{"RPM_ASSEMBLER_OUTPUT"},
	}
)

func stringSliceToRelations(s []string) ([]*rpmpack.Relation, error) {
	relations := []*rpmpack.Relation{}
	for _, v := range s {
		r, err := rpmpack.NewRelation(v)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse relation: %s", v)
		}
		relations = append(relations, r)
	}
	return relations, nil
}

func action(c *cli.Context) error {
	provides, err := stringSliceToRelations(c.StringSlice(FlagProvides.Name))
	if err != nil {
		return errors.Wrap(err, "failed to parse provides")
	}

	requires, err := stringSliceToRelations(c.StringSlice(FlagRequires.Name))
	if err != nil {
		return errors.Wrap(err, "failed to parse requires")
	}

	conflicts, err := stringSliceToRelations(c.StringSlice(FlagConflicts.Name))
	if err != nil {
		return errors.Wrap(err, "failed to parse conflicts")
	}

	r, err := rpmpack.NewRPM(rpmpack.RPMMetaData{
		Name:        c.String(FlagName.Name),
		Summary:     c.String(FlagSummary.Name),
		Description: c.String(FlagDescription.Name),
		Version:     c.String(FlagVersion.Name),
		Release:     c.String(FlagRelease.Name),
		Arch:        c.String(FlagArch.Name),
		OS:          c.String(FlagOS.Name),
		Vendor:      c.String(FlagVendor.Name),
		URL:         c.String(FlagURL.Name),
		Packager:    c.String(FlagPackager.Name),
		Group:       c.String(FlagGroup.Name),
		Licence:     c.String(FlagLicence.Name),
		Epoch:       uint32(c.Int(FlagEpoch.Name)),
		Provides:    provides,
		Requires:    requires,
		Conflicts:   conflicts,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create rpm")
	}

	for _, v := range c.Args().Slice() {
		inputFile, err := ParseInputFile(v)
		if err != nil {
			return errors.Wrap(err, "failed to parse input file")
		}
		content, err := os.ReadFile(inputFile.Path)
		if err != nil {
			return errors.Wrap(err, "failed to read input file")
		}

		r.AddFile(rpmpack.RPMFile{
			Name:  inputFile.Destination,
			Body:  content,
			Mode:  uint(inputFile.Mode),
			Owner: inputFile.Owner,
			Group: inputFile.Group,
		})
	}

	output := c.String(FlagOutput.Name)
	rpmName := fmt.Sprintf("%s-%s-%s.%s.rpm", r.Name, r.Version, r.Release, r.Arch)
	if info, err := os.Stat(output); (err == nil && info.IsDir()) || (err != nil && os.IsNotExist(err) && strings.Contains(output, "/")) {
		dir := path.Dir(output)
		if path.Base(dir) == "." {
			output = fmt.Sprintf("%s/%s", output, rpmName)
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return errors.Wrap(err, "failed to create output directory")
		}
	} else if output == "" {
		output = rpmName
	} else if err != nil {
		return errors.Wrap(err, "failed to stat output file")
	}

	f, err := os.Create(output)
	if err != nil {
		return errors.Wrapf(err, "failed to create output file: %s", output)
	}
	defer f.Close()
	return r.Write(f)
}

func main() {
	app := &cli.App{
		Name:   "rpm-assembler",
		Usage:  "assemble rpm packages from artifacts",
		Action: action,
		Flags: []cli.Flag{
			FlagName, FlagSummary, FlagDescription, FlagVersion, FlagRelease, FlagArch, FlagOS, FlagVendor, FlagURL, FlagPackager, FlagGroup, FlagLicence, FlagEpoch, FlagProvides, FlagRequires, FlagConflicts, FlagOutput,
		},
		Args:            true,
		ArgsUsage:       "[input files...]\n" + "  input files are specified as: <path>:<destination>[:<mode>[:<owner>[:<group>]]]",
		HideHelpCommand: true,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
