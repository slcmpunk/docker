// +build !windows

package daemon

import (
	"github.com/docker/docker/container"
	"github.com/docker/docker/pkg/archive"
	"github.com/opencontainers/runc/libcontainer/user"
)

func (daemon *Daemon) tarCopyOptions(container *container.Container, noOverwriteDirNonDir bool) (*archive.TarOptions, error) {
	if container.Config.User == "" {
		return daemon.defaultTarCopyOptions(noOverwriteDirNonDir), nil
	}

	// TODO(runcom): maybe use the newer docker/pkg/idtools which supports lookup
	// users by getent also.
	user, err := user.LookupUser(container.Config.User)
	if err != nil {
		return nil, err
	}

	return &archive.TarOptions{
		NoOverwriteDirNonDir: noOverwriteDirNonDir,
		ChownOpts: &archive.TarChownOptions{
			UID: user.Uid,
			GID: user.Gid,
		},
	}, nil
}
