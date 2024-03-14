// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package kernels

import (
	"context"

	"github.com/cilium/little-vm-helper/pkg/kernels"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	archFlag = "arch"
	archHelp = "target architecture to configure the kernel, e.g. 'amd64' or 'arm64' (default to native architecture)"
)

func configureCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "configure <kernel>",
		Short: "configure kernel",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log := logrus.New()
			kd, err := kernels.LoadDir(dirName)
			if err != nil {
				log.Fatal(err)
			}

			kname := args[0]
			if err := kd.ConfigureKernel(context.Background(), log, kname, cmd.Flag(archFlag).Value.String()); err != nil {
				log.Fatal(err)
			}

		},
	}

	cmd.Flags().String(archFlag, "", archHelp)

	return cmd
}

func rawConfigureCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "raw_configure <kernel_dir> [<kernel_name>]",
		Short: "configure a kernel prepared by means other than lvh",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			log := logrus.New()
			kd, err := kernels.LoadDir(dirName)
			if err != nil {
				log.Fatal(err)
			}

			kdir := args[0]
			kname := ""
			if len(args) > 1 {
				kname = args[1]
			}
			if err := kd.RawConfigure(context.Background(), log, kdir, kname, cmd.Flag(archFlag).Value.String()); err != nil {
				log.Fatal(err)
			}

		},
	}

	cmd.Flags().String(archFlag, "", archHelp)

	return cmd
}
