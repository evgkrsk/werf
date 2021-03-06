package cleaning

import (
	"context"
	"fmt"
	"time"

	"github.com/werf/lockgate"
	"github.com/werf/logboek"
	"github.com/werf/logboek/pkg/style"
	"github.com/werf/logboek/pkg/types"

	"github.com/werf/werf/pkg/stages_manager"
	"github.com/werf/werf/pkg/storage"
	"github.com/werf/werf/pkg/werf"
)

type StagesPurgeOptions struct {
	RmContainersThatUseWerfImages bool
	DryRun                        bool
}

func StagesPurge(ctx context.Context, projectName string, storageLockManager storage.LockManager, stagesManager *stages_manager.StagesManager, options StagesPurgeOptions) error {
	m := newStagesPurgeManager(projectName, stagesManager, options)

	if lock, err := storageLockManager.LockStagesAndImages(ctx, projectName, storage.LockStagesAndImagesOptions{GetOrCreateImagesOnly: false}); err != nil {
		return fmt.Errorf("unable to lock stages and images: %s", err)
	} else {
		defer storageLockManager.Unlock(ctx, lock)
	}

	return logboek.Context(ctx).Default().LogProcess("Running stages purge").
		Options(func(options types.LogProcessOptionsInterface) {
			options.Style(style.Highlight())
		}).
		DoError(func() error {
			return m.run(ctx)
		})
}

func newStagesPurgeManager(projectName string, stagesManager *stages_manager.StagesManager, options StagesPurgeOptions) *stagesPurgeManager {
	return &stagesPurgeManager{
		StagesManager:                 stagesManager,
		ProjectName:                   projectName,
		RmContainersThatUseWerfImages: options.RmContainersThatUseWerfImages,
		DryRun:                        options.DryRun,
	}
}

type stagesPurgeManager struct {
	StagesManager                 *stages_manager.StagesManager
	ProjectName                   string
	RmContainersThatUseWerfImages bool
	DryRun                        bool
}

func (m *stagesPurgeManager) run(ctx context.Context) error {
	deleteImageOptions := storage.DeleteImageOptions{
		RmiForce:                 true,
		SkipUsedImage:            false,
		RmForce:                  m.RmContainersThatUseWerfImages,
		RmContainersThatUseImage: m.RmContainersThatUseWerfImages,
	}

	lockName := fmt.Sprintf("stages-purge.%s", m.ProjectName)
	return werf.WithHostLock(ctx, lockName, lockgate.AcquireOptions{Timeout: time.Second * 600}, func() error {
		logProcess := logboek.Context(ctx).Default().LogProcess("Deleting stages")
		logProcess.Start()

		stages, err := m.StagesManager.GetAllStages(ctx)
		if err != nil {
			logProcess.Fail()
			return err
		}

		if err := deleteStageInStagesStorage(ctx, m.StagesManager, deleteImageOptions, m.DryRun, stages...); err != nil {
			logProcess.Fail()
			return err
		} else {
			logProcess.End()
		}

		logProcess = logboek.Context(ctx).Default().LogProcess("Deleting managed images")
		logProcess.Start()

		managedImages, err := m.StagesManager.StagesStorage.GetManagedImages(ctx, m.ProjectName)
		if err != nil {
			logProcess.Fail()
			return err
		}

		for _, managedImage := range managedImages {
			if !m.DryRun {
				if err := m.StagesManager.StagesStorage.RmManagedImage(ctx, m.ProjectName, managedImage); err != nil {
					return err
				}
			}

			logTag := managedImage
			if logTag == "" {
				logTag = storage.NamelessImageRecordTag
			}

			logboek.Context(ctx).Default().LogFDetails("  tag: %s\n", logTag)
			logboek.Context(ctx).LogOptionalLn()
		}

		logProcess.End()

		return nil
	})
}
