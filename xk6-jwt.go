package jwt

import (
	"fmt"

	"github.com/grafana/sobek"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

func init() {
	// 註冊模組到 k6 系統 (強制前綴 k6/x/)
	modules.Register("k6/x/jwt", new(RootModule))
}

// RootModule - 全局模組
type RootModule struct{}

// NewModuleInstance implements modules.Module.
func (r *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ModuleInstance{vu: vu}
}

// ModuleInstance - 模組實例
type ModuleInstance struct {
	vu modules.VU
}

// Exports implements modules.Instance.
func (m *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"JwtGenerator": m.newJwtGenerator, // JS 中 new OtpGenerator()
		},
	}
}

func (m *ModuleInstance) newJwtGenerator(c sobek.ConstructorCall) *sobek.Object {
	rt := m.vu.Runtime()
	if len(c.Arguments) != 0 {
		// 統一使用 common.Throw 拋出異常
		common.Throw(rt, fmt.Errorf("JwtGenerator requires no argument"))
	}
	generator := &JwtGenerator{
		vu: m.vu,
	}

	// Create JS object and bind methods
	obj := rt.NewObject()
	// Bind generate method
	if err := obj.Set("signToken", generator.signToken); err != nil {
		common.Throw(rt, err)
	}
	return obj
}

// 必須實現的 interface 驗證
var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &ModuleInstance{}
)

type JwtGenerator struct {
	vu modules.VU // 必備的 VU 引用
}

func (g *JwtGenerator) signToken(c sobek.FunctionCall) sobek.Value {
	rt := g.vu.Runtime()

	// Validate constructor arguments
	if len(c.Arguments) != 4 {
		// 統一使用 common.Throw 拋出異常
		common.Throw(rt, fmt.Errorf("OtpGenerator requires 4 argument (jwtSecret, memberID, uuid, jti)"))
	}

	jwtSecret := c.Argument(0).String()
	if jwtSecret == "" {
		common.Throw(rt, fmt.Errorf("jwtSecret cannot be empty"))
	}

	memberID := c.Argument(1).String()
	if memberID == "" {
		common.Throw(rt, fmt.Errorf("memberID cannot be empty"))
	}

	uuid := c.Argument(2).String()
	if uuid == "" {
		common.Throw(rt, fmt.Errorf("uuid cannot be empty"))
	}

	jti := c.Argument(3).String()
	if jti == "" {
		common.Throw(rt, fmt.Errorf("jti cannot be empty"))
	}

	token, err := SignAccessToken(jwtSecret, memberID, uuid, jti)
	if err != nil {
		common.Throw(rt, err)
	}

	return rt.ToValue(token)
}
