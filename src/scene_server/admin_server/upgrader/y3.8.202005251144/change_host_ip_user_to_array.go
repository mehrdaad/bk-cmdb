/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package y3_8_202005251144

import (
	"context"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// change host inner ip and outer ip and operator and bak operator value from string split by comma to array
func changeHostIPAndUserToArray(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	count, err := db.Table(common.BKTableNameBaseHost).Find(nil).Count(ctx)
	if err != nil {
		blog.Errorf("count hosts failed, err: %s", err.Error())
		return err
	}
	needChangeFields := []string{common.BKHostInnerIPField, common.BKHostOuterIPField, common.BKOperatorField, common.BKBakOperatorField}
	for i := uint64(0); i < count; i += common.BKMaxPageSize {
		hosts := make([]metadata.HostMapStr, 0)
		if err := db.Table(common.BKTableNameBaseHost).Find(nil).Start(i).Limit(common.BKMaxPageSize).
			Fields(append(needChangeFields, common.BKHostIDField)...).All(ctx, &hosts); err != nil {
			blog.Errorf("find hosts starting from %d failed, err: %s", i, err.Error())
			return err
		}
		for _, host := range hosts {
			filter := map[string]interface{}{
				common.BKHostIDField: host[common.BKHostIDField],
			}
			doc := make(map[string]interface{})
			for _, field := range needChangeFields {
				if host[field] == nil {
					doc[field] = make([]string, 0)
					continue
				}
				if value, ok := host[field].(string); ok {
					if len(value) == 0 {
						doc[field] = make([]string, 0)
					} else {
						doc[field] = strings.Split(value, ",")
					}
				}
			}
			if len(doc) == 0 {
				continue
			}
			if err := db.Table(common.BKTableNameBaseHost).Update(ctx, filter, doc); err != nil {
				blog.ErrorJSON("update host ip to array failed, filter: %s, doc: %s, err: %s", filter, doc, err)
				return err
			}
		}
	}
	return nil
}
