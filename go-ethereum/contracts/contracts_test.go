package contracts

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var buffer []byte = `{
	"dpos": {
		"validators": ["0x1aa228dde26f02e1cc5551cb4f1d74d0e998d24a"]
	},
	"contracts": [{
			"contractType": "1",
			"deployAddress": "0x1aa228dde26f02e1cc5551cb4f1d74d0e998d24a",
			"contractAddress": "0x1aa228dde26f02e1cc5551cb4f1d74d0e998d24a"
		},
		{
			"contractType": "2",
			"deployAddress": "0x1aa228dde26f02e1cc5551cb4f1d74d0e998d24a",
			"contractAddress": "0x1aa228dde26f02e1cc5551cb4f1d74d0e998d24a"
		}
	]
}`

func TestLoadConfig() error {

	config := BokerConfig{}
	err = json.Unmarshal(buffer, &config)
	if err != nil {
		//log.Error("Unmarshal boker.json File Error", "Error", err)
		return err
	}

}

func main() {

	err := TestLoadConfig()

	for {

		time.Sleep(1 * time.Second)

	}

}
