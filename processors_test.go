package raven

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"testing"
)

func customProcessor(packet *Packet) *Packet {
	packet.Message = strings.Replace(packet.Message, "Password", "********", -1)
	return packet
}

func TestProcessor(t *testing.T) {
	// ... i.e. raisedErr is incoming error
	raisedErr := errors.New("Test Password Error")
	// sentry DSN generated by Sentry server
	var sentryDSN string
	// r is a request performed when error occured

	r, _ := http.NewRequest("GET", "http://example.com/", nil)

	r.Header.Set("api_key", "SGFycnlQb3R0ZXI6RHVtYmxlZG9yZXNBcm15")
	r.Header.Set("apikey", "SGFycnlQb3R0ZXI6RHVtYmxlZG9yZXNBcm15")
	r.Header.Set("Authorization", "Basic SGFycnlQb3R0ZXI6RHVtYmxlZG9yZXNBcm15")
	r.Header.Set("X-Passwd", "62442")
	r.Header.Set("X-Password", "62442")
	r.Header.Set("X-Secret-Key", "SGFycnlQb3R0ZXI6RHVtYmxlZG9yZXNBcm15")
	r.Header.Set("X-PASS", "Harry Potter")

	r.URL.RawQuery = "api_key=SGFycnlQb3R0ZXI6RHVtYmxlZG9yZXNBcm15&apikey=SGFycnlQb3R0ZXI6RHVtYmxlZG9yZXNBcm15&auth=SGFycnlQb3R0ZXI6RHVtYmxlZG9yZXNBcm15&passwd=62442&userpassword=62442&secret-keeper=Sirius&pass=62442"

	config := &ClientConfig{&[]Processor{customProcessor}, nil}
	client, err := NewClient(sentryDSN, config)
	if err != nil {
		log.Fatal(err)
	}
	trace := NewStacktrace(0, 2, nil)
	packet := NewPacket(raisedErr.Error(), NewException(raisedErr, trace), NewHttp(r))

	scrubbedPacket := client.Scrub(packet)

	scrubbedJSON := string(scrubbedPacket.JSON())
	expectedJSON := `{"message":"Test ******** Error","event_id":"","project":"","timestamp":"0001-01-01T00:00:00","level":"","logger":"","extra":{"runtime.GOMAXPROCS":1,"runtime.NumCPU":8,"runtime.NumGoroutine":6,"runtime.Version":"go1.3.1"},"sentry.interfaces.Exception":{"value":"Test Password Error","type":"*errors.errorString","stacktrace":{"frames":[{"filename":"runtime/proc.c","function":"goexit","module":"runtime","lineno":1445,"abs_path":"/usr/local/Cellar/go/1.3.1/libexec/src/pkg/runtime/proc.c","context_line":"runtime·goexit(void)","pre_context":["#pragma textflag NOSPLIT","void"],"post_context":["{","\u0009if(g-\u003estatus != Grunning)"],"in_app":false},{"filename":"testing/testing.go","function":"tRunner","module":"testing","lineno":422,"abs_path":"/usr/local/Cellar/go/1.3.1/libexec/src/pkg/testing/testing.go","context_line":"\u0009test.F(t)","pre_context":["","\u0009t.start = time.Now()"],"post_context":["\u0009t.finished = true","}"],"in_app":false},{"filename":"Github.com/mattgerstman/raven-go/processors_test.go","function":"TestProcessor","module":"Github.com/mattgerstman/raven-go","lineno":40,"abs_path":"/Users/Matthew/Documents/Go/src/Github.com/mattgerstman/raven-go/processors_test.go","context_line":"\u0009trace := NewStacktrace(0, 2, nil)","pre_context":["\u0009\u0009log.Fatal(err)","\u0009}"],"post_context":["\u0009packet := NewPacket(raisedErr.Error(), NewException(raisedErr, trace), NewHttp(r))",""],"in_app":false}]}},"sentry.interfaces.Http":{"url":"http://example.com/","method":"GET","query_string":"api_key=%2A%2A%2A%2A%2A%2A%2A%2A\u0026apikey=%2A%2A%2A%2A%2A%2A%2A%2A\u0026auth=SGFycnlQb3R0ZXI6RHVtYmxlZG9yZXNBcm15\u0026pass=62442\u0026passwd=%2A%2A%2A%2A%2A%2A%2A%2A\u0026secret-keeper=%2A%2A%2A%2A%2A%2A%2A%2A\u0026userpassword=%2A%2A%2A%2A%2A%2A%2A%2A","headers":{"Api_key":"********","Apikey":"********","Authorization":"********","X-Pass":"Harry Potter","X-Passwd":"********","X-Password":"********","X-Secret-Key":"********"}}}`

	if scrubbedJSON != expectedJSON {
		t.Errorf("incorrect Value: got %s, want %s", scrubbedJSON, expectedJSON)
	}

}