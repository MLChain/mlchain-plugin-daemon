package real

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mlchain/mlchain-plugin-daemon/internal/core/mlchain_invocation"
	"github.com/mlchain/mlchain-plugin-daemon/internal/types/entities/plugin_entities"
	"github.com/mlchain/mlchain-plugin-daemon/internal/utils/network"
)

func TestEncryptRequired(t *testing.T) {
	data := map[string]any{
		"key": "value",
	}

	payload := &mlchain_invocation.InvokeEncryptRequest{
		BaseInvokeMlchainRequest: mlchain_invocation.BaseInvokeMlchainRequest{
			TenantId: "123",
			UserId:   "456",
			Type:     mlchain_invocation.INVOKE_TYPE_ENCRYPT,
		},
		InvokeEncryptSchema: mlchain_invocation.InvokeEncryptSchema{
			Opt:       mlchain_invocation.ENCRYPT_OPT_ENCRYPT,
			Namespace: mlchain_invocation.ENCRYPT_NAMESPACE_ENDPOINT,
			Identity:  "test123",
			Data:      data,
			Config: []plugin_entities.ProviderConfig{
				{
					Name: "key",
					Type: plugin_entities.CONFIG_TYPE_SECRET_INPUT,
				},
			},
		},
	}

	if !payload.EncryptRequired(data) {
		t.Errorf("EncryptRequired should return true")
	}

	payload.Config = []plugin_entities.ProviderConfig{
		{
			Name: "key",
			Type: plugin_entities.CONFIG_TYPE_TEXT_INPUT,
		},
	}

	if payload.EncryptRequired(data) {
		t.Errorf("EncryptRequired should return false")
	}
}

func TestInvokeEncrypt(t *testing.T) {
	server := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	port, err := network.GetRandomPort()
	if err != nil {
		t.Errorf("GetRandomPort failed: %v", err)
	}

	httpInvoked := false

	server.POST("/inner/api/invoke/encrypt", func(ctx *gin.Context) {
		data := make(map[string]any)
		if err := ctx.BindJSON(&data); err != nil {
			t.Errorf("BindJSON failed: %v", err)
		}

		if data["data"].(map[string]any)["key"] != "value" {
			t.Errorf("data[key] should be `value`, but got %v", data["data"].(map[string]any)["key"])
		}

		httpInvoked = true

		ctx.JSON(http.StatusOK, gin.H{
			"data": map[string]any{
				"data": map[string]any{
					"key": "encrypted",
				},
			},
		})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: server,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	defer srv.Close()

	time.Sleep(1 * time.Second)

	i, err := NewMlchainInvocationDaemon(fmt.Sprintf("http://localhost:%d", port), "test")
	if err != nil {
		t.Errorf("InitMlchainInvocationDaemon failed: %v", err)
		return
	}

	payload := &mlchain_invocation.InvokeEncryptRequest{
		BaseInvokeMlchainRequest: mlchain_invocation.BaseInvokeMlchainRequest{
			TenantId: "123",
			UserId:   "456",
			Type:     mlchain_invocation.INVOKE_TYPE_ENCRYPT,
		},
		InvokeEncryptSchema: mlchain_invocation.InvokeEncryptSchema{
			Opt:       mlchain_invocation.ENCRYPT_OPT_ENCRYPT,
			Namespace: mlchain_invocation.ENCRYPT_NAMESPACE_ENDPOINT,
			Identity:  "test123",
			Data:      map[string]any{"key": "value"},
			Config: []plugin_entities.ProviderConfig{
				{
					Name: "key",
					Type: plugin_entities.CONFIG_TYPE_SECRET_INPUT,
				},
			},
		},
	}

	if encrypted, err := i.InvokeEncrypt(payload); err != nil {
		t.Errorf("InvokeEncrypt failed: %v", err)
	} else {
		if encrypted["key"] != "encrypted" {
			t.Errorf("encrypted[key] should be `encrypted`, but got %v", encrypted["key"])
		}
	}

	if !httpInvoked {
		t.Errorf("http_invoked should be true")
	}
}
