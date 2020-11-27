package main

import (
	"context"
	"errors"
	addon "github.com/galenliu/gateway-addon-golang"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yeelight-adapter/lib"
)

var (
	On               = "on"
	Brightness       = "brightness"
	Hue              = "hue"
	ColorTemperature = "ct"
	ColorModel       = "ColorMode"
)

func main() {

	var adapter = NewYeeLightAdapter("yeeLight-adapter", "yeeLight-adapter")
	adapter.StartPairing(2000)

	var systemCallCloseFunc = func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-c
		if adapter != nil {
			adapter.CloseProxy()
		}
		os.Exit(0)
	}

	go systemCallCloseFunc()

	for {
		if !adapter.IsProxyRunning() {
			time.Sleep(time.Duration(2) * time.Second)
			return
		}
	}

}

type YeeLightAdapter struct {
	*addon.AdapterProxy
}

type YeeLightDevice struct {
	*addon.DeviceProxy
	lightLib *lib.Light
}

type YeeLightProperty struct {
	*addon.PropertyProxy
	*YeeLightDevice
}

func NewYeeLightAdapter(id, packageName string) *YeeLightAdapter {
	adapter := &YeeLightAdapter{
		addon.NewAdapterProxy(id, packageName),
	}
	adapter.StartPairing(10)
	return adapter
}

func NewYeeLightProperty(proxy *addon.PropertyProxy, device *YeeLightDevice) *YeeLightProperty {
	return &YeeLightProperty{PropertyProxy: proxy, YeeLightDevice: device}
}

func NewYeeLightDevice(proxy *addon.DeviceProxy, light *lib.Light) *YeeLightDevice {
	return &YeeLightDevice{proxy, light}
}

func (adapter *YeeLightAdapter) StartPairing(timeout int) {
	ctx := context.Background()
	ctx1, _ := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	var pairing = func(ctx context.Context) {
		select {
		case <-ctx1.Done():
			adapter.CancelParing()
			return
		default:
			adapter.AdapterProxy.StartPairing(timeout)
			lights := lib.Discover()
			devs := initDevices(adapter, lights)
			for _, d := range devs {
				adapter.HandleDeviceAdded(d)
				return
			}

		}
	}
	go pairing(ctx1)
}

func initDevices(adapter addon.IAdapter, lights []lib.Light) (YeeLights []*YeeLightDevice) {

	for _, light := range lights {
		var yeeLight = NewYeeLightDevice(addon.NewDevice(adapter, light.ID), &light)
		yeeLight.AppendType("Light")
		prop := addon.NewColorModeProperty(ColorModel, light.ColorMode-1)
		var p = NewYeeLightProperty(prop, yeeLight)
		yeeLight.AddProperty(p)

		for _, prop := range light.Support {
			switch prop {
			case "set_power":
				proxyProp := addon.NewOnOffProperty(On)

				var _p = NewYeeLightProperty(proxyProp, yeeLight)

				yeeLight.AddProperty(_p)
				_p.SetValue(true)
			case "set_bright":
				p := addon.NewBrightnessProperty(Brightness, 0, 100)
				var _p = NewYeeLightProperty(p, yeeLight)
				_p.SetValue(10)
				yeeLight.AddProperty(_p)
			case "set_rgb":
				p := addon.NewColorProperty(ColorTemperature)
				var _p = NewYeeLightProperty(p, yeeLight)
				yeeLight.AddProperty(_p)
				_p.SetValue("#FF0000")
				yeeLight.AddProperty(_p)
			case "set_ct_abx":
				p := addon.NewColorTemperatureProperty(ColorTemperature, 170000, 680000)
				var _p = NewYeeLightProperty(p, yeeLight)
				_p.SetValue(680000)
				yeeLight.AddProperty(_p)
			default:
				continue
			}

		}
		YeeLights = append(YeeLights, yeeLight)

	}
	return YeeLights

}

func (prop *YeeLightProperty) SetValue(value interface{}) error {
	switch prop.Name {
	case On:
		_ = prop.PropertyProxy.SetValue(value)
		v, ok := value.(bool)
		if !ok {
			return errors.New("value err")

		}
		if v {
			_, _ = prop.lightLib.PowerOn(0)
		} else {
			_, _ = prop.lightLib.PowerOff(0)
		}
	case ColorTemperature:
		_ = prop.PropertyProxy.SetValue(value)
		v, ok := value.(int)
		if !ok {
			return errors.New("value err")

		}
		prop.lightLib.SetTemp(v, 0)
	case Hue:
		prop.PropertyProxy.SetValue(value)
		v, ok := value.(string)
		if !ok {
			return errors.New("value err")

		}
		rgb, err := lib.HTMLToRGB(v)
		if err != nil {
			return err
		}
		prop.lightLib.SetRGB(rgb.R, rgb.G, rgb.B, 0)

	}
	return nil
}
