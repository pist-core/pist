package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"unicode"

	"gopkg.in/urfave/cli.v1"

	"git.taiyue.io/pist/go-pist/cmd/utils"
	"git.taiyue.io/pist/go-pist/crypto"
	"git.taiyue.io/pist/go-pist/dashboard"
	"git.taiyue.io/pist/go-pist/node"
	"git.taiyue.io/pist/go-pist/params"
	"git.taiyue.io/pist/go-pist/pist"
	"github.com/naoina/toml"
)

var (
	dumpConfigCommand = cli.Command{
		Action:      utils.MigrateFlags(dumpConfig),
		Name:        "dumpconfig",
		Usage:       "Show configuration values",
		ArgsUsage:   "",
		Flags:       append(append(nodeFlags, rpcFlags...)),
		Category:    "MISCELLANEOUS COMMANDS",
		Description: `The dumpconfig command shows configuration values.`,
	}

	configFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TOML configuration file",
	}
)

// These settings ensure that TOML keys use the same names as Go struct fields.
var tomlSettings = toml.Config{
	NormFieldName: func(rt reflect.Type, key string) string {
		return key
	},
	FieldToKey: func(rt reflect.Type, field string) string {
		return field
	},
	MissingField: func(rt reflect.Type, field string) error {
		link := ""
		if unicode.IsUpper(rune(rt.Name()[0])) && rt.PkgPath() != "main" {
			link = fmt.Sprintf(", see https://godoc.org/%s#%s for available fields", rt.PkgPath(), rt.Name())
		}
		return fmt.Errorf("field '%s' is not defined in %s%s", field, rt.String(), link)
	},
}

type piststatsConfig struct {
	URL string `toml:",omitempty"`
}

type gethConfig struct {
	Pist      pist.Config
	Node      node.Config
	Piststats piststatsConfig
	Dashboard dashboard.Config
}

func loadConfig(file string, cfg *gethConfig) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tomlSettings.NewDecoder(bufio.NewReader(f)).Decode(cfg)
	// Add file name to errors that have a line number.
	if _, ok := err.(*toml.LineError); ok {
		err = errors.New(file + ", " + err.Error())
	}
	return err
}

func defaultNodeConfig() node.Config {
	cfg := node.DefaultConfig
	cfg.Name = clientIdentifier
	cfg.Version = params.VersionWithCommit(gitCommit, gitDate)
	cfg.HTTPModules = append(cfg.HTTPModules, "pist", "eth", "impawn", "shh")
	cfg.WSModules = append(cfg.WSModules, "pist")
	cfg.IPCPath = "gpist.ipc"
	return cfg
}

func makeConfigNode(ctx *cli.Context) (*node.Node, gethConfig) {
	// Load defaults.
	cfg := gethConfig{
		Pist:      pist.DefaultConfig,
		Node:      defaultNodeConfig(),
		Dashboard: dashboard.DefaultConfig,
	}
	if ctx.GlobalBool(utils.SingleNodeFlag.Name) {
		// set pistconfig
		prikey, _ := crypto.HexToECDSA("c1581e25937d9ab91421a3e1a2667c85b0397c75a195e643109938e987acecfc")
		cfg.Pist.PrivateKey = prikey
		cfg.Pist.CommitteeKey = crypto.FromECDSA(prikey)

		//cfg.Pist.NetworkId =400
		//set node config
		cfg.Node.HTTPPort = 8888
		cfg.Node.HTTPHost = "127.0.0.1"
		cfg.Node.HTTPModules = []string{"db", "pist", "net", "web3", "personal", "admin", "miner", "eth"}

		ctx.GlobalSet("datadir", "./data")
	}
	// Load config file.
	if file := ctx.GlobalString(configFileFlag.Name); file != "" {
		if err := loadConfig(file, &cfg); err != nil {
			utils.Fatalf("%v", err)
		}
	}

	// Apply flags.
	utils.SetNodeConfig(ctx, &cfg.Node)
	stack, err := node.New(&cfg.Node)
	if err != nil {
		utils.Fatalf("Failed to create the protocol stack: %v", err)
	}
	utils.SetPistchainConfig(ctx, stack, &cfg.Pist)
	if ctx.GlobalIsSet(utils.PistStatsURLFlag.Name) {
		cfg.Piststats.URL = ctx.GlobalString(utils.PistStatsURLFlag.Name)
	}

	utils.SetDashboardConfig(ctx, &cfg.Dashboard)

	return stack, cfg
}

func makeFullNode(ctx *cli.Context) *node.Node {
	stack, cfg := makeConfigNode(ctx)

	utils.RegisterPistService(stack, &cfg.Pist)

	if ctx.GlobalBool(utils.DashboardEnabledFlag.Name) {
		utils.RegisterDashboardService(stack, &cfg.Dashboard, gitCommit)
	}

	// Add the Pistchain Stats daemon if requested.
	if cfg.Piststats.URL != "" {
		utils.RegisterPistStatsService(stack, cfg.Piststats.URL)
	}
	return stack
}

// dumpConfig is the dumpconfig command.
func dumpConfig(ctx *cli.Context) error {
	_, cfg := makeConfigNode(ctx)
	comment := ""

	if cfg.Pist.Genesis != nil {
		cfg.Pist.Genesis = nil
		comment += "# Note: this config doesn't contain the genesis block.\n\n"
	}

	out, err := tomlSettings.Marshal(&cfg)
	if err != nil {
		return err
	}
	io.WriteString(os.Stdout, comment)
	os.Stdout.Write(out)
	return nil
}
