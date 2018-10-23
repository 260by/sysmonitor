package handlefunc

import (
	"fmt"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintln(w, "Hello World!!!")
}

func HostMonitor(w http.ResponseWriter, r *http.Request) {

}