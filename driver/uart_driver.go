package driver

import (
	"context"
	"rulex/typex"
	"time"

	"github.com/goburrow/serial"
	"github.com/ngaut/log"
)

//------------------------------------------------------------------------
// 内部函数
//------------------------------------------------------------------------

//
// 正点原子的 Lora 模块封装
//
type UartDriver struct {
	state      typex.DriverState
	serialPort serial.Port
	ctx        context.Context
	In         *typex.InEnd
	RuleEngine typex.RuleX
}

//
// 初始化一个驱动
//
func NewUartDriver(serialPort serial.Port, in *typex.InEnd, e typex.RuleX) typex.XExternalDriver {
	m := &UartDriver{}
	m.In = in
	m.RuleEngine = e
	m.serialPort = serialPort
	m.ctx = context.Background()
	m.state = typex.STOP
	return m
}

//
//
//
func (a *UartDriver) Init() error {
	return nil
}
func (a *UartDriver) Work() error {

	go func(ctx context.Context) {
		acc := 0
		ticker := time.NewTicker(time.Duration(time.Microsecond * 100))
		buffer := [512]byte{}
		for {
			<-ticker.C
			select {
			case <-ctx.Done():
				{
					break
				}
			default:
				{
				}
			}
			data := make([]byte, 1)
			size, err0 := a.serialPort.Read(data)
			if err0 != nil {
				// log.Error("UartDriver error: ", err0)
				if a.state == typex.STOP {
					return
				}
			}
			if size == 1 {
				if data[0] == '#' {
					// log.Info("bytes => ", string(buffer[:acc]), buffer[:acc], acc)
					a.RuleEngine.PushQueue(typex.QueueData{
						In:   a.In,
						Out:  nil,
						E:    a.RuleEngine,
						Data: string(buffer[1:acc]),
					})
					// 重新初始化缓冲区
					for i := 0; i < acc-1; i++ {
						buffer[i] = 0
					}
					acc = 0
				}

				if (data[0] != 0) && (data[0] != '\r') && (data[0] != '\n') {
					buffer[acc] = data[0]
					acc += 1
				}
			}
		}

	}(a.ctx)
	a.state = typex.RUNNING
	return nil

}
func (a *UartDriver) State() typex.DriverState {
	return a.state

}
func (a *UartDriver) Stop() error {
	a.state = typex.STOP
	a.ctx.Done()
	return a.serialPort.Close()
}

func (a *UartDriver) Test() error {
	return nil
}

//
func (a *UartDriver) Read(b []byte) (int, error) {
	return a.serialPort.Read(b)
}

//
func (a *UartDriver) Write(b []byte) (int, error) {
	n, err := a.serialPort.Write(b)
	if err != nil {
		log.Error(err)
		return 0, err
	} else {
		return n, nil
	}

}
func (a *UartDriver) DriverDetail() *typex.DriverDetail {
	return &typex.DriverDetail{
		Name:        "UartDriver",
		Type:        "UartDriver",
		Description: "UartDriver",
	}
}