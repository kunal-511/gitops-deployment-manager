package deploy

import (
	"fmt"
	"os/exec"
)

func ApplyManifests(path string) error {
	cmd := exec.Command(("kubectl"), "apply", "-f", path)
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	return err
}
