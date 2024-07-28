package cli

import (
	"fmt"
	"loadtester/internal/lib/models"
	"loadtester/internal/lib/service"
)

func main() {
	tp := models.TestPlan{}
	_ = service.NewTestRunner(tp)
	fmt.Println("test for now")
}
