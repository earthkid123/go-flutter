package main

import (
	"bufio"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-flutter-desktop/go-flutter"
	"github.com/go-flutter-desktop/go-flutter/embedder"

	"github.com/go-gl/glfw/v3.2/glfw"
)

func setIcon(window *glfw.Window) error {
	_, currentFilePath, _, _ := runtime.Caller(0)
	dir := path.Dir(currentFilePath)
	imgFile, err := os.Open(dir + "/assets/icon.png")
	if err != nil {
		return err
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return err
	}
	window.SetIcon([]image.Image{img})
	return nil
}

func main() {

	_, currentFilePath, _, _ := runtime.Caller(0)
	dir := path.Dir(currentFilePath)

	initialApplicationHeight := 600
	initialApplicationWidth := 800

	options := []flutter.Option{
		flutter.ProjectAssetPath(dir + "/flutter_project/demo/build/flutter_assets"),

		// This path should not be changed. icudtl.dat is handled by engineDownloader.go
		flutter.ApplicationICUDataPath(dir + "/icudtl.dat"),

		flutter.ApplicationWindowDimension(initialApplicationWidth, initialApplicationHeight),
		// gutter.OptionPixelRatio(1.2),
		flutter.OptionWindowInitializer(setIcon),
		flutter.OptionVMArguments([]string{
			// "--disable-dart-asserts", // release mode flag
			// "--disable-observatory",
			"--observatory-port=50300",
		}),

		flutter.OptionAddPluginReceiver(ownPlugin, "plugin_demo"),

		// Default keyboard is Qwerty, if you want to change it, you can check keyboard.go in gutter package.
		// Otherwise you can create your own by usinng `KeyboardShortcuts` struct.
		//flutter.OptionKeyboardLayout(flutter.KeyboardAzertyLayout),
	}

	if err := flutter.Run(options...); err != nil {
		log.Fatalln(err)
	}

}

// Plugin that read the stdin and send the number to the dart side
func ownPlugin(
	platMessage *embedder.PlatformMessage,
	flutterEngine *embedder.FlutterEngine,
	window *glfw.Window,
) bool {

	if platMessage.Message.Method != "getNumber" {
		log.Printf("Unhandled platform method: %#v from channel %#v\n",
			platMessage.Message.Method, platMessage.Channel)
		return false
	}

	time.Sleep(2 * time.Second)
	go func() {
		fmt.Printf("Reading (A number): ")
		reader := bufio.NewReader(os.Stdin)
		s, _ := reader.ReadString('\n')
		s = strings.TrimRight(s, "\r\n")
		if _, err := strconv.Atoi(s); err == nil {

			flutterEngine.SendPlatformMessageResponse(platMessage, []byte("[ "+s+" ]"))

		}
	}()

	return true

}
