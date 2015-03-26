package client

import (
	"fmt"
	"io"
	"strings"

	"github.com/flynn/flynn/Godeps/_workspace/src/github.com/docker/docker/engine"
	"github.com/flynn/flynn/Godeps/_workspace/src/github.com/docker/docker/pkg/archive"
	flag "github.com/flynn/flynn/Godeps/_workspace/src/github.com/docker/docker/pkg/mflag"
	"github.com/flynn/flynn/Godeps/_workspace/src/github.com/docker/docker/utils"
)

func (cli *DockerCli) CmdCp(args ...string) error {
	cmd := cli.Subcmd("cp", "CONTAINER:PATH HOSTDIR|-", "Copy files/folders from a PATH on the container to a HOSTDIR on the host\nrunning the command. Use '-' to write the data\nas a tar file to STDOUT.", true)
	cmd.Require(flag.Exact, 2)

	utils.ParseFlags(cmd, args, true)

	var copyData engine.Env
	info := strings.Split(cmd.Arg(0), ":")

	if len(info) != 2 {
		return fmt.Errorf("Error: Path not specified")
	}

	copyData.Set("Resource", info[1])
	copyData.Set("HostPath", cmd.Arg(1))

	stream, statusCode, err := cli.call("POST", "/containers/"+info[0]+"/copy", copyData, false)
	if stream != nil {
		defer stream.Close()
	}
	if statusCode == 404 {
		return fmt.Errorf("No such container: %v", info[0])
	}
	if err != nil {
		return err
	}

	if statusCode == 200 {
		dest := copyData.Get("HostPath")

		if dest == "-" {
			_, err = io.Copy(cli.out, stream)
		} else {
			err = archive.Untar(stream, dest, &archive.TarOptions{NoLchown: true})
		}
		if err != nil {
			return err
		}
	}
	return nil
}
