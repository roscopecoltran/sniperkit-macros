package config

import (
    // "errors"
    log "github.com/sirupsen/logrus"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "testing"
    "reflect"
    "strconv"
)

func TestFromNutPackageV9(t *testing.T) {
    log.SetLevel(log.DebugLevel)

    var volume Volume
    volume = &VolumeV9{}
    log.Debug("OK ", volume)

    var baseEnv BaseEnvironment
    baseEnv = &BaseEnvironmentBase{}
    baseEnv = &BaseEnvironmentV9{}
    log.Debug("OK ", baseEnv)

    var conf Config
    // conf = &ConfigBase{}
    // conf = &ProjectBase{}
    // conf = &ProjectBase{}
    conf = NewConfigV9(nil) // make sure that ConfigV9 implements Config
    log.Debug("OK ", conf)

    projV9 := NewProjectV9(nil)
    projV9.Macros = map[string]*MacroV9{
        "run": &MacroV9 {Usage: "run, from project",},
        "build": &MacroV9 {Usage: "build, from project",},
    }
    // *NewConfigV9(
    //         // "ProjectV9.WorkingDir",
    //         // []string{"ProjectV9.Ports"},
    //         nil,
    //     ),
    projV9.ProjectName = "NewNut"
    var proj Project
    proj = projV9
    conf = proj // make sure that Project implements Config
    log.Debug("OK ", proj)

    var macro Macro
    // macro = &MacroBase{}
    macro = &MacroV9{
        Usage: "UsageV9",
        ConfigV9: *NewConfigV9(
            // "MacroV9.WorkingDir",
            // []string{"MacroV9.Ports"},
            proj,
        ),
    }

    conf = macro // make sure that Macro implements Config
    log.Debug("OK ", macro)

    macros := GetMacros(proj)
    log.Debug("OK merge macros ", macros)

    ports := GetPorts(macro)
    log.Debug("OK merge ports ", ports)

    priv := IsPrivileged(macro)
    log.Debug("OK privileged ", priv)

    // macroExec := CreateMacro(proj, "--exec", []string{"echo hello"})
    // log.Debug("OK macroExec ", macroExec)
    log.Debug("proj before YAML ", proj)

//     bytes, err := yaml.Marshal(proj)
//     log.Debug("proj ", err, string(bytes))

//     input := `ports:
// - ProjectV9.Ports
// macros:
//   build:
//     usage: build, from project
//   run:
//     usage: run, from project`
//     err = yaml.Unmarshal([]byte(input), proj)
//     log.Debug("proj ", err, proj)
//     bytes, err = yaml.Marshal(proj)
//     log.Debug("proj ", err, string(bytes))
//     err = yaml.Unmarshal([]byte(input), proj)
//     log.Debug("proj ", err, proj)
}

func makeProjectV9(name string, workingDir string, enableGui string) Project {
    p := NewProjectV9(nil)
    p.ConfigV9.Volumes["main"] = &VolumeV9{
        Host:".",
        Container:"/go/src/project",
    }
    p.Macros["build"] = &MacroV9{
        Usage: "build the project",
        Actions: []string{"go build -o nut"},
    }
    p.ProjectName = name
    p.WorkingDir = workingDir
    p.EnableGUI = enableGui
    return p
}

func makeChildV9(name string, workingDir string, enableGui string) Project {
    return makeProjectV9(name, workingDir, enableGui)
}
func makeParentV9(name string, workingDir string, enableGui string) Project {
    return makeProjectV9(name, workingDir, enableGui)
}

func TestInheritanceEnableGuiV9(t *testing.T) {
    possible := []string{"true", "false", "", "junk"}
    for _, vp := range possible {
        for _, vc := range possible {

            p1 := makeChildV9("child", "child", vc)
            p2 := makeParentV9("parent", "parent", vp)

            p1.setParentProject(p2)

            res := IsGUIEnabled(p1)
            var expectation bool
            if vc == "" {
                expectation = strToBool(vp)
            } else {
                expectation = strToBool(vc)
            }

            if res != expectation {
                t.Error(
                    "For", vc, vp,
                    "expected", expectation,
                    "got", res,
                )
            }
        }
    }
}

