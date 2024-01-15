package app

import (
	"fmt"
	"github.com/rs/zerolog"
)

type BackupStorage interface {
	GetLatest() ([]byte, error)
	Get(version string) ([]byte, error)
	Put(obj []byte) error
}

type BackupZipper func(path string) ([]byte, error)
type BackupUnzipper func(path string, zip []byte) error

type BackupJob struct {
	Path   string
	Zipper BackupZipper

	BackupStorage BackupStorage
	Logger        zerolog.Logger
}

func (j BackupJob) Run() {
	err := j.Backup()
	if err != nil {
		j.Logger.Error().Err(err).Msgf("backup failed, filesystem path: %s", j.Path)
	}
}

func (j BackupJob) Backup() error {
	obj, err := j.Zipper(j.Path)
	if err != nil {
		return fmt.Errorf("error reading backup from filesystem [path: %s]: %w", j.Path, err)
	}

	err = j.BackupStorage.Put(obj)
	if err != nil {
		return fmt.Errorf("error uploading backup to remote: %w", err)
	}
	return nil
}

type RestoreJob struct {
	Path     string
	Unzipper BackupUnzipper

	BackupStorage BackupStorage
	Logger        zerolog.Logger
}

func (j RestoreJob) Run() {
	err := j.Restore()
	if err != nil {
		j.Logger.Error().Err(err).Msg("")
	}
}

func (j RestoreJob) RestoreVersion(backupVersion string) error {
	obj, err := j.BackupStorage.Get(backupVersion)
	if err != nil {
		return fmt.Errorf("error downloading backup from remote: %w", err)
	}
	err = j.Unzipper(j.Path, obj)
	if err != nil {
		return fmt.Errorf("error writing backup to filesystem [path: %s]: %w", j.Path, err)
	}
	return nil
}

func (j RestoreJob) Restore() error {
	obj, err := j.BackupStorage.GetLatest()
	if err != nil {
		return fmt.Errorf("error downloading backup from remote: %w", err)
	}
	err = j.Unzipper(j.Path, obj)
	if err != nil {
		return fmt.Errorf("error writing backup to filesystem [path: %s]: %w", j.Path, err)
	}
	return nil
}
