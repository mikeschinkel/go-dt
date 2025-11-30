package test

import (
	"testing"

	"github.com/mikeschinkel/go-dt"
	"github.com/mikeschinkel/go-dt/appinfo"
)

func TestNew(t *testing.T) {
	args := appinfo.Args{
		Name:        "TestApp",
		Description: "A test application",
		Version:     dt.Version("1.0.0"),
		AppSlug:     dt.PathSegment("testapp"),
		ConfigSlug:  dt.PathSegment("testapp-config"),
		ConfigFile:  dt.RelFilepath("config/app.json"),
		ExeName:     dt.Filename("testapp"),
		LogFile:     dt.Filename("testapp.log"),
		LogPath:     dt.PathSegments("logs"),
		InfoURL:     dt.URL("https://example.com/testapp"),
		ExtraInfo:   map[string]any{"key": "value"},
	}

	ai := appinfo.New(args)

	if ai.Name() != args.Name {
		t.Errorf("Name() = %v, want %v", ai.Name(), args.Name)
	}
	if ai.Description() != args.Description {
		t.Errorf("Description() = %v, want %v", ai.Description(), args.Description)
	}
	if ai.Version() != args.Version {
		t.Errorf("Version() = %v, want %v", ai.Version(), args.Version)
	}
	if ai.AppSlug() != args.AppSlug {
		t.Errorf("AppSlug() = %v, want %v", ai.AppSlug(), args.AppSlug)
	}
	if ai.ConfigSlug() != args.ConfigSlug {
		t.Errorf("ConfigSlug() = %v, want %v", ai.ConfigSlug(), args.ConfigSlug)
	}
	if ai.ConfigFile() != args.ConfigFile {
		t.Errorf("ConfigFile() = %v, want %v", ai.ConfigFile(), args.ConfigFile)
	}
	if ai.ExeName() != args.ExeName {
		t.Errorf("ExeName() = %v, want %v", ai.ExeName(), args.ExeName)
	}
	if ai.LogFile() != args.LogFile {
		t.Errorf("LogFile() = %v, want %v", ai.LogFile(), args.LogFile)
	}
	if ai.LogPath() != args.LogPath {
		t.Errorf("LogPath() = %v, want %v", ai.LogPath(), args.LogPath)
	}
	if ai.InfoURL() != args.InfoURL {
		t.Errorf("InfoURL() = %v, want %v", ai.InfoURL(), args.InfoURL)
	}

	extraInfo := ai.ExtraInfo()
	if extraInfo["key"] != "value" {
		t.Errorf("ExtraInfo()[\"key\"] = %v, want %v", extraInfo["key"], "value")
	}
}

func TestNewWithNilExtraInfo(t *testing.T) {
	args := appinfo.Args{
		Name:    "TestApp",
		Version: dt.Version("1.0.0"),
		// ExtraInfo is nil
	}

	ai := appinfo.New(args)

	extraInfo := ai.ExtraInfo()
	if extraInfo == nil {
		t.Error("ExtraInfo() should not be nil, expected empty map")
	}
	if len(extraInfo) != 0 {
		t.Errorf("ExtraInfo() length = %v, want 0", len(extraInfo))
	}
}
