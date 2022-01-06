package validate

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

// CheckPort 检查端口是否可用
func CheckPort(port int) error {
	checkStatement := fmt.Sprintf(`netstat -anp | grep -q %d ; echo $?`, port)
	output, err := exec.Command("sh", "-c", checkStatement).CombinedOutput()
	if err != nil {
		return err
	}

	result, err := strconv.Atoi(strings.TrimSuffix(string(output), "\n"))
	if err != nil {
		return err
	}
	if result == 0 {
		return fmt.Errorf("port %d already in use", port)
	}

	return nil
}
