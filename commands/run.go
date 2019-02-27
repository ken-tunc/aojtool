package commands

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/ken-tunc/aojtool/api"
	"github.com/ken-tunc/aojtool/models"

	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/ken-tunc/aojtool/util"

	"github.com/spf13/cobra"
)

var (
	TimeOutSec int
)

var runCmd = &cobra.Command{
	Use:   "run [problem-id] [source-file]",
	Short: "Run program with sample inputs.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("requires at least two args")
		}
		if TimeOutSec < 1 {
			return fmt.Errorf("invalid timeout seconds: %d", TimeOutSec)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var problemId = args[0]
		var sourceFile = args[1]

		samples, err := getSamples(problemId)
		if err != nil {
			abort(err)
		}

		timeout := time.Duration(TimeOutSec) * time.Second
		runner, err := util.NewCodeRunner(sourceFile, timeout)
		if err != nil {
			abort(err)
		}

		for _, sample := range samples {
			cmd.Printf("[Sample %d]\n", sample.Serial)
			out, err := runner.Run(sample.In)
			out = strings.TrimSpace(out)
			oracle := strings.TrimSpace(sample.Out)
			if err != nil {
				abort(err)
			}
			if out == oracle {
				cmd.Println("Pass.")
			} else {
				cmd.Println("Wrong answer...")
				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(out, oracle, true)
				cmd.Println(dmp.DiffPrettyText(diffs))
				break
			}
		}
	},
}

func init() {
	runCmd.Flags().IntVarP(&TimeOutSec, "timeout", "t", 60, "execution timeout seconds")
	rootCmd.AddCommand(runCmd)
}

func getSamples(problemId string) ([]models.TestCase, error) {
	var samples []models.TestCase

	cachePath := getSamplesSavePath(problemId)
	exist, err := util.Exists(cachePath)
	if err != nil {
		return nil, err
	}

	if exist {
		samples, err = loadSample(problemId)
		if err != nil {
			return nil, err
		}
	} else {
		client, err := api.NewClient()
		if err != nil {
			return nil, err
		}

		ctx := context.Background()
		samples, err = client.Test.FindSamples(ctx, problemId)
		if err != nil {
			return nil, err
		}

		err = saveSamples(samples)
		if err != nil {
			return nil, err
		}
	}

	return samples, nil
}

func getSamplesSavePath(problemId string) string {
	return filepath.Join(sampleCache, problemId)
}

func saveSamples(samples []models.TestCase) error {
	path := getSamplesSavePath(samples[0].ProblemId)
	return util.SaveData(path, &samples)
}

func loadSample(problemId string) ([]models.TestCase, error) {
	var samples []models.TestCase

	path := getSamplesSavePath(problemId)
	if err := util.LoadData(path, &samples); err != nil {
		return nil, err
	}

	return samples, nil
}
