package util

import "github.com/spf13/cobra"

type Options struct {
	Cwd            string
	SearchDir      string
	OutputFile     string
	OutputEndpoint string
	Match          string
	Ignore         string
	Pattern        string
	DryRun         bool
	Debug          bool
}

func InitFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("cwd", "c", ".", "Path working directory, filepaths will be relative to this")
	cmd.Flags().StringP("searchdir", "s", "", "Path directory to search")
	cmd.Flags().StringP("outputfile", "f", "", "Output file")
	cmd.Flags().StringP("outputendpoint", "e", "", "Output endpoint")
	cmd.MarkFlagRequired("searchdir")
	cmd.Flags().StringP("match", "m", `.*\.(zip|cbr|cbz)`, "Regex pattern to match files")
	cmd.Flags().StringP("ignore", "i", "", "Regex pattern to ignore files")
	cmd.Flags().StringP("pattern", "p", `\[(?P<author>.+?)\] (?P<title>[^(.]+)(\((?P<publication>.+)?\))?`, "Regex pattern to extract fields")
	cmd.Flags().Bool("dryrun", false, "Only index, do not submit")
	cmd.Flags().Bool("debug", false, "Enable debug logging")
}

func InitOptions(cmd *cobra.Command) Options {
	return Options{
		Cwd:            cmd.Flag("cwd").Value.String(),
		SearchDir:      cmd.Flag("searchdir").Value.String(),
		OutputFile:     cmd.Flag("outputfile").Value.String(),
		OutputEndpoint: cmd.Flag("outputendpoint").Value.String(),
		Match:          cmd.Flag("match").Value.String(),
		Ignore:         cmd.Flag("ignore").Value.String(),
		Pattern:        cmd.Flag("pattern").Value.String(),
		DryRun:         cmd.Flag("dryrun").Value.String() == "true",
		Debug:          cmd.Flag("debug").Value.String() == "true",
	}
}
