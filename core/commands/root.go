package commands

import (
	"io"
	"strings"

	oldcmds "github.com/ipfs/go-ipfs/commands"
	dag "github.com/ipfs/go-ipfs/core/commands/dag"
	e "github.com/ipfs/go-ipfs/core/commands/e"
	files "github.com/ipfs/go-ipfs/core/commands/files"
	ocmd "github.com/ipfs/go-ipfs/core/commands/object"
	unixfs "github.com/ipfs/go-ipfs/core/commands/unixfs"

	logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	"gx/ipfs/QmYAri1etRFXVrE7bkFSdpAMpoTiSCUTQnKpuDMPA2feCp/go-ipfs-cmds"
	"gx/ipfs/QmadYQbq2fJpaRE3XhpMLH68NNxmWMwfMQy1ntr1cKf7eo/go-ipfs-cmdkit"
)

var log = logging.Logger("core/commands")

const (
	ApiOption = "api"
)

var Root = &cmds.Command{
	Helptext: cmdsutil.HelpText{
		Tagline:  "Global p2p merkle-dag filesystem.",
		Synopsis: "ipfs [--config=<config> | -c] [--debug=<debug> | -D] [--help=<help>] [-h=<h>] [--local=<local> | -L] [--api=<api>] <command> ...",
		Subcommands: `
BASIC COMMANDS
  init          Initialize ipfs local configuration
  add <path>    Add a file to IPFS
  cat <ref>     Show IPFS object data
  get <ref>     Download IPFS objects
  ls <ref>      List links from an object
  refs <ref>    List hashes of links from an object

DATA STRUCTURE COMMANDS
  block         Interact with raw blocks in the datastore
  object        Interact with raw dag nodes
  files         Interact with objects as if they were a unix filesystem
  dag           Interact with IPLD documents (experimental)

ADVANCED COMMANDS
  daemon        Start a long-running daemon process
  mount         Mount an IPFS read-only mountpoint
  resolve       Resolve any type of name
  name          Publish and resolve IPNS names
  key           Create and list IPNS name keypairs
  dns           Resolve DNS links
  pin           Pin objects to local storage
  repo          Manipulate the IPFS repository
  stats         Various operational stats
  ptp           Libp2p stream mounting
  filestore     Manage the filestore (experimental)

NETWORK COMMANDS
  id            Show info about IPFS peers
  bootstrap     Add or remove bootstrap peers
  swarm         Manage connections to the p2p network
  dht           Query the DHT for values or peers
  ping          Measure the latency of a connection
  diag          Print diagnostics

TOOL COMMANDS
  config        Manage configuration
  version       Show ipfs version information
  update        Download and apply go-ipfs updates
  commands      List all available commands

Use 'ipfs <command> --help' to learn more about each command.

ipfs uses a repository in the local file system. By default, the repo is located
at ~/.ipfs. To change the repo location, set the $IPFS_PATH environment variable:

  export IPFS_PATH=/path/to/ipfsrepo

EXIT STATUS

The CLI will exit with one of the following values:

0     Successful execution.
1     Failed executions.
`,
	},
	Options: []cmdsutil.Option{
		cmdsutil.StringOption("config", "c", "Path to the configuration file to use."),
		cmdsutil.BoolOption("debug", "D", "Operate in debug mode.").Default(false),
		cmdsutil.BoolOption("help", "Show the full command help text.").Default(false),
		cmdsutil.BoolOption("h", "Show a short version of the command help text.").Default(false),
		cmdsutil.BoolOption("local", "L", "Run the command locally, instead of using the daemon.").Default(false),
		cmdsutil.StringOption(ApiOption, "Use a specific API instance (defaults to /ip4/127.0.0.1/tcp/5001)"),
	},
}

// commandsDaemonCmd is the "ipfs commands" command for daemon
var CommandsDaemonCmd = CommandsCmd(Root)

var rootSubcommands = map[string]*cmds.Command{
	"add":       AddCmd,
	"block":     BlockCmd,
	"cat":       CatCmd,
	"commands":  CommandsDaemonCmd,
	"get":       GetCmd,
	"filestore": FileStoreCmd,
}

var rootOldSubcommands = map[string]*oldcmds.Command{
	"bootstrap": BootstrapCmd,
	"config":    ConfigCmd,
	"dag":       dag.DagCmd,
	"dht":       DhtCmd,
	"diag":      DiagCmd,
	"dns":       DNSCmd,
	"files":     files.FilesCmd,
	"id":        IDCmd,
	"key":       KeyCmd,
	"log":       LogCmd,
	"ls":        LsCmd,
	"mount":     MountCmd,
	"name":      NameCmd,
	"object":    ocmd.ObjectCmd,
	"pin":       PinCmd,
	"ping":      PingCmd,
	"ptp":       PTPCmd,
	"pubsub":    PubsubCmd,
	"refs":      RefsCmd,
	"repo":      RepoCmd,
	"resolve":   ResolveCmd,
	"stats":     StatsCmd,
	"swarm":     SwarmCmd,
	"tar":       TarCmd,
	"tour":      tourCmd,
	"file":      unixfs.UnixFSCmd,
	"update":    ExternalBinary(),
	"version":   VersionCmd,
	"bitswap":   BitswapCmd,
	"shutdown":  daemonShutdownCmd,
}

// RootRO is the readonly version of Root
var RootRO = &cmds.Command{}

var CommandsDaemonROCmd = CommandsCmd(RootRO)

var RefsROCmd = &oldcmds.Command{}

var rootROSubcommands = map[string]*cmds.Command{
	"commands": CommandsDaemonROCmd,
	"cat":      CatCmd,
	"block": &cmds.Command{
		Subcommands: map[string]*cmds.Command{
			"stat": blockStatCmd,
			"get":  blockGetCmd,
		},
	},
	"get": GetCmd,
}

var rootROOldSubcommands = map[string]*oldcmds.Command{
	"dns": DNSCmd,
	"ls":  LsCmd,
	"name": &oldcmds.Command{
		Subcommands: map[string]*oldcmds.Command{
			"resolve": IpnsCmd,
		},
	},
	"object": &oldcmds.Command{
		Subcommands: map[string]*oldcmds.Command{
			"data":  ocmd.ObjectDataCmd,
			"links": ocmd.ObjectLinksCmd,
			"get":   ocmd.ObjectGetCmd,
			"stat":  ocmd.ObjectStatCmd,
			"patch": ocmd.ObjectPatchCmd,
		},
	},
	"dag": &oldcmds.Command{
		Subcommands: map[string]*oldcmds.Command{
			"get": dag.DagGetCmd,
		},
	},
	"refs":    RefsROCmd,
	"resolve": ResolveCmd,
	"version": VersionCmd,
}

func init() {
	Root.ProcessHelp()
	*RootRO = *Root

	// sanitize readonly refs command
	*RefsROCmd = *RefsCmd
	RefsROCmd.Subcommands = map[string]*oldcmds.Command{}

	Root.OldSubcommands = rootOldSubcommands
	Root.Subcommands = rootSubcommands

	RootRO.OldSubcommands = rootROOldSubcommands
	RootRO.Subcommands = rootROSubcommands
}

type MessageOutput struct {
	Message string
}

func MessageTextMarshaler(res oldcmds.Response) (io.Reader, error) {
	v, err := unwrapOutput(res.Output())
	if err != nil {
		return nil, err
	}

	out, ok := v.(*MessageOutput)
	if !ok {
		return nil, e.TypeErr(out, v)
	}

	return strings.NewReader(out.Message), nil
}