func TestInheritanceDockerImageV9(t *testing.T) {
    possible := []string{"", "golang:1.7", "golang:1.7"}
    for _, vp := range possible {
        for _, vc := range possible {
            pc := NewProjectV9(nil)
            pc.DockerImage = vc
            pp := NewProjectV9(nil)
            pp.DockerImage = vp

            pc.setParentProject(pp)

            res := GetDockerImage(pc)
            expectation := "undefined"
            if vc == "" {
                expectation = vp
            } else {
                expectation = vc
            }

            if res != expectation {
                t.Error(
                    "For child=", vc, "parent=", vp,
                    "expected", expectation,
                    "got", res,
                )
            }
        }
    }
}

func TestInheritanceWorkingDirV9(t *testing.T) {
    possible := []string{"here", "", "or here"}
    for _, vp := range possible {
        for _, vc := range possible {

            p1 := makeChildV9("child", vc, "child")
            p2 := makeParentV9("parent", vp, "parent")

            p1.setParentProject(p2)

            res := GetWorkingDir(p1)
            var expectation string
            if vc == "" {
                expectation = vp
            } else {
                expectation = vc
            }

            if res != expectation {
                t.Error(
                    "For", vc, vp,
                    "expected", expectation,
                    "got", res,
                )
            }
        }
    }
}

func TestInheritanceNameV9(t *testing.T) {
    possible := []string{"name1", "", "name2"}
    for _, vp := range possible {
        for _, vc := range possible {

            p1 := makeChildV9(vc, "child", "child")
            p2 := makeParentV9(vp, "parent", "parent")

            p1.setParentProject(p2)

            res := GetProjectName(p1)
            var expectation string
            if vc == "" {
                expectation = vp
            } else {
                expectation = vc
            }

            if res != expectation {
                t.Error(
                    "For", vc, vp,
                    "expected", expectation,
                    "got", res,
                )
            }
        }
    }
}

func makeProjectMacrosV9(macros []string, usage string) Project {
    p := NewProjectV9(nil)
    for _, name := range macros {
        p.Macros[name] = &MacroV9{
            Usage: usage,
            Actions: []string{name},
            ConfigV9: *NewConfigV9(p),
        }
    }
    return p
}


func TestInheritanceMacrosV9(t *testing.T) {
    possible := [][]string{[]string{}, []string{"run"}, []string{"build"}, []string{"build", "run"}}
    for _, vp := range possible {
        for _, vc := range possible {

            pc := makeProjectMacrosV9(vc, "child")
            pp := makeProjectMacrosV9(vp, "parent")

            pc.setParentProject(pp)

            nameOccurences := make(map[string]int)
            for _, name := range append(vc, vp...) {
                nameOccurences[name] += 1
            }

            macrosc := GetMacros(pc)
            macrosp := GetMacros(pp)

            for name, count := range nameOccurences {
                if count == 2 {
                    // it must be in parent and in child
                    // and we expect the one of the child
                    if macrosc[name] == nil ||
                       macrosp[name] == nil ||
                       macrosc[name].getUsage() != "child" {
                        t.Error("Macro", name, "defined in parent and in child",
                            "expected to be overriden by the child")
                    }
                } else if count == 1 {
                    // it comes either from parent or from child
                    if macrosp[name] != nil {
                        if macrosp[name].getUsage() != "parent" {
                            t.Error("Macro", name, "defined in parent",
                                "expected to be overriden by the child")
                        }
                    } else if macrosc[name] != nil{
                        if macrosc[name].getUsage() != "child" {
                            t.Error("Macro", name, "defined in child",
                                "expected to be overriden by the child")
                        }
                    } else {
                        t.Error("Macro", name, "has been ignored.")
                    }
                }
            }
        }
    }
}

