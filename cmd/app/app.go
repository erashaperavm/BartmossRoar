package app

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/erashaperavm/BartmossRoar/core/api/configs"
	"github.com/erashaperavm/BartmossRoar/core/database"
	"github.com/erashaperavm/BartmossRoar/core/transport"
	"github.com/erashaperavm/BartmossRoar/core/utils"
	"github.com/spf13/cobra"
)

type App struct {
	cfg Config
	db  database.DB
	wd  string
	apiConfig configs.ApiConfig
	ts  []transport.Transport
}

func NewApp() *App {
	wd, err := os.Getwd()
	if err != nil {
		panic("can not get workdir")
	}

	// Load api config, br config
	brPath := path.Join(wd, "br.json")
	brCfg, err := ReadBrConfig(brPath)
	if err != nil {
		panic(err)
	}

	apiPath := path.Join(wd, "api_config.yaml")
	apiCfg, err := configs.ReadApiConfig(apiPath)
	if err != nil {
		panic(err)
	}

	var ts []transport.Transport
	for _, proto := range brCfg.Protocols {
		switch proto {
		case transport.ProtoToString(transport.TCP):
			for _, addr := range brCfg.AddressesIfTCP {
				t, err := transport.NewTCP(addr)
				if err != nil {
					panic(err)
				}
				ts = append(ts, t)
			}
		case transport.ProtoToString(transport.Tor):
			for _, port := range brCfg.LocalPortsIfTor {
				t, err := transport.NewTor(port, false) // do not use existing tor hidden service for now
				if err != nil {
					panic(err)
				}
				ts = append(ts, t)
			}
		default:
			panic("unsupported network protocol")
		}
	}

	return &App{
		wd:        wd,
		cfg:       brCfg,
		apiConfig: apiCfg,
		ts:        ts,
		// todo
	}
}

func (a *App) Start() error {
	// check there is badger db folder
	ok, err := utils.IsBadgerDir(a.wd)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%s has no valid badger db. ", a.wd)
	}

	rootCmd := &cobra.Command{
		Use:   "br",
		Short: "BartmossRoar Core CLI tool. ",
		Long:  "Hidden your identity, chat free. ",
	}

	rootCmd.AddCommand(a.joinCmd())
	rootCmd.AddCommand(startCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	return nil
}

func (a *App) joinCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "join",
		Short:   "join a group",
		Long:    "join a bartmoss roar group",
		Example: "br join network.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			// check if the ntwk is joined
			networkID := args[0]
			strings.Trim(networkID, ".json")
			pathToNtwk := path.Join(a.wd, "networks", networkID)
			if err := os.MkdirAll(filepath.Dir(pathToNtwk), 0750); err != nil {
				return fmt.Errorf("create parent dir: %w", err)
			}
			err := os.Mkdir(pathToNtwk, 0750)
			switch {
			case err == nil:
				// create dir successfully
			case errors.Is(err, fs.ErrExist):
				return errors.New("the network is joined")
			default:
				return fmt.Errorf("create network dir: %w", err)
			}

			// randomly create connection with given nodes first, if nobody response
			// it indicates a broken network or out-date network
			// todo

			return nil
		},
	}
}

func startCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start the node",
		Long: `Start the BartmossRoar node.
			The node will:
				1. Discover neighbor nodes in the channel
				2. Sync and rotate ML-KEM public keys
				3. Tick dummy packets at 3-7s intervals
				4. Handle send, gossip and receive`,
		Example: "br start",
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
}
