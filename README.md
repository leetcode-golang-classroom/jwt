# xk6-jwt

這個 repository 主要用來撰寫透過 xk6 這個標準套件來擴充原本 jwt 功能到 k6 runtime 之中
目標未來如果需有一些功能，需要在 k6 引用都可以透過 xk6 來擴充功能到 k6 處理

## 步驟ㄧ 安裝 xk6

```shell
go install go.k6.io/xk6/cmd/xk6@latest
```

## 步驟二 在專案安裝 k6

```shell
go get go.k6.io/k6
```
## 實作

```golang
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
```

## 編譯
```shell
xk6 build --with github.com/leetcode-golang-classroom/jwt=. --output ../k6 
```

## 使用方式
```js
import { JwtGenerator } from 'k6/x/jwt';
const jwtGenerator = new JwtGenerator();

const accessToken = jwtGenerator.signToken(jwtSecret, memberID, uuid, jti)
```
## 參考

https://ganhua.wang/grafana-xk6