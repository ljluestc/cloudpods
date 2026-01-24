// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nat

import (
	"context"

	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"

	api "yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/cloudcommon/db/taskman"
	"yunion.io/x/onecloud/pkg/compute/models"
	"yunion.io/x/onecloud/pkg/util/logclient"
)

type NatGatewayChangeSpecTask struct {
	taskman.STask
}

func init() {
	taskman.RegisterTask(NatGatewayChangeSpecTask{})
}

func (self *NatGatewayChangeSpecTask) taskFailed(ctx context.Context, nat *models.SNatGateway, err error) {
	nat.SetStatus(ctx, self.UserCred, api.NAT_STAUTS_AVAILABLE, err.Error())
	logclient.AddActionLogWithStartable(self, nat, logclient.ACT_CHANGE_CONFIG, err, self.UserCred, false)
	self.SetStageFailed(ctx, jsonutils.NewString(err.Error()))
}

func (self *NatGatewayChangeSpecTask) OnInit(ctx context.Context, obj db.IStandaloneModel, body jsonutils.JSONObject) {
	nat := obj.(*models.SNatGateway)

	natSpec, _ := self.GetParams().GetString("nat_spec")
	if len(natSpec) == 0 {
		self.taskFailed(ctx, nat, errors.Wrap(errors.ErrInvalidStatus, "empty nat_spec"))
		return
	}

	self.SetStage("OnNatGatewayChangeSpecComplete", nil)

	taskman.LocalTaskRun(self, func() (jsonutils.JSONObject, error) {
		iNat, err := nat.GetINatGateway(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "nat.GetINatGateway")
		}

		err = iNat.SetNatSpec(natSpec)
		if err != nil {
			return nil, errors.Wrap(err, "iNat.SetNatSpec")
		}

		return nil, nil
	})
}

func (self *NatGatewayChangeSpecTask) OnNatGatewayChangeSpecComplete(ctx context.Context, nat *models.SNatGateway, body jsonutils.JSONObject) {
	natSpec, _ := self.GetParams().GetString("nat_spec")
	db.Update(nat, func() error {
		nat.NatSpec = natSpec
		return nil
	})
	nat.SetStatus(ctx, self.UserCred, api.NAT_STAUTS_AVAILABLE, "")
	logclient.AddActionLogWithStartable(self, nat, logclient.ACT_CHANGE_CONFIG, nil, self.UserCred, true)
	self.SetStageComplete(ctx, nil)
}

func (self *NatGatewayChangeSpecTask) OnNatGatewayChangeSpecCompleteFailed(ctx context.Context, nat *models.SNatGateway, body jsonutils.JSONObject) {
	self.taskFailed(ctx, nat, errors.Errorf("%s", body.String()))
}
