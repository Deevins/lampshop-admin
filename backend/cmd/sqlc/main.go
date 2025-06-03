package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const defaultConfigName = "sqlc.sqlc.yaml"

func generateConfig(engine *viper.Viper, file string) (string, error) {
	var (
		dir, _      = filepath.Split(strings.Join(filepath.SplitList(filepath.Clean(file)), string(os.PathSeparator)))
		parts       = strings.Split(dir, string(os.PathSeparator))
		packageName = strings.Replace(parts[len(parts)-2], "-", "_", -1)
	)
	engine.Set("gen.go.package", packageName)
	engine.Set("queries", file)

	engine.Set("gen.go.out", dir)
	engineSettings := engine.AllSettings()
	delete(engineSettings, "source")

	resultConfig := viper.New()

	resultConfig.Set("version", viper.GetString("version"))
	resultConfig.Set("sql", []interface{}{engineSettings})

	allSettings := resultConfig.AllSettings()

	bs, err := yaml.Marshal(allSettings)
	if err != nil {
		return "", errors.Wrap(err, "marshal config to sqlc.yaml")
	}
	content := string(bs)
	_ = os.Remove(defaultConfigName)
	temp, err := os.Create(defaultConfigName)
	if err != nil {
		return "", errors.Wrap(err, "create sqlc.sqlc.yaml file")
	}
	if _, err = temp.WriteString(content); err != nil {
		_ = os.Remove(temp.Name())
		return "", errors.Wrap(err, "write content")
	}
	return temp.Name(), nil
}

func callSqlc(config string) error {
	cmd := exec.Command("sqlc", "generate", "--file", config)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("call sqlc: %s", string(output)))
	}
	return nil
}

func main() {
	viper.SetConfigName("sqlc")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	srcConfigValue := viper.GetStringSlice("sql.0.source")
	if len(srcConfigValue) == 0 {
		panic("has no sql.0.source in config")
	}
	files := make([]string, 0)
	for _, pattern := range srcConfigValue {
		f, err := filepath.Glob(pattern)
		if err != nil {
			panic(fmt.Errorf("get file glob: %w", err))
		}
		files = append(files, f...)
	}
	schemaConfigValue := viper.GetString("sql.0.schema")

	engine := viper.Sub("sql.0")
	engine.Set("schema", schemaConfigValue)

	for _, file := range files {
		configFile, gErr := generateConfig(engine, file)
		if gErr != nil {
			panic(fmt.Errorf("can't generate result config: %w", gErr))
		}
		if cErr := callSqlc(configFile); cErr != nil {
			panic(fmt.Errorf("call sqlc: %w", cErr))
		}
		fmt.Printf("%s file complete\n", file)
	}
	_ = os.Remove(defaultConfigName)
	fmt.Println("done")
}
