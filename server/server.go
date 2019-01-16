package server

import (
	"net/http"
	"fmt"

	"github.com/kraftykai/breview/configs"
)

func Init(c *configs.Definitions) error {
	//TODO: Add hostname to addy (Issue #1)
	addy := fmt.Sprintf(":%d", c.Port)
	return http.ListenAndServe(addy, nil);
}
