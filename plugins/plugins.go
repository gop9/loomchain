package plugins

import (
	"fmt"
	"path/filepath"
	"plugin"

	"github.com/loomnetwork/loom"
)

// SimpleContract is our interface for inmemory contracts
type SimpleContract interface {
	// Plugin name to display in logs etc.
	Name() string
	Init() error
	Routes() []int
	HandleRoutes(loom.State, string, []byte) ([]byte, error)
}

func AttachBuiltinPlugins(plugins []SimpleContract, router *loom.TxRouter) error {
	for _, p := range plugins {
		if err := attachBuiltinPlugin(p, router); err != nil {
			fmt.Printf("error loading built-in plugin -%s-%v\n", p.Name(), err)
		}
	}
	return nil
}

func AttachLocalPlugins(path string, router *loom.TxRouter) error {
	files, err := filepath.Glob(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		fmt.Println(f)
		err := attachLocalPlugin(f, router)
		if err != nil {
			fmt.Printf("error loading local plugin -%s-%v\n", f, err)
		}
	}
	return nil
}

func attachLocalPlugin(filename string, router *loom.TxRouter) error {
	//TODO iterate over the folder and load all the plugins

	// load module
	// 1. open the so file to load the symbols
	plug, err := plugin.Open(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 2. look up a symbol (an exported function or variable)
	contractsPlug, err := plug.Lookup("Contract")
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 3. Assert that loaded symbol is of a desired type
	// in this case interface type SimpleContract (defined above)
	var contract SimpleContract
	contract, ok := contractsPlug.(SimpleContract)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		return err
	}
	// 4. init the module
	err = contract.Init()
	if err != nil {
		return err
	}
	// 5. use the module
	res := contract.Routes()
	fmt.Printf("weee -%v\n", res)

	return nil
}

func attachBuiltinPlugin(plugin SimpleContract, router *loom.TxRouter) error {
	if err := plugin.Init(); err != nil {
		return err
	}
	res := plugin.Routes()
	fmt.Printf("weee -%v\n", res)
	return nil
}
