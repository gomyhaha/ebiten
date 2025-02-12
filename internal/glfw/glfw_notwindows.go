// Copyright 2018 The Ebiten Authors
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

// +build !windows
// +build !js

package glfw

import (
	"image"
	"sync"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type windows map[*glfw.Window]*Window

var (
	theWindows = windows{}
	windowsM   sync.Mutex
)

func (w windows) add(win *glfw.Window) *Window {
	if win == nil {
		return nil
	}
	ww := &Window{win}
	windowsM.Lock()
	w[win] = ww
	windowsM.Unlock()
	return ww
}

func (w windows) remove(win *glfw.Window) {
	windowsM.Lock()
	delete(w, win)
	windowsM.Unlock()
}

func (w windows) get(win *glfw.Window) *Window {
	if win == nil {
		return nil
	}
	windowsM.Lock()
	ww := w[win]
	windowsM.Unlock()
	return ww
}

type Monitor struct {
	m *glfw.Monitor
}

func (m *Monitor) GetPos() (x, y int) {
	return m.m.GetPos()
}

func (m *Monitor) GetVideoMode() *VidMode {
	v := m.m.GetVideoMode()
	if v == nil {
		return nil
	}
	return &VidMode{
		Width:       v.Width,
		Height:      v.Height,
		RedBits:     v.RedBits,
		GreenBits:   v.GreenBits,
		BlueBits:    v.BlueBits,
		RefreshRate: v.RefreshRate,
	}
}

type Window struct {
	w *glfw.Window
}

func (w *Window) Destroy() {
	w.w.Destroy()
	theWindows.remove(w.w)
}

func (w *Window) GetAttrib(attrib Hint) int {
	return w.w.GetAttrib(glfw.Hint(attrib))
}

func (w *Window) GetCursorPos() (x, y float64) {
	return w.w.GetCursorPos()
}

func (w *Window) GetInputMode(mode InputMode) int {
	return w.w.GetInputMode(glfw.InputMode(mode))
}

func (w *Window) GetKey(key Key) Action {
	return Action(w.w.GetKey(glfw.Key(key)))
}

func (w *Window) GetMonitor() *Monitor {
	m := w.w.GetMonitor()
	if m == nil {
		return nil
	}
	return &Monitor{m}
}

func (w *Window) GetMouseButton(button MouseButton) Action {
	return Action(w.w.GetMouseButton(glfw.MouseButton(button)))
}

func (w *Window) GetPos() (x, y int) {
	return w.w.GetPos()
}

func (w *Window) GetSize() (width, height int) {
	return w.w.GetSize()
}

func (w *Window) MakeContextCurrent() {
	w.w.MakeContextCurrent()
}

func (w *Window) SetCharModsCallback(cbfun CharModsCallback) (previous CharModsCallback) {
	var gcb glfw.CharModsCallback
	if cbfun != nil {
		gcb = func(window *glfw.Window, char rune, mods glfw.ModifierKey) {
			cbfun(theWindows.get(window), char, ModifierKey(mods))
		}
	}
	w.w.SetCharModsCallback(gcb)
	return nil // TODO
}

func (w *Window) SetFramebufferSizeCallback(cbfun FramebufferSizeCallback) (previous FramebufferSizeCallback) {
	var gcb glfw.FramebufferSizeCallback
	if cbfun != nil {
		gcb = func(window *glfw.Window, width int, height int) {
			cbfun(theWindows.get(window), width, height)
		}
	}
	w.w.SetFramebufferSizeCallback(gcb)
	return nil // TODO
}

func (w *Window) SetScrollCallback(cbfun ScrollCallback) (previous ScrollCallback) {
	var gcb glfw.ScrollCallback
	if cbfun != nil {
		gcb = func(window *glfw.Window, xoff float64, yoff float64) {
			cbfun(theWindows.get(window), xoff, yoff)
		}
	}
	w.w.SetScrollCallback(gcb)
	return nil // TODO
}

func (w *Window) SetSizeCallback(cbfun SizeCallback) (previous SizeCallback) {
	var gcb glfw.SizeCallback
	if cbfun != nil {
		gcb = func(window *glfw.Window, width, height int) {
			cbfun(theWindows.get(window), width, height)
		}
	}
	w.w.SetSizeCallback(gcb)
	return nil // TODO
}

func (w *Window) SetIcon(images []image.Image) {
	w.w.SetIcon(images)
}

func (w *Window) SetInputMode(mode InputMode, value int) {
	w.w.SetInputMode(glfw.InputMode(mode), value)
}

func (w *Window) SetMonitor(monitor *Monitor, xpos, ypos, width, height, refreshRate int) {
	var m *glfw.Monitor
	if monitor != nil {
		m = monitor.m
	}
	w.w.SetMonitor(m, xpos, ypos, width, height, refreshRate)
}

func (w *Window) SetPos(xpos, ypos int) {
	w.w.SetPos(xpos, ypos)
}

func (w *Window) SetSize(width, height int) {
	w.w.SetSize(width, height)
}

func (w *Window) SetTitle(title string) {
	w.w.SetTitle(title)
}

func (w *Window) ShouldClose() bool {
	return w.w.ShouldClose()
}

func (w *Window) Show() {
	w.w.Show()
}

func (w *Window) SwapBuffers() {
	w.w.SwapBuffers()
}

func CreateWindow(width, height int, title string, monitor *Monitor, share *Window) (*Window, error) {
	var gm *glfw.Monitor
	if monitor != nil {
		gm = monitor.m
	}
	var gw *glfw.Window
	if share != nil {
		gw = share.w
	}

	w, err := glfw.CreateWindow(width, height, title, gm, gw)
	if err != nil {
		return nil, err
	}
	return theWindows.add(w), nil
}

func GetJoystickAxes(joy Joystick) []float32 {
	return glfw.GetJoystickAxes(glfw.Joystick(joy))
}

func GetJoystickButtons(joy Joystick) []byte {
	return glfw.GetJoystickButtons(glfw.Joystick(joy))
}

func GetMonitors() []*Monitor {
	ms := []*Monitor{}
	for _, m := range glfw.GetMonitors() {
		if m != nil {
			ms = append(ms, &Monitor{m})
		} else {
			ms = append(ms, nil)
		}
	}
	return ms
}

func GetPrimaryMonitor() *Monitor {
	m := glfw.GetPrimaryMonitor()
	if m == nil {
		return nil
	}
	return &Monitor{m}
}

func Init() error {
	return glfw.Init()
}

func JoystickPresent(joy Joystick) bool {
	return glfw.JoystickPresent(glfw.Joystick(joy))
}

func PollEvents() {
	glfw.PollEvents()
}

func SetMonitorCallback(cbfun func(monitor *Monitor, event MonitorEvent)) {
	var gcb func(monitor *glfw.Monitor, event glfw.MonitorEvent)
	if cbfun != nil {
		gcb = func(monitor *glfw.Monitor, event glfw.MonitorEvent) {
			var m *Monitor
			if monitor != nil {
				m = &Monitor{monitor}
			}
			cbfun(m, MonitorEvent(event))
		}
	}
	glfw.SetMonitorCallback(gcb)
}

func SwapInterval(interval int) {
	glfw.SwapInterval(interval)
}

func Terminate() {
	glfw.Terminate()
}

func WindowHint(target Hint, hint int) {
	glfw.WindowHint(glfw.Hint(target), hint)
}
