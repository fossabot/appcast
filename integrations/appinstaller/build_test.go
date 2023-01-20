package appinstaller_test

import (
	"context"
	"io"
	"testing"

	"github.com/abemedia/appcast/integrations/appinstaller"
	fileSource "github.com/abemedia/appcast/source/blob/file"
	fileTarget "github.com/abemedia/appcast/target/file"
	"github.com/google/go-cmp/cmp"
)

func TestBuild(t *testing.T) {
	dir := t.TempDir()
	src, _ := fileSource.New(fileSource.Config{Path: "../../testdata"})
	tgt, _ := fileTarget.New(fileTarget.Config{Path: dir})

	tests := []struct {
		name   string
		config *appinstaller.Config
		want   string
	}{
		{
			name: "Test-x64.appinstaller",
			config: &appinstaller.Config{
				Version: "1.0.0",
				Source:  src,
				Target:  tgt,
			},
			want: `<?xml version="1.0" encoding="UTF-8"?>
<AppInstaller xmlns="http://schemas.microsoft.com/appx/appinstaller/2017" Version="1.0.0.1" Uri="file://` + dir + `/Test-x64.appinstaller">
	<MainPackage Name="Test" Publisher="CN=Test" Version="1.0.0.1" ProcessorArchitecture="x64" Uri="file:///home/adam/Work/optimus/appcast/testdata/v1.0.0/test.msix" />
</AppInstaller>`,
		},
		{
			name: "Test.appinstaller",
			config: &appinstaller.Config{
				Version:    "1.0.0",
				Source:     src,
				Target:     tgt,
				ShowPrompt: true,
			},
			want: `<?xml version="1.0" encoding="UTF-8"?>
<AppInstaller xmlns="http://schemas.microsoft.com/appx/appinstaller/2018" Version="1.0.0.1" Uri="file://` + dir + `/Test.appinstaller">
	<MainBundle Name="Test" Publisher="CN=Test" Version="1.0.0.1" Uri="file:///home/adam/Work/optimus/appcast/testdata/v1.0.0/test.msixbundle" />
	<UpdateSettings>
		<OnLaunch ShowPrompt="true" />
	</UpdateSettings>
</AppInstaller>`,
		},
		{
			name: "Test-x64.appinstaller",
			config: &appinstaller.Config{
				Version:                 "1.0.0",
				Source:                  src,
				Target:                  tgt,
				AutomaticBackgroundTask: true,
			},
			want: `<?xml version="1.0" encoding="UTF-8"?>
<AppInstaller xmlns="http://schemas.microsoft.com/appx/appinstaller/2017/2" Version="1.0.0.1" Uri="file://` + dir + `/Test-x64.appinstaller">
	<MainPackage Name="Test" Publisher="CN=Test" Version="1.0.0.1" ProcessorArchitecture="x64" Uri="file:///home/adam/Work/optimus/appcast/testdata/v1.0.0/test.msix" />
	<UpdateSettings>
		<AutomaticBackgroundTask />
	</UpdateSettings>
</AppInstaller>`,
		},
		{
			name: "Test-x64.appinstaller",
			config: &appinstaller.Config{
				Version:                   "1.0.0",
				Source:                    src,
				Target:                    tgt,
				HoursBetweenUpdateChecks:  12,
				UpdateBlocksActivation:    true,
				ShowPrompt:                true,
				AutomaticBackgroundTask:   true,
				ForceUpdateFromAnyVersion: true,
			},
			want: `<?xml version="1.0" encoding="UTF-8"?>
<AppInstaller xmlns="http://schemas.microsoft.com/appx/appinstaller/2018" Version="1.0.0.1" Uri="file://` + dir + `/Test-x64.appinstaller">
	<MainPackage Name="Test" Publisher="CN=Test" Version="1.0.0.1" ProcessorArchitecture="x64" Uri="file:///home/adam/Work/optimus/appcast/testdata/v1.0.0/test.msix" />
	<UpdateSettings>
		<OnLaunch HoursBetweenUpdateChecks="12" UpdateBlocksActivation="true" ShowPrompt="true" />
		<AutomaticBackgroundTask />
		<ForceUpdateFromAnyVersion>true</ForceUpdateFromAnyVersion>
	</UpdateSettings>
</AppInstaller>`,
		},
	}

	for _, test := range tests {
		err := appinstaller.Build(context.Background(), test.config)
		if err != nil {
			t.Fatal(err)
		}

		r, err := test.config.Target.NewReader(context.Background(), test.name)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()

		got, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(test.want, string(got)); diff != "" {
			t.Fatal(diff)
		}
	}
}
