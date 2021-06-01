package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/x1nchen/portainer-cli/client"

	perr "github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/x1nchen/portainer-cli/cache"
)

const (
	NameDataDir = ".portainer-cli"
)

var (
	// portainer host
	Host    string
	Datadir string
	Verbose bool
	store   *cache.Store
	pclient *client.PortainerClient
	manager *Manager
)

func init() {
	// cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&Host, "host", "", "host base url such as http://localhost:9000")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(loginRegistryCmd)
	rootCmd.AddCommand(upgradeCmd)

	// os.UserHomeDir()
}

var rootCmd = &cobra.Command{
	Use:                "portainer-cli",
	Short:              "Portainer CLI",
	Long:               `Work seamlessly with Portainer from the command line.`,
	PersistentPreRunE:  prepare,
	PersistentPostRunE: cleanup,
	SilenceErrors:      true,
	SilenceUsage:       true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rootCmd.Println(err)
		os.Exit(-1)
	}
}

// 获取用户 $HOME 路径
// NOTE 如果没有定义 $HOME，在 macos 上可能会返回 "/"，危险
func prepare(cmd *cobra.Command, args []string) error {
	var err error
	var homePath string
	homePath, err = os.UserHomeDir()
	if err != nil {
		return perr.WithMessage(err, "get user home dir error")
	}
	if homePath == "/" {
		return errors.New("get user home dir error, you should define env $HOME")
	}
	// 数据目录：$HOME/.portainer-cli
	Datadir = path.Join(homePath, NameDataDir)
	// ensure datadir exists
	if err = os.MkdirAll(Datadir, 0755); err != nil {
		return perr.WithMessage(err, fmt.Sprintf("create data dir %s error", Datadir))
	}
	return nil
}

// 关闭本地存储
func cleanup(cmd *cobra.Command, args []string) error {
	var err error
	if store != nil {
		err = store.Close()
		if err != nil {
			return perr.WithMessage(err, "close cache store")
		}
	}
	return nil
}

// 初始化 app 外部依赖
func initUnauthorizedManager(cmd *cobra.Command, args []string) error {
	var err error
	Host = strings.TrimSuffix(Host, "/")
	store, err = cache.NewBoltStore(Datadir, Host)
	if err != nil {
		return err
	}
	pclient = client.NewPortainerClient(Host, "")
	manager = initManager(store, pclient, cmd)

	return nil
}

// 初始化 app 外部依赖，已授权的客户端
func initAuthorizedManager(cmd *cobra.Command, args []string) error {
	var err error
	Host = strings.TrimSuffix(Host, "/")
	store, err = cache.NewBoltStore(Datadir, Host)
	if err != nil {
		return err
	}
	token, err := store.TokenService.GetToken()
	if err != nil {
		return err
	}

	if token == "" {
		return errors.New("token invalid")
	}

	pclient = client.NewPortainerClient(Host, token)
	// TODO validate token
	manager = initManager(store, pclient, cmd)

	return nil
}
