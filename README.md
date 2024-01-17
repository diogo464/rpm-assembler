# rpm-assembler

```
NAME:
   rpm-assembler - assemble rpm packages from artifacts

USAGE:
   rpm-assembler [global options] [input files...]
  input files are specified as: <path>:<destination>[:<mode>[:<owner>[:<group>]]]

GLOBAL OPTIONS:
   --name value                             name of the package [$RPM_ASSEMBLER_NAME]
   --summary value                          summary of the package [$RPM_ASSEMBLER_SUMMARY]
   --description value                      description of the package [$RPM_ASSEMBLER_DESCRIPTION]
   --version value                          version of the package (default: "0.0.0") [$RPM_ASSEMBLER_VERSION]
   --release value                          release of the package (default: "0") [$RPM_ASSEMBLER_RELEASE]
   --arch value                             architecture of the package. this is usually one of: noarch, x86_64, aarch64, armv7hl, i686, ppc64, ppc64le, s390x (default: "noarch") [$RPM_ASSEMBLER_ARCH]
   --os value                               operating system of the package [$RPM_ASSEMBLER_OS]
   --vendor value                           vendor of the package [$RPM_ASSEMBLER_VENDOR]
   --url value                              url of the package [$RPM_ASSEMBLER_URL]
   --packager value                         packager of the package [$RPM_ASSEMBLER_PACKAGER]
   --group value                            group of the package [$RPM_ASSEMBLER_GROUP]
   --licence value                          licence of the package [$RPM_ASSEMBLER_LICENCE]
   --epoch value                            epoch of the package (default: 0) [$RPM_ASSEMBLER_EPOCH]
   --provides value [ --provides value ]    provides of the package [$RPM_ASSEMBLER_PROVIDES]
   --requires value [ --requires value ]    requires of the package [$RPM_ASSEMBLER_REQUIRES]
   --conflicts value [ --conflicts value ]  conflicts of the package [$RPM_ASSEMBLER_CONFLICTS]
   --output value                           output file. if not specified, the package will be written to the current working directory [$RPM_ASSEMBLER_OUTPUT]
   --help, -h                               show help
```
