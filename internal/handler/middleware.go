package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-gin-gorm-mysql/internal/core/config"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Request request from client.
func Request() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		r := c.Request
		contentType := r.Header.Get("Content-Type")

		dump, err := httputil.DumpRequest(c.Request, true)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, struct {
				Message string `json:"message"`
			}{
				Message: "ขออภัย ระบบไม่สามารถดำเนินการได้ในขณะนี้ กรุณาลองใหม่อีกครั้ง",
			})
			return
		}

		if strings.HasPrefix(contentType, "application/json") && r.Method != http.MethodGet {
			buffer, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logrus.Errorf("[Request] read body request error: %s", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, config.RR.Internal.BadRequest)
				return
			}

			rc := ioutil.NopCloser(bytes.NewBuffer(buffer))
			r.Body = rc
			err = JSONDuplicate(json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(buffer))), nil)
			if err != nil {
				logrus.Errorf("[Request] check duplicate json body request error: %s", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, config.RR.Internal.BadRequest)
				return
			}
		}

		acceptLanguage(c)

		// // intercept writing response and store it to write log.
		buff := bytes.Buffer{}
		bw := bufferWriter{c.Writer, &buff}
		c.Writer = bw
		c.Next()

		hostname, _ := os.Hostname()
		response := c.Writer.(gin.ResponseWriter)
		// write response log
		logs := logrus.Fields{
			"host":            hostname,
			"method":          r.Method,
			"path":            r.URL.Path,
			"Accept-Language": r.Header.Get("Accept-Language"),
			"clientIP":        GetIPAddress(r),
			"User-Agent":      r.Header.Get("User-Agent"),
			"statusCode":      fmt.Sprintf("%d %s", response.Status(), http.StatusText(response.Status())),
			"processTime":     time.Since(start),
		}

		logrus.WithFields(logs).Info(fmt.Sprintf("parameter: %+v", string(dump)))
	}
}

type bufferWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

type responseWriterModel struct {
	http.ResponseWriter
	body   *bytes.Buffer
	status int
}

func newResponseWriter(rw http.ResponseWriter, buffer *bytes.Buffer) *responseWriterModel {
	nrw := &responseWriterModel{
		ResponseWriter: rw,
		body:           buffer,
	}
	return nrw
}

// Write write data
func (rw responseWriterModel) Write(data []byte) (int, error) {
	rw.body.Write(data)
	return rw.ResponseWriter.Write(data)
}

// WriteHeader write header
func (rw *responseWriterModel) WriteHeader(s int) {
	rw.status = s
	rw.ResponseWriter.WriteHeader(s)
}

// Status get status from responseWriterModel
func (rw *responseWriterModel) Status() int {
	return rw.status
}

// IPFromRequest get ip address
func IPFromRequest(req *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	return userIP, nil
}

// GetIPAddress get ip address from request
func GetIPAddress(r *http.Request) string {
	clientIP := r.Header.Get("X-Forwarded-For")
	if userIP, ok := IPFromRequest(r); ok == nil {
		lastIP := userIP.String()
		if clientIP == "" {
			clientIP = lastIP
		} else {
			clientIP += fmt.Sprintf(", %s", lastIP)
		}
	}

	return clientIP
}

// JSONDuplicate check json is duplicate
func JSONDuplicate(d *json.Decoder, path []string) error {
	var duplicate string

	// Get next token from JSON
	t, err := d.Token()
	if err != nil {
		return err
	}

	delim, ok := t.(json.Delim)

	// There's nothing to do for simple values (strings, numbers, bool, nil)
	if !ok {
		return nil
	}

	switch delim {
	case '{':
		keys := make(map[string]bool)
		for d.More() {
			// Get field key
			t, err := d.Token()
			if err != nil {
				return err
			}
			key := t.(string)

			// Check for duplicates
			if keys[key] {
				duplicate = duplicate + fmt.Sprint(strings.Join(append(path, "dulicate field: "+key), ","))
			}
			keys[key] = true

			// Check value
			if err := JSONDuplicate(d, append(path, key)); err != nil {
				return err
			}
		}
		// Consume trailing }
		if _, err := d.Token(); err != nil {
			return err
		}

	case '[':
		i := 0
		for d.More() {
			if err := JSONDuplicate(d, append(path, strconv.Itoa(i))); err != nil {
				return err
			}
			i++
		}
		// Consume trailing ]
		if _, err := d.Token(); err != nil {
			return err
		}

	}
	if duplicate != "" {
		return errors.New(duplicate)
	}
	return nil
}

func acceptLanguage(c *gin.Context) {
	lang := c.Request.Header.Get("Accept-Language")
	if lang == "" || lang != "th" {
		lang = "en"
	}
	config.Set(c, config.LangKey, lang)
	c.Next()
}
