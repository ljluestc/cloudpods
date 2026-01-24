package models

import (
"context"
"testing"

api "yunion.io/x/onecloud/pkg/apis/compute"
"yunion.io/x/onecloud/pkg/cloudcommon/db"
)

func TestSNatGateway_PerformChangeSpec(t *testing.T) {
cases := []struct {
Name      string
Nat       *SNatGateway
Input     api.NatGatewayChangeSpecInput
WantError bool
}{
{
Name: "Invalid Status",
Nat: &SNatGateway{
SStatusInfrasResourceBase: db.SStatusInfrasResourceBase{
SStatusResourceBase: db.SStatusResourceBase{
Status: api.NAT_STATUS_DEPLOYING,
},
},
},
Input: api.NatGatewayChangeSpecInput{
NatSpec: "small",
},
WantError: true,
},
{
Name: "Empty Spec",
Nat: &SNatGateway{
SStatusInfrasResourceBase: db.SStatusInfrasResourceBase{
SStatusResourceBase: db.SStatusResourceBase{
Status: api.NAT_STAUTS_AVAILABLE,
},
},
},
Input: api.NatGatewayChangeSpecInput{
NatSpec: "",
},
WantError: true,
},
{
Name: "Same Spec",
Nat: &SNatGateway{
SStatusInfrasResourceBase: db.SStatusInfrasResourceBase{
SStatusResourceBase: db.SStatusResourceBase{
Status: api.NAT_STAUTS_AVAILABLE,
},
},
NatSpec: "small",
},
Input: api.NatGatewayChangeSpecInput{
NatSpec: "small",
},
WantError: false,
},
}

for _, c := range cases {
t.Run(c.Name, func(t *testing.T) {
_, err := c.Nat.PerformChangeSpec(context.Background(), nil, nil, c.Input)
if c.WantError && err == nil {
t.Errorf("Expected error but got nil")
}
if !c.WantError && err != nil {
t.Errorf("Expected nil error but got %v", err)
}
})
}
}
