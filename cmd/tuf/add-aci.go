package main

import (
	"encoding/json"
	"os"

	"github.com/appc/spec/aci"
	"github.com/flynn/go-tuf"
	"github.com/flynn/go-tuf/Godeps/_workspace/src/github.com/flynn/go-docopt"
)

func init() {
	register("add-aci", cmdAddACI, `
usage: tuf add-aci [--expires=<days>] [<path>...]

Add target file(s).

Options:
  --expires=<days>   Set the targets manifest to expire <days> days from now.
`)
}

func cmdAddACI(args *docopt.Args, repo *tuf.Repo) error {
	var custom json.RawMessage
	if c := args.String["--custom"]; c != "" {
		custom = json.RawMessage(c)
	}
	paths := args.All["<path>"].([]string)
	f, err := os.Open("staged/targets/" + paths[0])
	if err != nil {
		panic(err)
	}
	img, err := aci.ManifestFromImage(f)
	if err != nil {
		panic(err)
	}
	custom, err = img.MarshalJSON()
	if err != nil {
		panic(err)
	}
	if arg := args.String["--expires"]; arg != "" {
		expires, err := parseExpires(arg)
		if err != nil {
			return err
		}
		return repo.AddTargetsWithExpires(paths, custom, expires)
	}
	return repo.AddTargets(paths, custom)
}
