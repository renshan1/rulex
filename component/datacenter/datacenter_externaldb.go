// Copyright (C) 2023 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package datacenter

import (
	"context"
	"strconv"
	"time"

	"github.com/hootrhino/rulex/component/trailer"
	"github.com/hootrhino/rulex/glogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
*
* 外部扩展数据库
*
 */
type ExternalDb struct {
}

func (db *ExternalDb) Init() error {
	return nil
}

func (db *ExternalDb) Name() string {
	return "EXTERNALDB"
}
func (db *ExternalDb) GetSchemaDetail(goodsId string) SchemaDetail {
	return SchemaDetail{}
}
func (db *ExternalDb) Query(goodsId, query string) ([]map[string]any, error) {
	return []map[string]any{}, nil
}

/*
*
* 获取表格定义
*
 */
func SchemaList() []SchemaDetail {
	Schemas := []SchemaDetail{}
	trailer.AllGoods().Range(func(key, value any) bool {
		goodsPs := (value.(*trailer.GoodsProcess))
		Schemas = append(Schemas, SchemaDetail{
			UUID:        goodsPs.Info.UUID,
			Name:        goodsPs.Info.Name,
			LocalPath:   goodsPs.Info.LocalPath,
			NetAddr:     goodsPs.Info.NetAddr,
			CreateTs:    0,
			Size:        0,
			StorePath:   "",
			Description: goodsPs.Info.Description,
		})
		return true
	})
	return Schemas
}

/*
*
* 表结构
*
 */

func GetSchema(goodsId string) (SchemaDefine, error) {
	schemaDefine := SchemaDefine{}
	if goodsPs := trailer.Get(goodsId); goodsPs != nil {
		grpcConnection, err1 := grpc.Dial(goodsPs.Info.NetAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err1 != nil {
			glogger.GLogger.Error(err1)
			return schemaDefine, err1
		}
		defer grpcConnection.Close()
		client := trailer.NewTrailerClient(grpcConnection)
		columns, err2 := client.Schema(context.Background(), &trailer.SchemaRequest{})
		if err2 != nil {
			glogger.GLogger.Error(err2)
			return schemaDefine, err2
		}
		Columns := []Column{}
		for _, column := range columns.Columns {
			Columns = append(Columns, Column{
				Name:        string(column.Name),
				Type:        column.Type.String(),
				Description: string(column.Description),
			})
		}
		schemaDefine = SchemaDefine{
			UUID:    goodsPs.Info.UUID,
			Columns: Columns,
		}
	}
	return schemaDefine, nil

}
func SchemaDefineList() ([]SchemaDefine, error) {
	var err error
	ColumnsMap := []SchemaDefine{}
	Columns := []Column{}
	trailer.AllGoods().Range(func(key, value any) bool {
		goodsPs := (value.(*trailer.GoodsProcess))
		grpcConnection, err1 := grpc.Dial(goodsPs.Info.NetAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err1 != nil {
			glogger.GLogger.Error(err1)
			err = err1
			return false
		}
		defer grpcConnection.Close()
		client := trailer.NewTrailerClient(grpcConnection)
		columns, err2 := client.Schema(context.Background(), &trailer.SchemaRequest{})
		if err2 != nil {
			glogger.GLogger.Error(err2)
			err = err2
			return false
		}
		for _, column := range columns.Columns {
			Columns = append(Columns, Column{
				Name:        string(column.Name),
				Type:        column.Type.String(),
				Description: string(column.Description),
			})
		}
		Define := SchemaDefine{
			UUID:    goodsPs.Info.UUID,
			Columns: Columns,
		}
		ColumnsMap = append(ColumnsMap, Define)
		return true
	})
	return ColumnsMap, err
}

/*
*
* 获取仓库详情, 现阶段写死的, 后期会在proto中实现
*
 */
func GetSchemaDetail(goodsId string) SchemaDetail {
	return SchemaDetail{
		UUID:        "1122334455",
		Name:        "Test RPC",
		LocalPath:   "/root/app1",
		NetAddr:     "127.0.0.1:4567",
		CreateTs:    uint64(time.Now().Unix()),
		Size:        12.34,
		StorePath:   "/root/data/test.db",
		Description: "An simply demo",
	}
}

/*
*
* 查询，第一个参数是查询请求，针对Sqlite就是SQL，针对mongodb就是JS，根据具体情况而定
*
 */

func Query(goodsId, query string) ([]map[string]any, error) {
	var err error
	Rows := []map[string]any{}
	if goodsPs := trailer.Get(goodsId); goodsPs != nil {
		grpcConnection, err1 := grpc.Dial(goodsPs.Info.NetAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err1 != nil {
			glogger.GLogger.Error(err1)
			err = err1
			return Rows, err
		}
		defer grpcConnection.Close()
		client := trailer.NewTrailerClient(grpcConnection)
		columns, err2 := client.Query(context.Background(), &trailer.DataRowsRequest{
			Query: []byte(query),
		})
		if err2 != nil {
			glogger.GLogger.Error(err2)
			err = err2
			return Rows, err
		}
		for _, row := range columns.Row {
			Row := map[string]any{}
			for _, column := range row.Column {
				Row[string(column.GetName())] = CovertGoTypeToJsType(column)
			}
			Rows = append(Rows, Row)
		}
	}
	return Rows, err
}

/*
*
* 数据转换
*
 */
func CovertGoTypeToJsType(V *trailer.ColumnValue) any {
	if V.Type == trailer.ValueType_NUMBER {
		floatValue, err := strconv.ParseFloat(string(V.GetValue()), 64)
		if err != nil {
			glogger.GLogger.Error(err)
			return 0
		}
		return floatValue
	} // Bool 允许两种表示形式
	if V.Type == trailer.ValueType_BOOL {
		if string(V.Value) == "true" {
			return true
		}
		if string(V.Value) == "false" {
			return false
		}
		if string(V.Value) == "1" {
			return true
		}
		if string(V.Value) == "0" {
			return false
		}
		if len(V.Value) > 0 {
			if V.Value[0] == 0 {
				return false
			}
			if V.Value[0] == 1 {
				return true
			}
		}
		return false
	}
	if V.Type == trailer.ValueType_STRING {
		return string(V.Value)
	}
	// 如果到这里说明已经出问题了, 直接返回nil
	return nil
}