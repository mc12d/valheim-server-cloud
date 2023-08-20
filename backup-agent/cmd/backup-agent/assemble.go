package main

import (
	"backup-agent/internal/adapters/aws"
	"backup-agent/internal/adapters/filesystem"
	"backup-agent/internal/app"
	httpport "backup-agent/internal/ports/http"
)

var (
	config = MustLoadConfigFromEnv()
	logger = app.Zerolog()

	s3Storage = func() aws.S3BackupStorage {
		s3, err := aws.NewS3BackupStorage(config.BucketID, config.ObjectID)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to initialize s3 client")
		}
		return *s3
	}()

	cron = app.BackupCron(app.CronLogger(logger))

	backupJob = app.BackupJob{
		Path:          config.BackupDir,
		Zipper:        filesystem.ZipDirectory,
		BackupStorage: s3Storage,
		Logger:        logger,
	}

	restoreJob = app.RestoreJob{
		Path:          config.BackupDir,
		Unzipper:      filesystem.UnzipToDirectory,
		BackupStorage: s3Storage,
		Logger:        logger,
	}

	server = httpport.Server{
		BackupJob:  backupJob,
		RestoreJob: restoreJob,
		Logger:     logger,
	}
)

func main() {
	logger.Debug().Msgf("Starting with config: %+v", config)
	if _, err := cron.AddJob(config.BackupCron, backupJob); err != nil {
		logger.Fatal().Err(err).Msg("failed to add cron job")
	}

	logger.Info().Msg("restoring latest backup")
	if err := restoreJob.Restore(); err != nil {
		logger.Warn().Err(err).Msg("failed to restore latest backup")
	}

	cron.Start()

	server.Start(config.HttpPort)
}
