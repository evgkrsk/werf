package cleaning

import (
	"context"
	"fmt"

	"github.com/werf/logboek"
	"github.com/werf/logboek/pkg/style"
	"github.com/werf/logboek/pkg/types"

	"github.com/werf/werf/pkg/stages_manager"
	"github.com/werf/werf/pkg/storage"
)

type PurgeOptions struct {
	ImagesPurgeOptions
	StagesPurgeOptions
}

func Purge(ctx context.Context, projectName string, imagesRepo storage.ImagesRepo, storageLockManager storage.LockManager, stagesManager *stages_manager.StagesManager, options PurgeOptions) error {
	m := newPurgeManager(projectName, imagesRepo, stagesManager, options)

	if lock, err := storageLockManager.LockStagesAndImages(ctx, projectName, storage.LockStagesAndImagesOptions{GetOrCreateImagesOnly: false}); err != nil {
		return fmt.Errorf("unable to lock stages and images: %s", err)
	} else {
		defer storageLockManager.Unlock(ctx, lock)
	}

	if err := logboek.Context(ctx).Default().LogProcess("Running images purge").DoError(func() error {
		return m.imagesPurgeManager.run(ctx)
	}); err != nil {
		return err
	}

	if err := logboek.Context(ctx).Default().LogProcess("Running stages purge").
		Options(func(options types.LogProcessOptionsInterface) {
			options.Style(style.Highlight())
		}).
		DoError(func() error {
			return m.stagesPurgeManager.run(ctx)
		}); err != nil {
		return err
	}

	return nil
}

func newPurgeManager(projectName string, imagesRepo storage.ImagesRepo, stagesManager *stages_manager.StagesManager, options PurgeOptions) *purgeManager {
	return &purgeManager{
		imagesPurgeManager: newImagesPurgeManager(imagesRepo, options.ImagesPurgeOptions),
		stagesPurgeManager: newStagesPurgeManager(projectName, stagesManager, options.StagesPurgeOptions),
	}
}

type purgeManager struct {
	*imagesPurgeManager
	*stagesPurgeManager
}
