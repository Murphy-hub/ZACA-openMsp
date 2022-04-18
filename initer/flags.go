package initer

import (
	"flag"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	cfssl_config "gitlab.oneitfarm.com/bifrost/cfssl/config"

	"gitlab.oneitfarm.com/bifrost/capitalizone/core"
	"gitlab.oneitfarm.com/bifrost/capitalizone/core/config"
	"gitlab.oneitfarm.com/bifrost/capitalizone/pkg/influxdb"
)

const (
	G_                = "IS"
	ConfName          = "conf"
	CmdlineEnvDefault = "default"
	CmdlineEnvTest    = "test"
	CmdlineEnvProd    = "prod"
)

var (
	flagEnv     = flag.String("env", CmdlineEnvProd, "")
	flagEnvfile = flag.String("envfile", ".env", "")
	flagRootca  = flag.Bool("rootca", false, "")
)

func parseConfigs(c *cli.Context) (core.Config, error) {
	// Cmdline flags
	flag.Parse()
	// ENV 读取
	godotenv.Load(*flagEnvfile)
	// Default config
	viper.SetConfigName(fmt.Sprintf("%v.%v", ConfName, CmdlineEnvDefault))
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return core.Config{}, err
	}

	// Read ENV from cmdline
	env := *flagEnv
	// Merge ENV configs
	if env == CmdlineEnvTest || os.Getenv(G_+"_ENV") == CmdlineEnvTest {
		// test ENV
		viper.SetConfigName(fmt.Sprintf("%v.%v", ConfName, CmdlineEnvTest))
	} else {
		// prod ENV
		viper.SetConfigName(fmt.Sprintf("%v.%v", ConfName, CmdlineEnvProd))
	}
	viper.AddConfigPath(".")
	if err := viper.MergeInConfig(); err != nil {
		return core.Config{}, err
	}

	// Merge config frm ENV
	viper.SetEnvPrefix(G_)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	hs, err := os.Hostname()
	if err != nil {
		return core.Config{}, err
	}
	hostname := hs
	if v := os.Getenv("HOSTNAME"); v != "" {
		hostname = v
	}

	conf := core.Config{IConfig: config.IConfig{
		Log: config.Log{
			LogProxy: config.LogProxy{
				Host: viper.GetString("log.log-proxy.host"),
				Port: viper.GetInt("log.log-proxy.port"),
				Key:  viper.GetString("log.log-proxy.key"),
			},
		},
		Registry: config.Registry{
			SelfName: viper.GetString("registry.self-name") + "." + hs,
			Command:  c.Args().First(),
		},
		Redis: config.Redis{
			Nodes: viper.GetStringSlice("redis.nodes"),
		},
		Keymanager: config.Keymanager{
			UpperCa:  viper.GetStringSlice("keymanager.upper-ca"),
			SelfSign: viper.GetBool("keymanager.self-sign"),
			CsrTemplates: config.CsrTemplates{
				RootCa: config.RootCa{
					O:      viper.GetString("keymanager.csr-templates.root-ca.o"),
					Expiry: viper.GetString("keymanager.csr-templates.root-ca.expiry"),
				},
				IntermediateCa: config.IntermediateCa{
					O:      viper.GetString("keymanager.csr-templates.intermediate-ca.o"),
					Ou:     viper.GetString("keymanager.csr-templates.intermediate-ca.ou"),
					Expiry: viper.GetString("keymanager.csr-templates.intermediate-ca.expiry"),
				},
			},
		},
		Singleca: config.Singleca{
			ConfigPath: viper.GetString("singleca.config-path"),
		},
		Election: config.Election{
			Enabled:      viper.GetBool("election.enabled"),
			ID:           viper.GetString("election.id"),
			Baseon:       viper.GetString("election.baseon"),
			AlwaysLeader: viper.GetBool("election.always-leader"),
		},
		GatewayNervs: config.GatewayNervs{
			Enabled:  viper.GetBool("gateway-nervs.enabled"),
			Endpoint: viper.GetString("gateway-nervs.endpoint"),
		},
		OCSPHost: viper.GetString("ocsp-host"),
		HTTP: config.HTTP{
			OcspListen: viper.GetString("http.ocsp-listen"),
			CaListen:   viper.GetString("http.ca-listen"),
			Listen:     viper.GetString("http.listen"),
		},
		Mysql: config.Mysql{
			Dsn: viper.GetString("mysql.dsn"),
		},
		Vault: config.Vault{
			Enabled:  viper.GetBool("vault.enabled"),
			Addr:     viper.GetString("vault.addr"),
			Token:    viper.GetString("vault.token"),
			Prefix:   viper.GetString("vault.prefix"),
			Discover: viper.GetString("vault.discover"),
		},
		Influxdb: influxdb.CustomConfig{
			Enabled:             viper.GetBool("influxdb.enabled"),
			Address:             viper.GetString("influxdb.address"),
			Port:                viper.GetInt("influxdb.port"),
			UDPAddress:          viper.GetString("influxdb.udp_address"),
			Database:            viper.GetString("influxdb.database"),
			Precision:           viper.GetString("influxdb.precision"),
			UserName:            viper.GetString("influxdb.username"),
			Password:            viper.GetString("influxdb.password"),
			ReadUserName:        viper.GetString("influxdb.read-username"),
			ReadPassword:        viper.GetString("influxdb.read-password"),
			MaxIdleConns:        viper.GetInt("influxdb.max-idle-conns"),
			MaxIdleConnsPerHost: viper.GetInt("influxdb.max-idle-conns-per-host"),
			IdleConnTimeout:     viper.GetInt("influxdb.idle-conn-timeout"),
			FlushSize:           viper.GetInt("influxdb.flush-size"),
			FlushTime:           viper.GetInt("influxdb.flush-time"),
		},
		Mesh: config.Mesh{
			MSPPortalAPI: viper.GetString("mesh.msp-portal-api"),
		},
		SwaggerEnabled: viper.GetBool("swagger-enabled"),
		Debug:          viper.GetBool("debug"),
		Version:        viper.GetString("version"),
		Hostname:       hostname,
		Metrics: config.Metrics{
			CpuLimit: viper.GetFloat64("metrics.cpu-limit"),
			MemLimit: viper.GetFloat64("metrics.mem-limit"),
		},
		Ocsp: config.Ocsp{
			CacheTime: viper.GetInt("ocsp.cache-time"),
		},
	}}

	// ref: https://github.com/golang-migrate/migrate/issues/313
	if !strings.Contains(conf.Mysql.Dsn, "multiStatements") {
		conf.Mysql.Dsn += "&multiStatements=true"
	}

	cfg, err := cfssl_config.LoadFile(conf.Singleca.ConfigPath)
	if err != nil {
		return conf, fmt.Errorf("cfssl 配置文件 %s 错误: %s", conf.Singleca.ConfigPath, err)
	}

	cfg.Signing.Default.OCSP = conf.OCSPHost

	if *flagRootca {
		conf.Keymanager.SelfSign = true
	}

	conf.Singleca.CfsslConfig = cfg

	return conf, nil
}