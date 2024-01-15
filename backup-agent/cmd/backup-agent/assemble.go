package main

import (
	"backup-agent/internal/adapters/aws"
	"backup-agent/internal/adapters/filesystem"
	"backup-agent/internal/app"
	"backup-agent/internal/ports/httpserver"
	"github.com/robfig/cron/v3"
	"regexp"
)

var (
	config          = MustLoadConfigFromEnv()
	logger          = app.Zerolog()
	dirContentRegex = regexp.MustCompile(config.BackupFileRegex)
	cronLauncher    = cron.New(cron.WithLogger(app.CronLogger(logger)))

	s3Storage = func() aws.S3BackupStorage {
		s3, err := aws.NewS3BackupStorage(config.BucketID, config.ObjectID)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to initialize s3 client")
		}
		return *s3
	}()

	backupJob = app.BackupJob{
		Path:          config.BackupDir,
		Zipper:        filesystem.DirectoryRegexZipper(dirContentRegex),
		BackupStorage: s3Storage,
		Logger:        logger,
	}

	restoreJob = app.RestoreJob{
		Path:          config.BackupDir,
		Unzipper:      filesystem.DirectoryUnzipper,
		BackupStorage: s3Storage,
		Logger:        logger,
	}

	server = httpserver.Server{
		BackupJob:  backupJob,
		RestoreJob: restoreJob,
		Logger:     logger,
	}
)

func main() {
	logger.Debug().Msgf("Starting with config: %+v", config)
	if _, err := cronLauncher.AddJob(config.BackupCron, backupJob); err != nil {
		logger.Fatal().Err(err).Msg("failed to add cron job")
	}

	if config.BackupRestoreOnStartup {
		logger.Info().Msg("restoring latest backup")
		if err := restoreJob.Restore(); err != nil {
			logger.Warn().Err(err).Msg("failed to restore latest backup")
		}
	}

	cronLauncher.Start()

	server.Start(config.HttpPort)
}
