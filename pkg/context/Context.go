package context

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/qdnqn/smr/pkg/logger"
	"github.com/qdnqn/smr/pkg/static"
	"github.com/qdnqn/smr/pkg/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func NewContext(projectDir string) *Context {
	return &Context{
		ApiURL:        "",
		Name:          "",
		DirectoryPath: fmt.Sprintf("%s/%s", projectDir, static.CLIENT_CONTEXT_DIR),
		CertBundle:    "",
		ActiveContext: "",
		PrivateKey:    &bytes.Buffer{},
		Cert:          &bytes.Buffer{},
		Ca:            &bytes.Buffer{},
	}
}

func LoadContext(projectDir string) *Context {
	context := &Context{
		ApiURL:        "",
		Name:          "",
		DirectoryPath: fmt.Sprintf("%s/%s", projectDir, static.CLIENT_CONTEXT_DIR),
		CertBundle:    "",
		ActiveContext: "",
		PrivateKey:    &bytes.Buffer{},
		Cert:          &bytes.Buffer{},
		Ca:            &bytes.Buffer{},
	}

	context.GetActiveContext(projectDir)

	if context.GetActiveContext(projectDir) {
		if context.ReadFromFile() {
			var err error
			context.Client, err = context.GenerateHttpClient([]byte(context.CertBundle))

			if err != nil {
				logger.Log.Info("failed to generate http client", zap.String("error", err.Error()))
				return nil
			}

			return context
		}
	}

	return nil
}

func (context *Context) GenerateHttpClient(CertBundle []byte) (*http.Client, error) {
	for block, rest := pem.Decode(CertBundle); block != nil; block, rest = pem.Decode(rest) {
		switch block.Type {
		case "CERTIFICATE":
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}

			if cert.IsCA {
				pem.Encode(context.Ca, &pem.Block{
					Type:  "CERTIFICATE",
					Bytes: cert.Raw,
				})
			} else {
				pem.Encode(context.Cert, &pem.Block{
					Type:  "CERTIFICATE",
					Bytes: cert.Raw,
				})
			}

		case "PRIVATE KEY":
			pem.Encode(context.PrivateKey, &pem.Block{
				Type:  "PRIVATE KEY",
				Bytes: block.Bytes,
			})

		default:
			return nil, errors.New("unknown pem type in the cert bundle")
		}
	}

	context.CertBundle = string(CertBundle)

	cert, err := tls.X509KeyPair(context.Cert.Bytes(), context.PrivateKey.Bytes())
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(context.Ca.Bytes())

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}, nil
}

func (context *Context) GetActiveContext(projectDir string) bool {
	activeContextPath := fmt.Sprintf("%s/%s", context.DirectoryPath, ".active")

	activeContext, err := os.ReadFile(activeContextPath)
	if err != nil {
		logger.Log.Info("active context could not be determined", zap.String("error", err.Error()))
		return false
	}

	context.ActiveContext = fmt.Sprintf("%s/%s/%s", projectDir, static.CLIENT_CONTEXT_DIR, string(activeContext))
	return true
}

func (context *Context) SetActiveContext(contextName string) bool {
	activeContextPath := fmt.Sprintf("%s/%s", context.DirectoryPath, ".active")

	err := os.WriteFile(activeContextPath, []byte(contextName), 0755)
	if err != nil {
		logger.Log.Fatal("active context file not saved", zap.String("error", err.Error()))
	}

	return true
}

func (context *Context) ReadFromFile() bool {
	activeContext, err := os.ReadFile(context.ActiveContext)
	if err != nil {
		logger.Log.Info("active context file not found", zap.String("error", err.Error()))
		return false
	}

	if err = json.Unmarshal(activeContext, &context); err != nil {
		logger.Log.Info("active context file not valid json", zap.String("error", err.Error()))
		return false
	}

	return true
}

func (context *Context) SaveToFile(projectDir string) bool {
	context.Name = viper.GetString("context")
	jsonData, err := json.Marshal(context)

	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	contextPath := fmt.Sprintf("%s/%s/%s", projectDir, "contexts", viper.GetString("context"))

	if _, err := os.Stat(contextPath); err == nil {
		if utils.Confirm("Context with the same name already exists. Do you want to overwrite it?") {
			err = os.WriteFile(contextPath, jsonData, 0600)
			if err != nil {
				logger.Log.Fatal("active context file not saved", zap.String("error", err.Error()))
			}

			activeContextPath := fmt.Sprintf("%s/%s/%s", projectDir, "contexts", ".active")

			err = os.WriteFile(activeContextPath, []byte(viper.GetString("context")), 0755)
			if err != nil {
				logger.Log.Fatal("active context file not saved", zap.String("error", err.Error()))
			}

			return true
		} else {
			logger.Log.Info("action aborted")
			return false
		}
	}

	return false
}
