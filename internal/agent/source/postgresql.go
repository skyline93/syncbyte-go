package source

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/skyline93/syncbyte-go/internal/pkg/utils"
)

type PostgreSQL struct {
	Opts *Options
}

func NewPostgreSQL(opts *Options) *PostgreSQL {
	return &PostgreSQL{Opts: opts}
}

func (s *PostgreSQL) Uri() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		s.Opts.User, s.Opts.Password, s.Opts.Server, utils.IntToStr(s.Opts.Port), s.Opts.DBName,
	)
}

func (s *PostgreSQL) Dump(destFile string) error {
	path, filename := filepath.Split(destFile)

	u, err := user.Current()
	if err != nil {
		return err
	}

	c := fmt.Sprintf(
		"chmod -R o+w %s; docker run --rm --network host -v %s:/opt:rw postgres:%s bash -c 'pg_dump %s -Fc -f /opt/%s;chmod -R g+w /opt;chown %s:%s /opt/%s'",
		path, path, s.Opts.Version, s.Uri(), filename, u.Uid, u.Gid, filename,
	)
	_, err = utils.Exec("/bin/sh", "-c", c)
	return err
}
