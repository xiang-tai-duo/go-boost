@echo off
pushd "%~dp0"
set GO_BOOST_ROOT=C:\
set GO_BOOST=..\go-boost\
if not exist "go-boost" (
    mklink /d go-boost "%GO_BOOST_ROOT%go-boost\"
)
if not exist "__build.cmd" (
   copy "%GO_BOOST%__build.cmd" "__build.cmd"
)
if not exist "__build_electron.cmd" (
   copy "%GO_BOOST%electron\__build_electron.cmd" "__build_electron.cmd"
)
if not exist "src" (
    mkdir "src"
)
pushd src
    if not exist "winres" (
        mkdir "winres"
        pushd winres
            if not exist "icon.png" (
                copy "%GO_BOOST%winres\icon.png" "icon.png"
            )
            if not exist "winres.json" (
                copy "%GO_BOOST%winres\winres.json" "winres.json"
            )
        popd
    )
    if not exist "electron" (
        mkdir "electron"
        pushd electron
            copy "%GO_BOOST%electron\app.png.in" "app.png"
            copy "%GO_BOOST%electron\package.json.in" "package.json"
        popd
    )
    if not exist "wwwroot" (
        mkdir "wwwroot"
        xcopy "%GO_BOOST%electron\wwwroot\" "wwwroot\" /e /c /h /i /y
    )
popd

:: require github.com/xiang-tai-duo/go-bootstrap v0.0.0
:: replace github.com/xiang-tai-duo/go-bootstrap => ../go-bootstrap