func TestParsingV9(t *testing.T) {
    log.SetLevel(log.DebugLevel)

    type Tuple struct {
        file string
        env map[string]string
        ports []string
        docker_image string
        detached bool
        UTSMode string
        NetworkMode string
        Devices map[string]Device
        EnableCurrentUser bool
    }

    nutFiles := []Tuple{}

    nutFiles = append(nutFiles,
Tuple{ file:`
syntax_version: "7"
docker_image: golang:1.7
devices:
  first:
    host_path: "/dev/1"
    container_path: "/dev/1"
    options: "rw"
`,
  // second: /dev/2:/dev/2
ports: []string{},
docker_image: "golang:1.7",
detached: false,
Devices: map[string]Device{
        "first": &DeviceV9{
            Host: "/dev/1",
            Container: "/dev/1",
            Options: "rw",
        },
        // "second": "/dev/2:/dev/2",
    },
},

Tuple{ file:`
syntax_version: "7"
docker_image: golang:1.7
`,
ports: []string{},
docker_image: "golang:1.7",
detached: false,
},

Tuple{ file:`
syntax_version: "7"
docker_image: golang:1.7
`,
ports: []string{},
docker_image: "golang:1.7",
detached: false,
},

Tuple{ file:`
syntax_version: "7"
docker_image: golang:1.7
environment:
  A: 1
  B: 2
uts: host
`,
env: map[string]string{
    "A": "1",
    "B": "2",
}, ports: []string{},
docker_image: "golang:1.7",
detached: false,
UTSMode: "host",
},

Tuple{ file:`
syntax_version: "7"
docker_image: golang:1.7
detached: yes # it is not "true" !!
environment:
  A: 1
  B:
ports:
  - "3000:3000"
  - 100:100
net: none
enable_current_user: true
`,
env: map[string]string{
    "A": "1",
    "B": "",
}, ports: []string{
    "3000:3000",
    "100:100",
},
docker_image: "golang:1.7",
detached: false,
NetworkMode: "none",
EnableCurrentUser: true,
},

Tuple{ file:`
syntax_version: "7"
docker_image: golang:1.7
detached: true
environment:
  A: 1
  B:
ports:
  - "3000:3000"
`,
env: map[string]string{
    "A": "1",
    "B": "",
}, ports: []string{
    "3000:3000",
},
docker_image: "golang:1.7",
detached: true,
},
)

    for index, tuple := range nutFiles {
        byteArray := []byte(tuple.file)
        // project := NewProjectV9(nil)
        project, err := ProjectFromYAML(byteArray)
        assert.NotNil(t, project)
        assert.Nil(t, err)

        assert.Equal(t, GetDockerImage(project), tuple.docker_image,
                "Error with tuple " + strconv.Itoa(index) + ": not same docker_image")

        assert.Equal(t, IsDetached(project), tuple.detached,
                "Error with tuple " + strconv.Itoa(index) + ": not Detached")

        assert.Equal(t, GetUTSMode(project), tuple.UTSMode,
                "Error with tuple " + strconv.Itoa(index) + ": not same UTSMode")

        assert.Equal(t, GetNetworkMode(project), tuple.NetworkMode,
                "Error with tuple " + strconv.Itoa(index) + ": not same NetworkMode")

        assert.Equal(t, IsCurrentUserEnabled(project), tuple.EnableCurrentUser,
                "Error with tuple " + strconv.Itoa(index) + ": not same EnableCurrentUser")

        for name, value := range GetDevices(project) {
            assert.Equal(t, value, tuple.Devices[name],
                "Error with tuple " + strconv.Itoa(index) + ": not same device " + name)
        }

        envVariables := project.getEnvironmentVariables()
        assert.Equal(t, len(reflect.ValueOf(tuple.env).MapKeys()),
            len(reflect.ValueOf(envVariables).MapKeys()),
            "Error with tuple " + strconv.Itoa(index) + ": not same keys")

        for name, value := range envVariables {
            assert.Equal(t, value, tuple.env[name],
                "Error with tuple " + strconv.Itoa(index) + ": not same " + name)
            // execInContainer([]string{"echo $" + name}, project)
            // TODO: automate this test
        }

        ports := project.getPorts()
        require.Equal(t, len(tuple.ports), len(ports),
            "Error with tuple " + strconv.Itoa(index) + ": not same port quantity")
        for key, value := range ports {
            assert.Equal(t, value, tuple.ports[key],
                "Error with tuple " + strconv.Itoa(index) +
                ": not same port " + strconv.Itoa(key))
        }
    }

    // test nginx
    // TODO: automate this test
//     nutFile := `
// syntax_version: "7"
// based_on:
//   docker_image: nginx
// ports:
// #  - "80:80"  # works
//   - "80"  # works
// `
//     project := NewProjectV9()
//     byteArray := []byte(nutFile)
//     assert.Nil(t, project.fromYAML(byteArray))
//     execInContainer([]string{}, project)

//     log.Error("end")
}

