package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"sync"
)

var (
	// vaultFile is the vault file location
	configDir, _ = os.UserConfigDir()
	vaultFile    = configDir + "/portunus.json"

	// vault errors
	errorVaultExists      = errors.New("vault file already exists at " + vaultFile)
	errorVaultNotExists   = errors.New("no vault file found at " + vaultFile)
	errorVaultInvalid     = errors.New("invalid vault file at " + vaultFile)
	errorVaultNoSuchValue = errors.New("no such value in vault")

	// argument parsing errors
	errorBadArgs    = errors.New("possible subcommands 'vlt', 'get', 'set', 'new', 'lst', 'gen'")
	errorBadArgsSet = errors.New("'set' takes one argument, 'name'")
	errorBadArgsNew = errors.New("'new' takes one argument, 'name'")
	errorBadArgsGet = errors.New("'get' takes one argument, 'name'")
	errorBadArgsGen = errors.New("'gen' takes one argument, 'name'")
)

type vault struct {
	vlt  map[string]string
	lock sync.Mutex
}

func newVault() (*vault, error) {
	vlt := &vault{vlt: make(map[string]string)}
	fd, err := os.OpenFile(vaultFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, errorVaultExists
		}
		return nil, err
	}
	defer fd.Close()
	_, err = fd.WriteString("{}")
	return vlt, err
}

func openVault() (*vault, error) {
	vlt := &vault{vlt: make(map[string]string)}
	data, err := ioutil.ReadFile(vaultFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errorVaultNotExists
		}
		return nil, err
	}
	err = json.Unmarshal(data, &vlt.vlt)
	if err != nil {
		return nil, errorVaultInvalid
	}
	return vlt, nil
}

func (vlt *vault) set(name string) {
	vlt.lock.Lock()
	defer vlt.lock.Unlock()
	vlt.vlt[name] = readPassword()
}

func (vlt *vault) new(name string) {
	vlt.lock.Lock()
	defer vlt.lock.Unlock()
	vlt.vlt[name] = generatePassword()
}

func (vlt *vault) get(name string) (string, error) {
	vlt.lock.Lock()
	defer vlt.lock.Unlock()
	pswd, ok := vlt.vlt[name]
	if !ok {
		return "", errorVaultNoSuchValue
	}
	return pswd, nil
}

func (vlt *vault) rem(name string) error {
	vlt.lock.Lock()
	defer vlt.lock.Unlock()
	if _, ok := vlt.vlt[name]; !ok {
		return errorVaultNoSuchValue
	}
	delete(vlt.vlt, name)
	return nil
}

func (vlt *vault) lst() []string {
	names := make([]string, len(vlt.vlt))
	var i int
	for name := range vlt.vlt {
		names[i] = name
		i++
	}
	sort.Strings(names)
	return names
}

func (vlt *vault) saveVault() error {
	data, _ := json.Marshal(vlt.vlt)
	err := ioutil.WriteFile(vaultFile, data, 0644)
	return err
}

func generatePassword() string {
	pswd := make([]byte, 12)
	rand.Read(pswd)
	return base64.RawURLEncoding.EncodeToString(pswd)
}

func readPassword() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func main() {
	if len(os.Args) < 2 {
		exit(errorBadArgs)
	}

	vlt, err := openVault()
	if err != nil && os.Args[1] != "new" {
		exit(err)
	}

	switch os.Args[1] {
	case "vlt":
		_, err := newVault()
		exit(err)
	case "set":
		if len(os.Args) != 3 {
			exit(errorBadArgsSet)
		}
		name := os.Args[2]
		vlt.set(name)
		exit(vlt.saveVault())
	case "new":
		if len(os.Args) != 3 {
			exit(errorBadArgsNew)
		}
		name := os.Args[2]
		vlt.new(name)
	case "get":
		if len(os.Args) != 3 {
			exit(errorBadArgsGet)
		}
		name := os.Args[2]
		pswd, err := vlt.get(name)
		exit(err)
		fmt.Println(pswd)
	case "lst":
		for _, name := range vlt.lst() {
			fmt.Println(name)
		}
	case "gen":
		fmt.Println(generatePassword())
	default:
		exit(errorBadArgs)
	}
}

func exit(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("portunus: %w", err))
		os.Exit(1)
	}
}
