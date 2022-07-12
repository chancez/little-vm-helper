package images

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/cilium/little-vm-helper/pkg/images"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func BuildCmd() *cobra.Command {
	var dirName string
	var forceRebuild, dryRun bool

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build VM images",
		Run: func(cmd *cobra.Command, _ []string) {
			log := logrus.New()
			configFname := path.Join(dirName, images.DefaultConfFile)

			configData, err := os.ReadFile(configFname)
			if err != nil {
				log.Fatal(err)
			}

			var cnf images.ImagesConf
			cnf.Dir = dirName
			err = json.Unmarshal(configData, &cnf.Images)
			if err != nil {
				log.Fatal(err)
			}

			forest, err := images.NewImageForest(&cnf, false)
			if err != nil {
				log.Fatal(err)
			}

			start := time.Now()
			res := forest.BuildAllImages(&images.BuildConf{
				Log:          log,
				DryRun:       dryRun,
				ForceRebuild: forceRebuild,
			})
			elapsed := time.Since(start)

			if err := res.Err(); err != nil {
				log.WithError(err).Error("building images failed")
			} else {
				log.WithField("time-elapsed", elapsed).Info("images built succesfully")
			}

			for img, ir := range res.ImageResults {
				if ir.Error == nil {
					fmt.Printf("image:%-10s cachedImageUsed:%t cachedImageDeleted:%s\n", img, ir.CachedImageUsed, ir.CachedImageDeleted)
				}
			}
		},
	}

	cmd.Flags().StringVar(&dirName, "dir", "", "directory to keep the images (configuration will be saved in <dir>/images.json and images in <dir>/images)")
	cmd.MarkFlagRequired("dir")
	cmd.Flags().BoolVar(&forceRebuild, "force-rebuild", false, "rebuild all images, even if they exist")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "do the whole thing, but instead of building actual images create empty files")
	return cmd
}
