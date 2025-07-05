/*
Copyright © 2025 Moinak Dey <moinak.dey8@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"mongo-backup/model"
	"mongo-backup/util"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var cfgpath string
var toDump bool
var toRestore bool

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration file to read before starting backup.",
	Long:  "Configuration to read before starting backup.",
	Run:   execute,
}

// Method to call the right config method to read necessary details
// and start with the backup process
func execute(cmd *cobra.Command, args []string) {
	// get the config type
	cfgpathSplits := strings.Split(cfgpath, ".")
	ext := cfgpathSplits[len(cfgpathSplits)-1]

	log.Printf("Provided config file is of type %s", strings.ToUpper(ext))

	switch strings.ToLower(ext) {
	case "json":
		cfg := getjsoncfg()

		log.Printf("Found %d clusters", len(cfg.Confs))

		cfgcount := 0
		for _, conf := range cfg.Confs {
			cfgcount++
			if cfg == nil {
				continue
			}
			log.Printf("Following cluster conf number %d", cfgcount)
			if toDump {
				for _, dumpCfg := range conf.DumpCfgs {
					dump(dumpCfg.DBName, conf.ClusterURI, dumpCfg.DownDir)
				}
			}

			if toRestore {
				for _, resCfg := range conf.RestoreCfgs {
					restore(resCfg.DBName, conf.ClusterURI, resCfg.UpDir)
				}
			}
		}

		return
	default:
		log.Println("Unsupported config file.")
		return
	}
}

// Method to handle db dump
func dump(dbname string, mongouri string, outdir string) {
	log.Println("Starting DB dump...")
	log.Printf("DB name - %s; Out Directory - %s", dbname, outdir)

	// Check if the provided out dir is actually a dir
	if !util.IsDir(outdir) {
		log.Println("Out dir dosen't exist. Creating one...")
		if err := os.MkdirAll(outdir, os.ModePerm); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	// Create the final directory over the out dir
	today := fmt.Sprintf("%d-%s-%d", time.Now().Day(), time.Now().Month().String(), time.Now().Year())
	finalDir := path.Join(outdir, today)

	// have a check if the final dir already exists
	if !util.IsDir(finalDir) {
		log.Println("Final dir dosen't exist. Creating one...")
		if err := os.Mkdir(finalDir, os.ModePerm); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	// Run the mongo dump command
	log.Println("Executing mongodump...")
	out, err := exec.Command("mongodump", "--uri", mongouri, "--db", dbname, "--out", finalDir).Output()
	if err != nil {
		log.Printf("Error while executing mongodump command - %v", err)
		os.Exit(1)
	}

	log.Println("mongodump executed successfully.")
	log.Println(string(out))
}

// Method to handle db restore
func restore(dbname string, mongouri string, updir string) {
	log.Println("Starting DB restore...")
	log.Printf("DB name - %s; Upload Directory - %s", dbname, updir)

	// Check if the provided upload dir is actually a dir
	if !util.IsDir(updir) {
		log.Println("Upload dir dosen't exist.")
		os.Exit(1)
	}

	// Create the final directory over the upload dir
	finalDir := path.Join(updir, dbname)

	// have a check if the final dir exists
	if !util.IsDir(finalDir) {
		log.Println("DB dir dosen't exist.")
		os.Exit(1)
	}

	// Run the mongo restore command
	log.Println("Executing mongorestore...")
	out, err := exec.Command("mongorestore", "--uri", mongouri, "--db", dbname, "--dir", finalDir).Output()
	if err != nil {
		log.Printf("Error while executing mongorestore command - %v", err)
		os.Exit(1)
	}

	log.Println("mongorestore executed successfully.")
	log.Println(string(out))
}

// Method to read json config and return the necessary details
func getjsoncfg() *model.JsonCfg {
	log.Println("Getting JSON configuration.")
	cfg := &model.JsonCfg{}

	bytes, err := os.ReadFile(cfgpath)
	if err != nil {
		log.Println(err)
		return nil
	}

	// Unmarshal the json bytes
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		log.Println(err)
		return nil
	}

	return cfg
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	configCmd.Flags().StringVarP(&cfgpath, "cfg-path", "p", "", "Absolute/Full path of the config file")
	configCmd.Flags().BoolVarP(&toDump, "db-dump", "d", false, "Dump a DB")
	configCmd.Flags().BoolVarP(&toRestore, "db-restore", "r", false, "Restore a DB")

	configCmd.MarkFlagRequired("cfg-path")
}
